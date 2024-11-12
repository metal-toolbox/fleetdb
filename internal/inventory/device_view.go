package inventory

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	rivets "github.com/metal-toolbox/rivets/v2/types"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/metrics"
	"github.com/metal-toolbox/fleetdb/internal/models"
)

/*
	XXX From the "this is why we can't have nice things" dept.:
	SQLBoiler does in-fact generate Upsert methods for the objects in its ORM. However, because many of our tables
	have partial constraints on them, Postgres (and by extension CRDB) requires a WHERE clause when you specify
	ON CONFLICT (columns). That is the query looks like: INSERT INTO table (col1, col2) VALUES (foo, bar) ON CONFLICT
	(col1, col2) WHERE col1 <is is-not something> DO UPDATE ...

	SQLBoiler doesn't have a provision for that WHERE clause in the ON CONFLICT and they probably won't add it: cf.
	https://github.com/volatiletech/sqlboiler/issues/856

	That means we do the upserts the hard way until we change change the tables that have only partial constraints.
*/

var (
	// historically these values were determined/set by alloy, even though they are
	// internal to the data storage layer, hence the names
	alloyVendorNamespace = "sh.hollow.alloy.server_vendor_attributes"
	// XXX: enable this when Server supports UEFI variables
	// alloyUefiVarsNamespace = "sh.hollow.alloy.server_uefi_variables" // this is a versioned attribute, we expect it to change
	serverStatusNamespace = "sh.hollow.alloy.server_status" // versioned

	// metadata keys
	modelKey  = "model"
	vendorKey = "vendor"
	serialKey = "serial"
	// XXX: again, enable after UEFI Variables are a thing. uefiVarsKey = "uefi-variables"

	errBadServer    = errors.New("data is missing required field")
	errBadComponent = errors.New("component data")

	ErrNoInventory = errors.New("no inventory stored")
)

// DeviceView encapsulates everything we need to get and set inventory data for servers
// A reminder for maintenance: this type needs to be able to contain all the
// relevant fields from Component-Inventory or Alloy.
type DeviceView struct {
	Inv      *rivets.Server
	Inband   bool // the method of inventory collection
	DeviceID uuid.UUID
}

// ServerSanityCheck handles verifying that all the details in the incoming data
// structure have been set so we maintain database invariants.
func ServerSanityCheck(srv *rivets.Server) error {
	srvReq := map[string]string{
		"model":  srv.Model,
		"vendor": srv.Vendor,
		"serial": srv.Serial,
	}
	for k, v := range srvReq {
		if v == "" {
			return errors.Wrap(errBadServer, k)
		}
	}

	for _, cmp := range srv.Components {
		if _, err := dbtools.ComponentTypeIDFromName(cmp.Name); err != nil {
			return errors.Wrap(errBadComponent, err.Error())
		}
		if cmp.Serial == "" {
			return errors.Wrap(errBadComponent, cmp.Name+" missing serial")
		}
	}
	return nil
}

func (dv *DeviceView) vendorAttributes() json.RawMessage {
	m := map[string]string{
		modelKey:  dv.Inv.Model,
		serialKey: dv.Inv.Serial,
		vendorKey: dv.Inv.Vendor,
	}
	byt, _ := json.Marshal(m)

	return byt
}

/* XXX: return this when rivet's Server datatype has a facility to store UEFI vars.
func (dv *DeviceView) uefiVariables() (json.RawMessage, error) {
	var varString string

	varString, ok := dv.Inv.Metadata[uefiVarsKey]
	if !ok {
		return nil, nil
	}

	// sanity check the data coming in from the caller
	m := map[string]any{}
	if err := json.Unmarshal([]byte(varString), &m); err != nil {
		return nil, errors.Wrap(err, "unmarshaling uefi-variables")
	}

	return []byte(varString), nil
}*/

func (dv *DeviceView) updateVendorAttributes(ctx context.Context, exec boil.ContextExecutor) error {
	return updateAnyAttribute(ctx, exec, true, dv.DeviceID.String(), alloyVendorNamespace, dv.vendorAttributes())
}

// write all the versioned-attributes from this server
func (dv *DeviceView) updateServerVAs(ctx context.Context, exec boil.ContextExecutor) error {
	statusData, _ := json.Marshal(dv.Inv.Status)
	// XXX: more VAs here
	return updateAnyVersionedAttribute(ctx, exec, true,
		dv.DeviceID.String(), serverStatusNamespace, statusData)
}

func (dv *DeviceView) UpsertInventory(ctx context.Context, exec boil.ContextExecutor) error {
	// yes, this is a dopey, repetitive style that should be easy for folks to extend or modify
	if err := dv.updateVendorAttributes(ctx, exec); err != nil {
		return errors.Wrap(err, "server vendor attributes update")
	}

	if err := dv.updateServerVAs(ctx, exec); err != nil {
		return errors.Wrap(err, "server versioned attribute update")
	}
	for _, cmp := range dv.Inv.Components {
		if err := composeRecords(ctx, exec, cmp, dv.Inband, dv.DeviceID.String()); err != nil {
			return err // we already wrapped some contextual data around the error
		}
	}
	return nil
}

func (dv *DeviceView) FromDatastore(ctx context.Context, exec boil.ContextExecutor) error {
	dv.Inv = &rivets.Server{
		ID:   dv.DeviceID.String(), // XXX: remove later?
		Name: dv.DeviceID.String(),
	}

	attrs, err := models.Attributes(
		models.AttributeWhere.ServerID.EQ(null.StringFrom(dv.DeviceID.String())),
	).All(ctx, exec)
	switch err {
	case nil:
	default:
		metrics.DBError("fetching attributes")
		return errors.Wrap(err, "fetching attributes")
	}

	if len(attrs) == 0 {
		return ErrNoInventory
	}

	for _, a := range attrs {
		switch a.Namespace {
		case alloyVendorNamespace:
			m := map[string]string{}
			if err := a.Data.Unmarshal(&m); err != nil {
				return errors.Wrap(err, "unmarshaling vendor attributes")
			}
			dv.Inv.Vendor = m[vendorKey]
			dv.Inv.Model = m[modelKey]
			dv.Inv.Serial = m[serialKey]
		default:
		}
	}

	statusVAttr, err := models.VersionedAttributes(
		models.VersionedAttributeWhere.ServerID.EQ(null.StringFrom(dv.DeviceID.String())),
		models.VersionedAttributeWhere.Namespace.EQ(serverStatusNamespace),
		qm.OrderBy("tally DESC"),
	).One(ctx, exec)

	switch err {
	case nil:
		var status string
		if err = json.Unmarshal(statusVAttr.Data, &status); err != nil {
			return errors.Wrap(err, "unmarshaling status attribute")
		}
		dv.Inv.Status = status
	case sql.ErrNoRows:
		// just skip it, status is optional
	default:
		metrics.DBError("fetch status VA")
		return errors.Wrap(err, "fetching versioned attibutes")
	}

	comps, err := componentsFromDatabase(ctx, exec, dv.Inband, dv.DeviceID.String())
	if err != nil {
		return errors.Wrap(err, "fetching components")
	}

	dv.Inv.Components = comps

	return nil
}
