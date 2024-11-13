package fleetdbapi_test

import (
	"context"
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
