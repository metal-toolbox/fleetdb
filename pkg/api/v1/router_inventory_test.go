package fleetdbapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	common "github.com/metal-toolbox/bmc-common"
	rivets "github.com/metal-toolbox/rivets/v2/types"
	"github.com/stretchr/testify/require"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
)

func TestIntegrationGetInventory(t *testing.T) {
	s := serverTest(t)
	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		srvID := uuid.MustParse(dbtools.FixtureInventoryServer.ID)
		srv, _, err := s.Client.GetServerInventory(ctx, srvID, true)
		if !expectError {
			// we should get back the fixture data for the inventory server
			// we don't care about the name of the server
			require.NoError(t, err)
			require.NotNil(t, srv)
			require.Equal(t, "myModel", srv.Model)
			require.Equal(t, "Awesome Computer, Inc.", srv.Vendor)
			require.Equal(t, "1234xyz", srv.Serial)
		}
		return err
	})
}

func TestIntegrationSetInventory(t *testing.T) {
	s := serverTest(t)
	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)

		srv := &rivets.Server{
			Model:  "myModel",
			Vendor: "AwesomeCo",
			Serial: "1234xyz",
			Components: []*rivets.Component{
				{
					Name:   common.SlugBIOS,
					Serial: "0",
					Firmware: &common.Firmware{
						Installed: "best version",
					},
				},
			},
		}
		srvID := uuid.MustParse(dbtools.FixtureInventoryServer.ID)
		r, err := s.Client.SetServerInventory(ctx, srvID, srv, true)
		if !expectError {
			// we should get back the fixture data for the inventory server
			// we don't care about the name of the server
			require.NoError(t, err)
			require.NotNil(t, r)
		}
		return err
	})
}

func TestIntegrationCompareInventory(t *testing.T) {
	s := serverTest(t)
	realClientTests(t, func(ctx context.Context, authToken string, respCode int, expectError bool) error {
		s.Client.SetToken(authToken)
		tgtSrv := &rivets.Server{
			Model:  "fake_model",
			Vendor: "fake_vendor",
			Serial: "fake_serial",
			Components: []*rivets.Component{
				{
					Name:   common.SlugBIOS,
					Serial: "fake_bios_0",
					Firmware: &common.Firmware{
						Installed: "fake_bios_firmware_0",
					},
				},
				{
					Name:   common.SlugBMC,
					Serial: "fake_bmc_0",
					Firmware: &common.Firmware{
						Installed: "fake_bmc_firmware_0",
					},
				},
			},
		}
		srvID := uuid.MustParse(dbtools.FixtureInventoryServer.ID)
		r, err := s.Client.SetServerInventory(ctx, srvID, tgtSrv, true)
		if !expectError {
			// we should get back the fixture data for the inventory server
			// we don't care about the name of the server
			require.NoError(t, err)
			require.NotNil(t, r)
		}
		if err != nil {
			return err
		}
		srcSrv := &rivets.Server{
			Model:  "fake_model",
			Vendor: "fake_vendor",
			Serial: "fake_serial",
			Components: []*rivets.Component{
				{
					Name:   common.SlugBIOS,
					Serial: "fake_bios_1",
					Firmware: &common.Firmware{
						Installed: "fake_bios_firmware_1",
					},
				},
				{
					Name:   common.SlugCPU,
					Serial: "fake_cpu_0",
					Firmware: &common.Firmware{
						Installed: "fake_cpu_firmware_0",
					},
				},
			},
		}
		res, err := s.Client.CompareServerInventory(ctx, srvID, srcSrv, true)
		if err != nil {
			return err
		}
		expectedGapString :=
			`BMC: - Firmware: <nil>
+ Firmware: &{fake_bmc_firmware_0   [] map[]}
- Name: 
+ Name: BMC
- Serial: 
+ Serial: fake_bmc_0

BIOS: - Firmware: &{fake_bios_firmware_1   [] map[]}
+ Firmware: &{fake_bios_firmware_0   [] map[]}
- Serial: fake_bios_1
+ Serial: fake_bios_0

CPU: - Firmware: &{fake_cpu_firmware_0   [] map[]}
+ Firmware: <nil>
- Name: CPU
+ Name: 
- Serial: fake_cpu_0
+ Serial: 

`
		gaps := res.Record.(string)
		if gaps != expectedGapString {
			return fmt.Errorf("client.CompareServerInventory(..) expected gaps:\n%v, got\n%v", expectedGapString, gaps)
		}
		return nil
	})

}
