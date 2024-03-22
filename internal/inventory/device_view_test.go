//go:build integration

package inventory

import (
	"context"
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
	t.Parallel()
	t.Run("insert and update device attributes", func(t *testing.T) {
		t.Parallel()
		db := dbtools.DatabaseTest(t)
		dv := DeviceView{
			Inv: &common.Device{
				Common: common.Common{
					Vendor: "CoolVendor",
					Model:  "BestModel 420",
					Serial: "0xdeadbeef",
					Metadata: map[string]string{
						"uefi-variables": "shouldn't be here",
						"metakey":        "value",
					},
				},
			},
		}
		srvID := uuid.New()
		// make a server for out attributes
		server := models.Server{
			ID:           srvID.String(),
			Name:         null.StringFrom("dvtest-server"),
			FacilityCode: null.StringFrom("test1"),
		}

		srvErr := server.Insert(context.TODO(), db, boil.Infer())
		require.NoError(t, srvErr, "server setup failed")

		err := dv.UpsertInventory(context.TODO(), db, srvID, false)
		require.NoError(t, err)

		// do it again to test the update
		dv.Inv.Common.Serial = "roastbeef"
		err = dv.UpsertInventory(context.TODO(), db, srvID, false)
		require.NoError(t, err)

		// validate the contents
		read := DeviceView{}
		err = read.FromDatastore(context.TODO(), db, srvID)
		require.NoError(t, err)
		require.Equal(t, "roastbeef", read.Inv.Serial)
		require.Equal(t, 1, len(read.Inv.Metadata))
	})
}