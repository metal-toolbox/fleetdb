// Package inventory provides utilities to manage inventory
// within Attributes and Versioned Attributes
package inventory

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

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
