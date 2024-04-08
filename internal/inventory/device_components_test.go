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

func TestComponents(t *testing.T) {
	db := dbtools.DatabaseTest(t)
	t.Run("writeBios", func(t *testing.T) {
		srvUUID := mustCreateServerRecord(t, db, "write-bios")

		// XXX: If the serial number changes, we will insert a component record instead of updating it
		orig := &common.BIOS{
			Common: common.Common{
				Oem:         true,
				Vendor:      "coolguy",
				Model:       "xxxx",
				ProductName: "the-product",
				Serial:      "some-serial",
				Firmware: &common.Firmware{
					Installed: "some-version",
				},
				Status: &common.Status{
					State:  "OK",
					Health: "decent",
				},
			},
			SizeBytes:     int64(1111),
			CapacityBytes: int64(2222),
		}

		update := &common.BIOS{
			Common: common.Common{
				Oem:         false,
				Vendor:      "OpenBios",
				ProductName: "super-cool",
				Serial:      "some-serial",
				Firmware: &common.Firmware{
					Installed: "installed-version",
				},
				Status: &common.Status{
					State:  "contented",
					Health: "healthy",
				},
			},
			SizeBytes:     int64(3333),
			CapacityBytes: int64(4444),
		}

		device := &DeviceView{
			DeviceID: srvUUID,
			Inv: &common.Device{
				BIOS: orig,
			},
		} // inband = false

		tx := db.MustBegin()
		err := device.writeBios(context.TODO(), tx)
		require.NoError(t, err)
		_ = tx.Commit()

		scc, err := models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(common.SlugBIOS)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), scc)

		scr, err := models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(common.SlugBIOS)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).One(context.TODO(), db)
		require.NoError(t, err)

		ac, err := models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.AttributeWhere.Namespace.EQ(outofbandComponentNamespace),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), ac)

		vac, err := models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(outofbandComponentNamespace),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), vac)

		// update the BIOS
		device.Inv.BIOS = update
		tx2 := db.MustBegin()
		err = device.writeBios(context.TODO(), tx2)
		require.NoError(t, err)
		err = tx2.Commit()
		require.NoError(t, err)

		scc, err = models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(common.SlugBIOS)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), scc)

		ac, err = models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.AttributeWhere.Namespace.EQ(outofbandComponentNamespace),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), ac)

		ar, err := models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.AttributeWhere.Namespace.EQ(outofbandComponentNamespace),
		).One(context.TODO(), db)
		require.NoError(t, err)
		// unpack the Data to validate the update
		var attr attributes
		err = json.Unmarshal([]byte(ar.Data), &attr)
		require.NoError(t, err)
		require.Equal(t, int64(4444), attr.CapacityBytes)
		require.Equal(t, int64(3333), attr.SizeBytes)

		vac, err = models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(outofbandComponentNamespace),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(2), vac)

		vrec, err := models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(outofbandComponentNamespace),
			qm.OrderBy("tally DESC"),
		).One(context.TODO(), db)
		require.NoError(t, err)

		vattr := versionedAttributes{}
		err = json.Unmarshal([]byte(vrec.Data), &vattr)
		require.NoError(t, err)
		require.Equal(t, "contented", vattr.Status.State)
	})
	t.Run("writeDimms", func(t *testing.T) {
		srvUUID := mustCreateServerRecord(t, db, "write-dimms")

		device := &DeviceView{
			DeviceID: srvUUID,
			Inv: &common.Device{
				Memory: []*common.Memory{
					&common.Memory{
						Common: common.Common{ // stutter much?
							Vendor:      "Elephant",
							ProductName: "ivory",
							Serial:      "my-serial-number",
							Firmware: &common.Firmware{
								Installed: "installed-version",
							},
							Status: &common.Status{
								State:  "carceral",
								Health: "meh",
							},
						},
						ID:           "first",
						Slot:         "DIMM.Socket.4",
						SizeBytes:    int64(1024),
						ClockSpeedHz: int64(60),
					},
					&common.Memory{
						ID: "bogus",
					},
				},
			},
			Inband: true,
		}

		tx := db.MustBegin()
		err := device.writeDimms(context.TODO(), tx)
		require.NoError(t, err)
		_ = tx.Commit()

		// interrogate the database to validate our transformation
		// we expect a server component record, an attributes record,
		// and a versioned attributes record
		scc, err := models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(common.SlugPhysicalMem)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), scc)

		scr, err := models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(common.SlugPhysicalMem)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).One(context.TODO(), db)
		require.NoError(t, err)

		ac, err := models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.AttributeWhere.Namespace.EQ(inbandComponentNamespace),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), ac)

		vac, err := models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(inbandComponentNamespace),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), vac)

		// update the DIMM record
		device.Inv.Memory = []*common.Memory{
			&common.Memory{
				Common: common.Common{
					Vendor:      "Elephant",
					ProductName: "ivory",
					Serial:      "my-serial-number",
					Firmware: &common.Firmware{
						Installed: "installed-version",
					},
					Status: &common.Status{
						State:  "contented",
						Health: "healthy",
					},
				},
				ID:           "first",
				Slot:         "DIMM.Socket.4",
				SizeBytes:    int64(4096),
				ClockSpeedHz: int64(120),
			},
		}

		tx = db.MustBegin()
		err = device.writeDimms(context.TODO(), tx)
		require.NoError(t, err)
		_ = tx.Commit()

		ac, err = models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.AttributeWhere.Namespace.EQ(inbandComponentNamespace),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), ac)

		ar, err := models.Attributes(
			models.AttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.AttributeWhere.Namespace.EQ(inbandComponentNamespace),
		).One(context.TODO(), db)
		require.NoError(t, err)
		// unpack the Data to validate the update
		var attr attributes
		err = json.Unmarshal([]byte(ar.Data), &attr)
		require.NoError(t, err)
		require.Equal(t, int64(120), attr.ClockSpeedHz)
		require.Equal(t, int64(4096), attr.SizeBytes)

		vac, err = models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(inbandComponentNamespace),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(2), vac)

		vrec, err := models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(inbandComponentNamespace),
			qm.OrderBy("tally DESC"),
		).One(context.TODO(), db)
		require.NoError(t, err)

		vattr := versionedAttributes{}
		err = json.Unmarshal([]byte(vrec.Data), &vattr)
		require.NoError(t, err)
		require.Equal(t, "contented", vattr.Status.State)
	})
}
