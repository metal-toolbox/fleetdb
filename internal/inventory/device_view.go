//nolint:all  // XXX remove this!
package inventory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

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
	versionedAttributesByServerID = "(namespace, created_at) IN (select namespace, max(created_at) from versioned_attributes where server_id=? group by namespace)"

	// historically these values were determined/set by alloy, even though they are
	// internal to the data storage layer, hence the names
	alloyVendorNamespace   = "sh.hollow.alloy.server_vendor_attributes"
	alloyMetadataNamespace = "sh.hollow.alloy.server_metadata_attributes"
)

// A reminder for maintenance: this type needs to be able to contain all the
// relevant fields from Component-Inventory or Alloy.
type DeviceView struct {
	Inv        *common.Device    `json:"inventory"`
	BiosConfig map[string]string `json:"bios_config,omitempty"`
}

func (dv *DeviceView) vendorAttributes() json.RawMessage {
	m := map[string]string{
		"model":  "unknown",
		"serial": "unknown",
		"vendor": "unknown",
	}

	if dv.Inv.Model != "" {
		m["model"] = dv.Inv.Model
	}

	if dv.Inv.Serial != "" {
		m["serial"] = dv.Inv.Serial
	}

	if dv.Inv.Vendor != "" {
		m["vendor"] = dv.Inv.Vendor
	}

	byt, _ := json.Marshal(m)

	return byt
}

func (dv *DeviceView) metadataAttributes() json.RawMessage {
	m := map[string]string{}

	// filter UEFI variables -- they go in a versioned-attribute
	for k, v := range dv.Inv.Metadata {
		if k != "uefi-variables" {
			m[k] = v
		}
	}

	if len(m) == 0 {
		return nil
	}

	byt, _ := json.Marshal(m)
	return byt
}

func (dv *DeviceView) updateAnyAttribute(ctx context.Context, exec boil.ContextExecutor,
	srv uuid.UUID, namespace string, data json.RawMessage) error {
	mods := []qm.QueryMod{
		qm.Where("server_id=?", srv),
		qm.And(fmt.Sprintf("namespace='%s'", namespace)),
	}
	now := time.Now()

	existing, err := models.Attributes(mods...).One(ctx, exec)
	switch err {
	case nil:
		// do update
		existing.Data = types.JSON(data)
		existing.UpdatedAt = null.TimeFrom(now)
		_, updErr := existing.Update(ctx, exec, boil.Infer())
		return updErr
	case sql.ErrNoRows:
		// do insert
		vendorAttr := models.Attribute{
			ServerID:  null.StringFrom(srv.String()),
			Namespace: namespace,
			Data:      types.JSON(data),
			CreatedAt: null.TimeFrom(now),
		}
		return vendorAttr.Insert(ctx, exec, boil.Infer())
	default:
		return err
	}
}

func (dv *DeviceView) updateVendorAttributes(ctx context.Context, exec boil.ContextExecutor, srv uuid.UUID) error {
	return dv.updateAnyAttribute(ctx, exec, srv, alloyVendorNamespace, dv.vendorAttributes())
}

func (dv *DeviceView) updateMetadataAttributes(ctx context.Context, exec boil.ContextExecutor, srv uuid.UUID) error {
	var err error
	if md := dv.metadataAttributes(); md != nil {
		err = dv.updateAnyAttribute(ctx, exec, srv, alloyMetadataNamespace, md)
	}
	return err
}

func (dv *DeviceView) UpsertInventory(ctx context.Context, exec boil.ContextExecutor, srv uuid.UUID, inband bool) error {
	// yes, this is a dopey repetitive style that should be easy for folks to extend or modify
	if err := dv.updateVendorAttributes(ctx, exec, srv); err != nil {
		return errors.Wrap(err, "vendor attributes update")
	}
	if err := dv.updateMetadataAttributes(ctx, exec, srv); err != nil {
		return errors.Wrap(err, "metadata attribute update")
	}
	return nil
}

func (dv *DeviceView) FromDatastore(ctx context.Context, exec boil.ContextExecutor, srv uuid.UUID) error {
	// populate the vendor attributes
	attrs, err := models.Attributes(qm.Where("server_id=?", srv)).All(ctx, exec)
	if err != nil {
		return err
	}

	if dv.Inv == nil {
		dv.Inv = &common.Device{}
	}

	for _, a := range attrs {
		switch a.Namespace {
		case alloyVendorNamespace:
			m := map[string]string{}
			if err := a.Data.Unmarshal(&m); err != nil {
				return errors.Wrap(err, "unmarshaling vendor attributes")
			}
			dv.Inv.Vendor = m["vendor"]
			dv.Inv.Model = m["model"]
			dv.Inv.Serial = m["serial"]
		case alloyMetadataNamespace:
			m := map[string]string{}
			if err := a.Data.Unmarshal(&m); err != nil {
				return errors.Wrap(err, "unmarshaling metadata attributes")
			}
			dv.Inv.Metadata = m
		default:
		}
	}
	return nil
}
