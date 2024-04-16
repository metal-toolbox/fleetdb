//go:build integration

package inventory

import (
	"context"
	"database/sql"
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

func TestComposeComponentRecords(t *testing.T) {
	db := dbtools.DatabaseTest(t)
	t.Run("common case", func(t *testing.T) {
		srvUUID := mustCreateServerRecord(t, db, "common-case")

		// this can be any real slug, but *must* be a real slug, otherwise
		// we will panic on the slug -> type-id lookup.
		slug := common.SlugBIOS

		var inband bool
		attributeNS := getAttributeNamespace(inband)
		fwns := getFirmwareNamespace(inband)
		sns := getStatusNamespace(inband)

		orig := &common.Common{
			Oem:         true,
			Vendor:      "the-vendor",
			ProductName: "super-product",
			Firmware: &common.Firmware{
				Installed: "old-version",
			},
			Status: &common.Status{
				State:  "OK",
				Health: "decent",
			},
		}

		attr := &attributes{
			WWN: "string payload",
		}

		// we will add the product name to the attributes en-passant as part of
		// normalizing any missing model information
		expAttr := &attributes{
			WWN:         "string payload",
			ProductName: "super-product",
		}

		tx := db.MustBegin()
		err := composeRecords(context.TODO(), tx, orig, inband, srvUUID.String(), slug, attr)
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
		require.False(t, scr.Model.IsZero())
		require.Equal(t, "super-product", scr.Model.String)
		require.False(t, scr.Serial.IsZero())
		require.Equal(t, "0", scr.Serial.String)

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
		recordData := &attributes{}
		require.NoError(t, recordData.FromJSON(ar.Data))
		require.Equal(t, expAttr, recordData)

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
		// XXX: If the serial number changes, we will insert a component record instead of updating it
		update := &common.Common{
			Oem:         true,
			Vendor:      "the-vendor",
			ProductName: "super-new-product",
			Firmware: &common.Firmware{
				Installed: "new-version",
			},
			Status: &common.Status{
				State:  "happy",
				Health: "content",
			},
		}

		expAttr.ProductName = "super-new-product"

		tx = db.MustBegin()
		err = composeRecords(context.TODO(), tx, update, inband, srvUUID.String(), slug, attr)
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
		require.False(t, scr.Model.IsZero())
		require.Equal(t, "super-new-product", scr.Model.String)
		require.False(t, scr.Serial.IsZero())
		require.Equal(t, "0", scr.Serial.String)

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
		recordData = &attributes{}
		require.NoError(t, recordData.FromJSON(ar.Data))
		require.Equal(t, expAttr, recordData)

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
		comps, err := componentsFromDatabase(context.TODO(), db, inband, srvUUID.String(), slug)
		require.NoError(t, err)
		require.Len(t, comps, 1)
		require.NotNil(t, comps[0].cmn)
		require.NotNil(t, comps[0].attr)
	})
}

// We've tested the Common->component/attribute/versioned-attribute thing already
// so here we only check the component-specific attributes
func TestComponents(t *testing.T) {
	db := dbtools.DatabaseTest(t)
	t.Run("writeBios", func(t *testing.T) {
		srvUUID := mustCreateServerRecord(t, db, "write-bios")

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

		//attributeNS := getAttributeNamespace(device.Inband)

		tx := db.MustBegin()
		err := device.writeBios(context.TODO(), tx)
		require.NoError(t, err)
		_ = tx.Commit()

		reader := &DeviceView{
			DeviceID: srvUUID,
			Inv:      &common.Device{},
		}
		err = reader.getBios(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, orig.Firmware, reader.Inv.BIOS.Firmware)
		require.Equal(t, orig.Status, reader.Inv.BIOS.Status)
		require.Equal(t, orig.SizeBytes, reader.Inv.BIOS.SizeBytes)
		require.Equal(t, orig.CapacityBytes, reader.Inv.BIOS.CapacityBytes)

		// update the BIOS
		device.Inv.BIOS = update
		tx2 := db.MustBegin()
		err = device.writeBios(context.TODO(), tx2)
		require.NoError(t, err)
		err = tx2.Commit()
		require.NoError(t, err)

		err = reader.getBios(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, update.Firmware, reader.Inv.BIOS.Firmware)
		require.Equal(t, update.Status, reader.Inv.BIOS.Status)
		require.Equal(t, update.SizeBytes, reader.Inv.BIOS.SizeBytes)
		require.Equal(t, update.CapacityBytes, reader.Inv.BIOS.CapacityBytes)

		// failed lookup
		nogo := &DeviceView{
			DeviceID: srvUUID,
			Inv:      &common.Device{},
			Inband:   true,
		}

		err = nogo.getBios(context.TODO(), db)
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
	t.Run("writeBMC", func(t *testing.T) {
		// writeBMC is basically a re-run of composeRecords() so skip testing the update
		srvUUID := mustCreateServerRecord(t, db, "write-bmc")

		orig := &common.BMC{
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
			// BMC contains a NIC and ID string as well, but not populated or stored?
		}

		device := &DeviceView{
			DeviceID: srvUUID,
			Inv: &common.Device{
				BMC: orig,
			},
		} // inband = false

		tx := db.MustBegin()
		err := device.writeBMC(context.TODO(), tx)
		require.NoError(t, err)
		_ = tx.Commit()

		reader := &DeviceView{
			DeviceID: srvUUID,
			Inv:      &common.Device{},
		}
	})
	/*t.Run("writeMainboard", func(t *testing.T) {
		srvUUID := mustCreateServerRecord(t, db, "write-mainboard")

		device := &DeviceView{
			DeviceID: srvUUID,
			Inband:   true,
			Inv: &common.Device{
				Mainboard: &common.Mainboard{
					Common: common.Common{
						Oem:         true,
						Vendor:      "theVendor",
						Model:       "theBest",
						ProductName: "super-board",
						Description: "a mainboard",
						Firmware: &common.Firmware{
							Installed: "old-version",
						},
						Status: &common.Status{
							State:  "OK",
							Health: "meh",
						},
					},
				},
			},
		}

		expAttr := {
			ProductName: "super-board",
			Description: "a mainboard",
		}

		tx := db.MustBegin()
		err := device.writeMainboard(context.TODO(), tx)
		require.NoError(t, err)
		_ = tx.Commit()

		scc, err := models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(common.SlugMainboard)),
			models.ServerComponentWhere.ServerID.EQ(srvUUID.String()),
		).Count(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, int64(1), scc)

		scr, err := models.ServerComponents(
			models.ServerComponentWhere.Name.EQ(null.StringFrom(common.SlugMainboard)),
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

		device.Inv.Mainboard.Firmware = &common.Firmware{
			Installed: "new-version",
		}

		tx = db.MustBegin()
		err = device.writeMainboard(context.TODO(), tx)
		require.NoError(t, err)
		_ = tx.Commit()

		vrec, err := models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerComponentID.EQ(null.StringFrom(scr.ID)),
			models.VersionedAttributeWhere.Namespace.EQ(inbandComponentNamespace),
			qm.OrderBy("tally DESC"),
		).One(context.TODO(), db)
		require.NoError(t, err)

		vattr := versionedAttributes{}
		err = json.Unmarshal([]byte(vrec.Data), &vattr)
		require.NoError(t, err)
		require.Equal(t, "new-version", vattr.Firmware.Installed)
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
	})*/
}
