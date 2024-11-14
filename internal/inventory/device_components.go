package inventory

import (
	"context"
	"database/sql"
	"encoding/json"

	common "github.com/metal-toolbox/bmc-common"
	rivets "github.com/metal-toolbox/rivets/v2/types"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go.uber.org/zap"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/metrics"
	"github.com/metal-toolbox/fleetdb/internal/models"
)

var (
	errComponentType = errors.New("component type error")
	errComponent     = errors.New("component error")
	errAttribute     = errors.New("attribute error")
	errVersionedAttr = errors.New("versioned attribute error")

	inbandNSTag    = "sh.hollow.alloy.inband"
	outofbandNSTag = "sh.hollow.alloy.outofband"
)

func getNamespace(inband bool) string {
	ns := outofbandNSTag
	if inband {
		ns = inbandNSTag
	}
	return ns
}

func getAttributeNamespace(inband bool) string {
	return getNamespace(inband) + ".metadata"
}

func getFirmwareNamespace(inband bool) string {
	return getNamespace(inband) + ".firmware"
}

func getStatusNamespace(inband bool) string {
	return getNamespace(inband) + ".status"
}

func createOrUpdateComponent(ctx context.Context, exec boil.ContextExecutor, sc *models.ServerComponent) error {
	existing, err := models.ServerComponents(
		models.ServerComponentWhere.Name.EQ(sc.Name),
		models.ServerComponentWhere.ServerID.EQ(sc.ServerID),
		models.ServerComponentWhere.Serial.EQ(sc.Serial),
		models.ServerComponentWhere.ServerComponentTypeID.EQ(sc.ServerComponentTypeID),
	).One(ctx, exec)

	switch err {
	case nil:
		sc.ID = existing.ID
		_, updErr := sc.Update(ctx, exec, boil.Infer())
		return updErr
	case sql.ErrNoRows:
		return sc.Insert(ctx, exec, boil.Infer())
	default:
		return err
	}
}

// This encapsulates much of the repetitive work of getting a component to the database layer.
func composeRecords(ctx context.Context, exec boil.ContextExecutor, cmp *rivets.Component,
	inband bool, deviceID string) error {
	name := cmp.Name
	typeID, err := dbtools.ComponentTypeIDFromName(name)
	if err != nil {
		return errors.Wrap(errComponentType, name+" not found")
	}

	sc := &models.ServerComponent{
		Name:                  null.StringFrom(name),
		Vendor:                null.NewString(cmp.Vendor, cmp.Vendor != ""),
		Model:                 null.NewString(cmp.Model, cmp.Model != ""),
		Serial:                null.NewString(cmp.Serial, cmp.Serial != ""),
		ServerID:              deviceID,
		ServerComponentTypeID: typeID,
	}

	if err := createOrUpdateComponent(ctx, exec, sc); err != nil {
		return errors.Wrap(errComponent, name+": "+err.Error())
	}

	if cmp.Attributes != nil {
		attrData := mustAttributesJSON(cmp.Attributes)

		// update the component attribute
		if err := updateAnyAttribute(ctx, exec, false, sc.ID, getAttributeNamespace(inband), attrData); err != nil {
			return errors.Wrap(errAttribute, name+": "+err.Error())
		}
	}

	// every component with firmware gets a firmware versioned attribute
	if cmp.Firmware != nil {
		payload := mustFirmwareJSON(cmp.Firmware)

		if err := updateAnyVersionedAttribute(ctx, exec, false, sc.ID, getFirmwareNamespace(inband), payload); err != nil {
			return errors.Wrap(errVersionedAttr, name+"-firmware: "+err.Error())
		}
	}

	// every component with status gets a status versioned attribute
	if cmp.Status != nil {
		payload := mustStatusJSON(cmp.Status)

		if err := updateAnyVersionedAttribute(ctx, exec, false, sc.ID, getStatusNamespace(inband), payload); err != nil {
			return errors.Wrap(errVersionedAttr, name+"-status: "+err.Error())
		}
	}

	return nil
}

func retrieveComponentAttributes(ctx context.Context, exec boil.ContextExecutor,
	componentID, namespace string) (*rivets.ComponentAttributes, error) {
	ar, err := models.Attributes(
		models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(componentID)),
		models.AttributeWhere.Namespace.EQ(namespace),
	).One(ctx, exec)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
	default:
		metrics.DBError("fetch component attributes")
		return nil, err
	}

	return componentAttributesFromJSON(ar.Data)
}

