//go:build integration

package inventory

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/bmc-toolbox/common"
	"github.com/google/uuid"
	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/stretchr/testify/require"
)

func Test_DeviceViewUpdate(t *testing.T) {
	t.Run("insert and update device attributes", func(t *testing.T) {
		db := dbtools.DatabaseTest(t)
		srvID := uuid.New()
		dv := DeviceView{
			Inv: &common.Device{
				Common: common.Common{
					Vendor: "CoolVendor",
					Model:  "BestModel 420",
					Serial: "0xdeadbeef",
					Metadata: map[string]string{
						"uefi-variables": `{ "msg":"hi there" }`,
						"metakey":        "value",
					},
				},
			},
			DeviceID: srvID,
		}
		// make a server for out attributes
		server := models.Server{
			ID:           srvID.String(),
			Name:         null.StringFrom("dvtest-server"),
			FacilityCode: null.StringFrom("test1"),
		}

		srvErr := server.Insert(context.TODO(), db, boil.Infer())
		require.NoError(t, srvErr, "server setup failed")

		err := dv.UpsertInventory(context.TODO(), db)
		require.NoError(t, err)

		// do it again to test the update
		dv.Inv.Common.Serial = "roastbeef"
		dv.Inv.Metadata["uefi-variables"] = `{ "msg": "hi again" }`
		err = dv.UpsertInventory(context.TODO(), db)
		require.NoError(t, err)

		// there should be 2 records for the UEFI variables versioned attribute

		count, err := models.VersionedAttributes(
			models.VersionedAttributeWhere.ServerID.EQ(null.StringFrom(srvID.String())),
			models.VersionedAttributeWhere.Namespace.EQ(alloyUefiVarsNamespace),
		).Count(context.TODO(), db)

		require.NoError(t, err, "counting uefi versioned attributes")
		require.Equal(t, int64(2), count)

		// validate the contents
		read := DeviceView{
			DeviceID: srvID,
		}
		err = read.FromDatastore(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, "roastbeef", read.Inv.Serial)
		require.Equal(t, 2, len(read.Inv.Metadata))

		varStr, ok := read.Inv.Metadata["uefi-variables"]
		require.True(t, ok)

		uefiVar := map[string]any{}
		require.NoError(t, json.Unmarshal([]byte(varStr), &uefiVar))
		require.Equal(t, "hi again", uefiVar["msg"])
	})
}
