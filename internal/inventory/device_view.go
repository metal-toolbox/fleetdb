//nolint:all  // XXX remove this!
package inventory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/bmc-toolbox/common"
	"github.com/google/uuid"
	"github.com/metal-toolbox/fleetdb/internal/models"
	"github.com/pkg/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
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
	alloyVendorNamespace   = "sh.hollow.alloy.server_vendor_attributes"
	alloyMetadataNamespace = "sh.hollow.alloy.server_metadata_attributes"
	alloyUefiVarsNamespace = "sh.hollow.alloy.server_uefi_variables" // this is a versioned attribute, we expect it to change

	// metadata keys
	modelKey    = "model"
	vendorKey   = "vendor"
	serialKey   = "serial"
	uefiVarsKey = "uefi-variables"
)

// A reminder for maintenance: this type needs to be able to contain all the
// relevant fields from Component-Inventory or Alloy.
type DeviceView struct {
	Inv        *common.Device    `json:"inventory"`
	BiosConfig map[string]string `json:"bios_config,omitempty"`
	Inband     bool              // the method of inventory collection
	DeviceID   uuid.UUID
}

func (dv *DeviceView) vendorAttributes() json.RawMessage {
	m := map[string]string{
		modelKey:  "unknown",
		serialKey: "unknown",
		vendorKey: "unknown",
	}

	if dv.Inv.Model != "" {
		m[modelKey] = dv.Inv.Model
	}

	if dv.Inv.Serial != "" {
		m[serialKey] = dv.Inv.Serial
	}

	if dv.Inv.Vendor != "" {
		m[vendorKey] = dv.Inv.Vendor
	}

	byt, _ := json.Marshal(m)

	return byt
}

func (dv *DeviceView) metadataAttributes() json.RawMessage {
	m := map[string]string{}

	// filter UEFI variables -- they go in a versioned-attribute
	for k, v := range dv.Inv.Metadata {
		if k != uefiVarsKey {
			m[k] = v
		}
	}

	if len(m) == 0 {
		return nil
	}

	byt, _ := json.Marshal(m)
	return byt
}

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
}

// the "either server or server-component" facet of attributes makes this function a
// little complicated
func updateAnyAttribute(ctx context.Context, exec boil.ContextExecutor,
	isServerAttr bool, id, namespace string, data json.RawMessage) error {
	var mods []qm.QueryMod

	idStr := null.StringFrom(id)
	attrData := types.JSON(data)

	// create an attribute in the event we need to make an insert
	attr := models.Attribute{
		Namespace: namespace,
		Data:      attrData,
	}

	if isServerAttr {
		attr.ServerID = idStr
		mods = append(mods, models.AttributeWhere.ServerID.EQ(idStr))
	} else {
		attr.ServerComponentID = idStr
		mods = append(mods, models.AttributeWhere.ServerComponentID.EQ(idStr))
	}
	mods = append(mods, models.AttributeWhere.Namespace.EQ(namespace))

	existing, err := models.Attributes(mods...).One(ctx, exec)
	switch err {
	case nil:
		attr.ID = existing.ID
		_, updErr := attr.Update(ctx, exec, boil.Infer())
		return updErr
	case sql.ErrNoRows:
		return attr.Insert(ctx, exec, boil.Infer())
	default:
		return err
	}
}

func (dv *DeviceView) updateVendorAttributes(ctx context.Context, exec boil.ContextExecutor) error {
	return updateAnyAttribute(ctx, exec, true, dv.DeviceID.String(), alloyVendorNamespace, dv.vendorAttributes())
}

func (dv *DeviceView) updateMetadataAttributes(ctx context.Context, exec boil.ContextExecutor) error {
	var err error
	if md := dv.metadataAttributes(); md != nil {
		err = updateAnyAttribute(ctx, exec, true, dv.DeviceID.String(), alloyMetadataNamespace, md)
	}
	return err
}

// insert a new versioned attribute record with the provided data. if this is not the first
// time we've seen a id/namespace tuple, increment the tally
func updateAnyVersionedAttribute(ctx context.Context, exec boil.ContextExecutor,
	isServerAttr bool, id, namespace string, data json.RawMessage) error {
	var mods []qm.QueryMod

	idStr := null.StringFrom(id)
	attrData := types.JSON(data)

	// we will always insert a new versioned attribute, just incrementing the tally
	vattr := models.VersionedAttribute{
		Namespace: namespace,
		Data:      attrData,
	}

	if isServerAttr {
		vattr.ServerID = idStr
		mods = append(mods, models.VersionedAttributeWhere.ServerID.EQ(idStr))
	} else {
		vattr.ServerComponentID = idStr
		mods = append(mods, models.VersionedAttributeWhere.ServerComponentID.EQ(idStr))
	}
	mods = append(mods, models.VersionedAttributeWhere.Namespace.EQ(namespace), qm.OrderBy("tally DESC"))

	lastVA, err := models.VersionedAttributes(mods...).One(ctx, exec)
	switch err {
	case nil:
		vattr.Tally = lastVA.Tally + 1
	case sql.ErrNoRows:
		// first time we've seen this vattr
	default:
		return err
	}
	return vattr.Insert(ctx, exec, boil.Infer())
}

// write a versioned-attribute containing the UEFI variables from this server
func (dv *DeviceView) updateUefiVariables(ctx context.Context, exec boil.ContextExecutor) error {
	uefiVarData, err := dv.uefiVariables()
	if err != nil {
		return err
	}
	return updateAnyVersionedAttribute(ctx, exec, true,
		dv.DeviceID.String(), alloyUefiVarsNamespace, uefiVarData)
}

func (dv *DeviceView) UpsertInventory(ctx context.Context, exec boil.ContextExecutor) error {
	// yes, this is a dopey, repetitive style that should be easy for folks to extend or modify
	if err := dv.updateVendorAttributes(ctx, exec); err != nil {
		return errors.Wrap(err, "vendor attributes update")
	}
	if err := dv.updateMetadataAttributes(ctx, exec); err != nil {
		return errors.Wrap(err, "metadata attribute update")
	}
	if err := dv.updateUefiVariables(ctx, exec); err != nil {
		return errors.Wrap(err, "uefi variables update")
	}
	return nil
}

func (dv *DeviceView) FromDatastore(ctx context.Context, exec boil.ContextExecutor) error {
	attrs, err := models.Attributes(qm.Where("server_id=?", dv.DeviceID)).All(ctx, exec)
	if err != nil {
		return err
	}

	if dv.Inv == nil {
		dv.Inv = &common.Device{
			Common: common.Common{
				Metadata: map[string]string{},
			},
		}
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
		case alloyMetadataNamespace:
			if err := a.Data.Unmarshal(&dv.Inv.Metadata); err != nil {
				return errors.Wrap(err, "unmarshaling metadata attributes")
			}
		default:
		}
	}

	uefiVarsAttr, err := models.VersionedAttributes(
		qm.Where("server_id=?", dv.DeviceID),
		qm.And(fmt.Sprintf("namespace='%s'", alloyUefiVarsNamespace)),
		qm.OrderBy("tally DESC"),
	).One(ctx, exec)

	if err != nil {
		return err
	}

	dv.Inv.Metadata[uefiVarsKey] = uefiVarsAttr.Data.String()

	// XXX: get components and component attributes and populate the dv.Inv

	return nil
}
