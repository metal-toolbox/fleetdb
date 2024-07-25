//go:build integration

package inventory

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/bmc-toolbox/common"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/models"
	rivets "github.com/metal-toolbox/rivets/types"
)

func mustCreateServerRecord(t *testing.T, db *sqlx.DB, name string) uuid.UUID {
	t.Helper()
	// we need to create a server in order to fulfill the foreign-key requirement
	// for server-components
	srv := models.Server{
		Name:         null.StringFrom(name),
		FacilityCode: null.StringFrom("tf2"),
	}

	err := srv.Insert(context.TODO(), db, boil.Infer())
	require.NoError(t, err, "server setup")
	srvUUID := uuid.MustParse(srv.ID)
	return srvUUID
}

func TestComposeComponentRecords(t *testing.T) {
	db := dbtools.DatabaseTest(t)
	t.Run("nil attr/vattr", func(t *testing.T) {
		srvUUID := mustCreateServerRecord(t, db, "nil-attr-case")

		var inband bool
		slug := common.SlugBIOS

		orig := &rivets.Component{
			Name:   slug,
			Vendor: "the-vendor",
		}

		attributeNS := getAttributeNamespace(inband)
		fwns := getFirmwareNamespace(inband)
		sns := getStatusNamespace(inband)

		tx := db.MustBegin()
		err := composeRecords(context.TODO(), tx, orig, inband, srvUUID.String())
		require.NoError(t, err)
		_ = tx.Commit()

		scc, err := models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(slug)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), scc)

		scr, err := models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(slug)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).One(context.TODO(), db)
		require.NoError(t, err)
		require.True(t, scr.Model.IsZero())

		rc, err := models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.AttributeWhere.Namespace.EQ(attributeNS),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(0), rc)

		rc, err = models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(fwns),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(0), rc)

		rc, err = models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(sns),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(0), rc)
	})
	t.Run("common case", func(t *testing.T) {
		srvUUID := mustCreateServerRecord(t, db, "common-case")

		var inband bool
		attributeNS := getAttributeNamespace(inband)
		fwns := getFirmwareNamespace(inband)
		sns := getStatusNamespace(inband)

		slug := common.SlugBIOS

		orig := &rivets.Component{
			// this can be any real slug, but *must* be a real slug, otherwise
			// we will panic on the slug -> type-id lookup.
			Name:   slug,
			Vendor: "the-vendor",
			Serial: "some-serial-number",
			Firmware: &common.Firmware{
				Installed: "old-version",
			},
			Status: &common.Status{
				State:  "OK",
				Health: "decent",
			},
			Attributes: &rivets.ComponentAttributes{
				PhysicalID: "my-id",
			},
		}

		tx := db.MustBegin()
		err := composeRecords(context.TODO(), tx, orig, inband, srvUUID.String())
		require.NoError(t, err)
		_ = tx.Commit()

		// interrogate the DB for our records
		scc, err := models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(slug)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), scc)

		scr, err := models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(slug)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).One(context.TODO(), db)
		require.NoError(t, err)
		require.True(t, scr.Model.IsZero())

		// validate the record counts of the attributes/versioned-attributes
		ac, err := models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.AttributeWhere.Namespace.EQ(attributeNS),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), ac, "attribute record")

		ar, err := models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.AttributeWhere.Namespace.EQ(attributeNS),
		).One(context.TODO(), db)
		require.NoError(t, err)

		vac, err := models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(fwns),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), vac, "firmware record")

		fwr, err := models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(fwns),
		).One(context.TODO(), db)
		require.NoError(t, err)
		fwData, err := firmwareFromJSON(fwr.Data)
		require.NoError(t, err)
		require.Equal(t, orig.Firmware, fwData)

		vac, err = models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(sns),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), vac, "status record")

		sr, err := models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(sns),
		).One(context.TODO(), db)
		require.NoError(t, err)
		srData, err := statusFromJSON(sr.Data)
		require.NoError(t, err)
		require.Equal(t, orig.Status, srData)

		// now do the update
		update := &rivets.Component{
			Name:   common.SlugBIOS,
			Vendor: "the-vendor",
			Serial: "some-serial-number",
			Firmware: &common.Firmware{
				Installed: "new-version",
			},
			Status: &common.Status{
				State:  "happy",
				Health: "content",
			},
			Attributes: &rivets.ComponentAttributes{
				PhysicalID: "my-id",
			},
		}

		tx = db.MustBegin()
		err = composeRecords(context.TODO(), tx, update, inband, srvUUID.String())
		require.NoError(t, err)
		_ = tx.Commit()

		scc, err = models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(slug)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), scc)

		scr, err = models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(slug)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).One(context.TODO(), db)
		require.NoError(t, err)

		// validate the record counts of the attributes/versioned-attributes
		ac, err = models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.AttributeWhere.Namespace.EQ(attributeNS),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), ac, "attribute record")

		ar, err = models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.AttributeWhere.Namespace.EQ(attributeNS),
		).One(context.TODO(), db)
		require.NoError(t, err)
		require.NotNil(t, ar.Data)
		attr, err := componentAttributesFromJSON(ar.Data)
		require.NoError(t, err)
		require.Equal(t, update.Attributes, attr)

		vac, err = models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(fwns),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(2), vac, "firmware record")

		fwr, err = models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(fwns),
			qm.OrderBy("tally DESC"),
		).One(context.TODO(), db)
		require.NoError(t, err)
		fwData, err = firmwareFromJSON(fwr.Data)
		require.NoError(t, err)
		require.Equal(t, update.Firmware, fwData)

		vac, err = models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(sns),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(2), vac, "status record")

		sr, err = models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(sns),
			qm.OrderBy("tally DESC"),
		).One(context.TODO(), db)
		require.NoError(t, err)
		srData, err = statusFromJSON(sr.Data)
		require.NoError(t, err)
		require.Equal(t, update.Status, srData)

		// validate that we can get all the component data we expect
		comps, err := componentsFromDatabase(context.TODO(), db, inband, srvUUID.String())
		require.NoError(t, err)
		require.Len(t, comps, 1)
	})
	t.Run("nonconforming existing data", func(t *testing.T) {
		// write an attribute record that doesn't have rivets.ComponentAttribute as a basis
		srvUUID := mustCreateServerRecord(t, db, "bad-attribute-json")

		var inband bool
		attributeNS := getAttributeNamespace(inband)

		slug := common.SlugBIOS

		orig := &rivets.Component{
			// this can be any real slug, but *must* be a real slug, otherwise
			// we will panic on the slug -> type-id lookup.
			Name:   slug,
			Vendor: "the-vendor",
			Serial: "some-serial-number",
			Firmware: &common.Firmware{
				Installed: "old-version",
			},
			Status: &common.Status{
				State:  "OK",
				Health: "decent",
			},
		}

		tx := db.MustBegin()
		err := composeRecords(context.TODO(), tx, orig, inband, srvUUID.String())
		require.NoError(t, err)
		_ = tx.Commit()

		// add the crazy attribute record
		compRecs, err := models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(slug)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).All(context.TODO(), db)

		require.NoError(t, err)
		require.Len(t, compRecs, 1)
		// get id and inject attributes
		compID := compRecs[0].ID
		badData := []map[string]string{
			{
				"msg": "this is not a rivets component attributes structure",
			},
			{
				"msg": "this is also not a rivets component attributes structure",
			},
		}
		payload, err := json.Marshal(badData)
		require.NoError(t, err)

		err = updateAnyAttribute(context.TODO(), db, false, compID, attributeNS, payload)
		require.NoError(t, err)

		// be pedantic and validate that there is an attribute record
		attrs, err := models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(compID)),
			models.AttributeWhere.Namespace.EQ(attributeNS),
		).All(context.TODO(), db)
		require.NoError(t, err)
		require.Len(t, attrs, 1)

		// now ask for this component via the API
		comps, err := componentsFromDatabase(context.TODO(), db, inband, srvUUID.String())
		require.NoError(t, err, "received error with type %T", err)
		require.Len(t, comps, 1)
	})
}
