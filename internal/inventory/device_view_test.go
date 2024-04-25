//go:build integration

package inventory

import (
	"context"
	"testing"

	"github.com/bmc-toolbox/common"
	"github.com/google/uuid"
	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/models"
	rivets "github.com/metal-toolbox/rivets/types"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/stretchr/testify/require"
)

func Test_ServerSanityCheck(t *testing.T) {
	srv := &rivets.Server{}

	require.Error(t, ServerSanityCheck(srv), "zero value")
	srv.Model = "some model"
	require.Error(t, ServerSanityCheck(srv), "model only")
	srv.Vendor = "some vendor"
	require.Error(t, ServerSanityCheck(srv), "model+vendor")
	srv.Serial = "some-serial"
	// theoretically a server could have no components, from the FleetDB perspective
	require.NoError(t, ServerSanityCheck(srv))
	srv.Components = []*rivets.Component{
		&rivets.Component{},
	}
	require.Error(t, ServerSanityCheck(srv), "zero value component")
	srv.Components = []*rivets.Component{
		&rivets.Component{Name: common.SlugPhysicalMem},
	}
	require.Error(t, ServerSanityCheck(srv), "component name only")
	srv.Components = []*rivets.Component{
		&rivets.Component{
			Name:   common.SlugPhysicalMem,
			Serial: "some-serial-number",
		},
	}
	require.NoError(t, ServerSanityCheck(srv))
}

func Test_DeviceViewUpdate(t *testing.T) {
	t.Run("insert and update device attributes", func(t *testing.T) {
		db := dbtools.DatabaseTest(t)
		srvID := uuid.New()
		dv := DeviceView{
			Inv: &rivets.Server{
				Vendor: "CoolVendor",
				Model:  "BestModel 420",
				Serial: "0xdeadbeef",
				Status: "some status string",
				Components: []*rivets.Component{
					&rivets.Component{
						Name:   common.SlugBIOS,
						Serial: "some-serial",
						Firmware: &common.Firmware{
							Installed: "version",
						},
						Status: &common.Status{
							State:  "great",
							Health: "very yes",
						},
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
		dv.Inv.Status = "different status"
		err = dv.UpsertInventory(context.TODO(), db)
		require.NoError(t, err)

		// validate the contents
		read := DeviceView{
			DeviceID: srvID,
		}
		err = read.FromDatastore(context.TODO(), db)
		require.NoError(t, err)
		require.Equal(t, "different status", read.Inv.Status)

		require.Len(t, read.Inv.Components, 1)
		bios := read.Inv.Components[0]
		require.NotNil(t, bios.Firmware)
		require.NotNil(t, bios.Status)
	})
}
