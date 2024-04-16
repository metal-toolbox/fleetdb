//nolint:all  // XXX remove this!
package inventory

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bmc-toolbox/common"
	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/models"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
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
// The caller needs to compose the correct attributes for its given component.
func composeRecords(ctx context.Context, exec boil.ContextExecutor, cmn *common.Common,
	inband bool, deviceID, slug string, attr *attributes) error {
	typeID := dbtools.MustComponentTypeID(ctx, exec, slug)

	sc := &models.ServerComponent{
		Name:                  null.StringFrom(slug),
		Vendor:                null.NewString(cmn.Vendor, cmn.Vendor != ""),
		Model:                 null.NewString(cmn.Model, cmn.Model != ""),
		Serial:                null.NewString(cmn.Serial, cmn.Serial != ""),
		ServerID:              deviceID,
		ServerComponentTypeID: typeID,
	}

	prodName := strings.ToLower(strings.TrimSpace(cmn.ProductName))
	if sc.Model.IsZero() && prodName != "" {
		sc.Model.SetValid(prodName)
	}

	if sc.Serial.IsZero() {
		sc.Serial = null.StringFrom("0")
	}

	if err := createOrUpdateComponent(ctx, exec, sc); err != nil {
		return errors.Wrap(errComponent, slug+": "+err.Error())
	}

	// avoid computing this twice
	attr.ProductName = prodName

	attrData := attr.MustJSON()

	// update the component attribute
	if err := updateAnyAttribute(ctx, exec, false, sc.ID, getAttributeNamespace(inband), attrData); err != nil {
		return errors.Wrap(errAttribute, slug+": "+err.Error())
	}

	// every component with firmware gets a firmware versioned attribute
	if cmn.Firmware != nil {
		payload := mustFirmwareJSON(cmn.Firmware)

		if err := updateAnyVersionedAttribute(ctx, exec, false, sc.ID, getFirmwareNamespace(inband), payload); err != nil {
			return errors.Wrap(errVersionedAttr, slug+"-firmware: "+err.Error())
		}
	}

	// every component with status gets a status versioned attribute
	if cmn.Status != nil {
		payload := mustStatusJSON(cmn.Status)

		if err := updateAnyVersionedAttribute(ctx, exec, false, sc.ID, getStatusNamespace(inband), payload); err != nil {
			return errors.Wrap(errVersionedAttr, slug+"-status: "+err.Error())
		}
	}

	return nil
}

func retrieveComponentAttributes(ctx context.Context, exec boil.ContextExecutor,
	componentID, namespace string) (*attributes, error) {
	ar, err := models.Attributes(
		models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(componentID)),
		models.AttributeWhere.Namespace.EQ(namespace),
	).One(ctx, exec)

	if err != nil {
		return nil, err
	}

	attr := &attributes{}
	if err := attr.FromJSON(ar.Data); err != nil {
		return nil, err
	}
	return attr, nil
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
		return nil, err
	}
	return fwr.Data, nil
}

func retrieveComponentFirmwareVA(ctx context.Context, exec boil.ContextExecutor,
	parentID, slug, namespace string) (*common.Firmware, error) {
	data, err := retrieveVersionedAttribute(ctx, exec, parentID, namespace, false)
	if err != nil {
		return nil, err
	}
	fw, err := firmwareFromJSON(data)
	if err != nil {
		return nil, err
	}
	return fw, nil
}

func retrieveComponentStatusVA(ctx context.Context, exec boil.ContextExecutor, parentID,
	slug, namespace string) (*common.Status, error) {
	data, err := retrieveVersionedAttribute(ctx, exec, parentID, namespace, false)
	if err != nil {
		return nil, err
	}
	st, err := statusFromJSON(data)
	if err != nil {
		return nil, err
	}
	return st, nil
}

// this is returned by more-or-less generic database routines and the caller can specialize into a type
type dbComponent struct {
	cmn  *common.Common
	attr *attributes
}

// As generically as possible, retrieve this component from the database. Status and Firmware
// are composed into the *Common. We return the attributes so that the caller can reconsitute
// the specific device type. This is basically the reverse of composeRecords.
func componentsFromDatabase(ctx context.Context, exec boil.ContextExecutor,
	inband bool, deviceID, slug string) ([]*dbComponent, error) {
	records, err := models.ServerComponents(
		models.ServerComponentWhere.Name.EQ(null.StringFrom(slug)),
		models.ServerComponentWhere.ServerID.EQ(deviceID),
		qm.OrderBy(models.ServerComponentColumns.CreatedAt+" DESC"),
	).All(ctx, exec)

	if err != nil {
		return nil, err
	}

	comps := []*dbComponent{}

	for _, rec := range records {
		// We should always have attributes, even if it's only "ProductName" (b/c it comes from common)
		attr, err := retrieveComponentAttributes(ctx, exec, rec.ID, getAttributeNamespace(inband))
		if err != nil {
			return nil, err
		}

		// Either firmware or status might have no stored data. That's fine.
		fw, err := retrieveComponentFirmwareVA(ctx, exec, rec.ID, slug, getFirmwareNamespace(inband))
		switch err {
		case nil, sql.ErrNoRows:
		default:
			return nil, err
		}

		st, err := retrieveComponentStatusVA(ctx, exec, rec.ID, slug, getStatusNamespace(inband))
		switch err {
		case nil, sql.ErrNoRows:
		default:
			return nil, err
		}
		// Despite the schema, serial is required. It is set on storing the component if it was empty coming in.
		serial := rec.Serial.String
		comp := &dbComponent{
			cmn: &common.Common{
				Vendor:      rec.Vendor.String,
				Model:       rec.Model.String,
				Serial:      serial,
				ProductName: attr.ProductName,
				Firmware:    fw,
				Status:      st,
			},
			attr: attr,
		}
		comps = append(comps, comp)
	}

	return comps, nil
}