// we rely on the caller to know what v-attributes to retrieve and how to deserialize that data
func retrieveVersionedAttribute(ctx context.Context, exec boil.ContextExecutor,
	parentID, namespace string, isServer bool) ([]byte, error) {
	var mods []qm.QueryMod
	if isServer {
		mods = append(mods, models.VersionedAttributeWhere.ServerID.EQ(null.StringFrom(parentID)))
	} else {
		mods = append(mods, models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(parentID)))
	}
	mods = append(mods,
		models.VersionedAttributeWhere.Namespace.EQ(namespace),
		qm.OrderBy("tally DESC"), // get the most recent record
	)

	fwr, err := models.VersionedAttributes(mods...).One(ctx, exec)
	if err != nil {
		metrics.DBError("fetch versioned attributes")
		return nil, err
	}
	return fwr.Data, nil
}

func retrieveComponentFirmwareVA(ctx context.Context, exec boil.ContextExecutor,
	parentID, namespace string) (*common.Firmware, error) {
	data, err := retrieveVersionedAttribute(ctx, exec, parentID, namespace, false)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
	default:
		return nil, err
	}

	fw, err := firmwareFromJSON(data)
	if err != nil {
		return nil, err
	}
	return fw, nil
}

func retrieveComponentStatusVA(ctx context.Context, exec boil.ContextExecutor, parentID,
	namespace string) (*common.Status, error) {
	data, err := retrieveVersionedAttribute(ctx, exec, parentID, namespace, false)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
	default:
		return nil, err
	}

	st, err := statusFromJSON(data)
	if err != nil {
		return nil, err
	}
	return st, nil
}

// As generically as possible, retrieve this component from the database. Status and Firmware
// are composed into the *Common. We return the attributes so that the caller can reconsitute
// the specific device type. This is basically the reverse of composeRecords.
func componentsFromDatabase(ctx context.Context, exec boil.ContextExecutor,
	inband bool, deviceID string) ([]*rivets.Component, error) {
	records, err := models.ServerComponents(
		models.ServerComponentWhere.ServerID.EQ(deviceID),
		qm.OrderBy(models.ServerComponentColumns.CreatedAt+" DESC"),
	).All(ctx, exec)

	if err != nil {
		return nil, err
	}

	var comps []*rivets.Component

	var ute *json.UnmarshalTypeError
	for _, rec := range records {
		// attributes/firmware/status might not be stored because it was missing in the original data.
		attr, err := retrieveComponentAttributes(ctx, exec, rec.ID, getAttributeNamespace(inband))
		switch {
		case err == nil, errors.Is(err, sql.ErrNoRows):
		case errors.As(err, &ute):
			// attributes are a bit of the wild-west. if the JSON we stored doesn't deserialize
			// cleanly into an attributes structure, just complain about it but don't stop.
			zap.L().With(
				zap.String("server.id", deviceID),
				zap.String("component.id", rec.ID),
				zap.String("component.type", rec.Name.String),
			).Warn("bad json attributes")
		default:
			return nil, errors.Wrap(err, "retrieving "+rec.Name.String+"-"+rec.ID+" attributes"+":"+err.Error())
		}

		fw, err := retrieveComponentFirmwareVA(ctx, exec, rec.ID, getFirmwareNamespace(inband))
		switch err {
		case nil, sql.ErrNoRows:
		default:
			return nil, errors.Wrap(err, "retrieving "+rec.Name.String+"-"+rec.ID+" firmware"+":"+err.Error())
		}

		comp := &rivets.Component{
			Name:       rec.Name.String,
			Vendor:     rec.Vendor.String,
			Model:      rec.Model.String,
			Serial:     rec.Serial.String,
			Firmware:   fw,
			Attributes: attr,
		}

		st, err := retrieveComponentStatusVA(ctx, exec, rec.ID, getStatusNamespace(inband))
		if err != nil {
			// Relax error
			zap.L().With(
				zap.String("rec.ID", rec.ID),
				zap.String("rec.Name", rec.Name.String),
			).Error(err.Error())
		} else {
			comp.Status = st
		}
		comps = append(comps, comp)
	}

	return comps, nil
}