func (dv *DeviceView) ComposeComponents(ctx context.Context, exec boil.ContextExecutor) error {
	if err := dv.writeBios(ctx, exec); err != nil {
		return err
	}
	if err := dv.writeBMC(ctx, exec); err != nil {
		return err
	}
	if err := dv.writeMainboard(ctx, exec); err != nil {
		return err
	}
	if err := dv.writeDimms(ctx, exec); err != nil {
		return err
	}
	return nil
}

func (dv *DeviceView) writeBios(ctx context.Context, exec boil.ContextExecutor) error {
	bios := dv.Inv.BIOS

	attr := &attributes{
		Capabilities:  bios.Capabilities,
		CapacityBytes: bios.CapacityBytes,
		Description:   bios.Description,
		Metadata:      bios.Metadata,
		Oem:           bios.Oem,
		SizeBytes:     bios.SizeBytes,
	}

	return composeRecords(ctx, exec, &bios.Common, dv.Inband, dv.DeviceID.String(), common.SlugBIOS, attr)
}

func (dv *DeviceView) getBios(ctx context.Context, exec boil.ContextExecutor) error {
	bios := &common.BIOS{}
	components, err := componentsFromDatabase(ctx, exec, dv.Inband, dv.DeviceID.String(), common.SlugBIOS)
	if err != nil {
		return err
	}
	// We should never have more BIOS component, but a defect could result in multiple records. The
	// components slice should come back in order of most recent records first, so first record wins.
	for _, comp := range components {
		bios.Common = *comp.cmn
		bios.Capabilities = comp.attr.Capabilities
		bios.CapacityBytes = comp.attr.CapacityBytes
		bios.Description = comp.attr.Description
		bios.Oem = comp.attr.Oem
		bios.SizeBytes = comp.attr.SizeBytes
		break
	}
	dv.Inv.BIOS = bios
	return nil
}

func (dv *DeviceView) writeBMC(ctx context.Context, exec boil.ContextExecutor) error {
	bmc := dv.Inv.BMC

	attr := &attributes{
		Capabilities: bmc.Capabilities,
		Description:  bmc.Description,
		Metadata:     bmc.Metadata,
		Oem:          bmc.Oem,
	}

	return composeRecords(ctx, exec, &bmc.Common, dv.Inband, dv.DeviceID.String(), common.SlugBMC, attr)
}

func (dv *DeviceView) getBMC(ctx context.Context, exec boil.ContextExecutor) error {
}

func (dv *DeviceView) writeMainboard(ctx context.Context, exec boil.ContextExecutor) error {
	mb := dv.Inv.Mainboard

	attr := &attributes{
		Capabilities: mb.Capabilities,
		Description:  mb.Description,
		Metadata:     mb.Metadata,
		Oem:          mb.Oem,
		PhysicalID:   mb.PhysicalID,
	}

	return composeRecords(ctx, exec, &mb.Common, dv.Inband, dv.DeviceID.String(), common.SlugMainboard, attr)
}

func (dv *DeviceView) writeDimms(ctx context.Context, exec boil.ContextExecutor) error {
	for idx, dimm := range dv.Inv.Memory {
		// skip bogus dimms
		if dimm.Vendor == "" &&
			dimm.ProductName == "" &&
			dimm.SizeBytes == 0 &&
			dimm.ClockSpeedHz == 0 {
			continue
		}

		if strings.TrimSpace(dimm.Serial) == "" {
			dimm.Serial = fmt.Sprintf("%d", idx)
		}

		attr := &attributes{
			Capabilities: dimm.Capabilities,
			ClockSpeedHz: dimm.ClockSpeedHz,
			Description:  dimm.Description,
			FormFactor:   dimm.FormFactor,
			Metadata:     dimm.Metadata, // maybe this should be versioned?
			PartNumber:   dimm.PartNumber,
			SizeBytes:    dimm.SizeBytes,
			Slot:         strings.TrimPrefix(dimm.Slot, "DIMM.Socket."),
		}

		if err := composeRecords(ctx, exec, &dimm.Common, dv.Inband, dv.DeviceID.String(), common.SlugPhysicalMem, attr); err != nil {
			return err
		}
	}
	return nil
}
