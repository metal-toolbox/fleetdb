package fleetdbapi_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

var BiosConfigSetTest fleetdbapi.BiosConfigSet = fleetdbapi.BiosConfigSet{
	ID:      uuid.NewString(),
	Name:    "Test",
	Version: "version",
	Components: []fleetdbapi.BiosConfigComponent{
		{
			Name:   "SM Motherboard",
			Vendor: "SUPERMICRO",
			Serial: "BIOS",
			Model:  "ATX",
			Settings: []fleetdbapi.BiosConfigSetting{
				{
					Key:   "BootOrder",
					Value: "dev2,dev3,dev4",
				},
				{
					Key:   "Mode",
					Value: "UEFI",
				},
			},
		},
		{
			Name:   "Intel Network Adapter",
			Vendor: "Intel",
			Serial: "NIC",
			Model:  "PCIE",
			Settings: []fleetdbapi.BiosConfigSetting{
				{
					Key:    "PXEEnable",
					Value:  "true",
					Raw: []byte(`{}`),
				},
				{
					Key:   "SRIOVEnable",
					Value: "false",
				},
				{
					Key:    "position",
					Value:  "1",
					Raw: []byte(`{ "lanes": 8 }`),
				},
			},
		},
	},
}

func TestBiosConfigSetCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource created", Slug: BiosConfigSetTest.ID})
		require.NoError(t, err)

		c := mockClient(string(jsonResponse), respCode)
		resp, err := c.CreateServerBiosConfigSet(ctx, BiosConfigSetTest)

		if expectError {
			assert.Error(t, err)

			return err
		}

		require.NotNil(t, resp) // stop testing if resp is nil
		assert.Equal(t, BiosConfigSetTest.ID, resp.Slug)

		return err
	})
}

func TestBiosConfigSetDelete(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, _ bool) error {
		jsonResponse := json.RawMessage([]byte(`{"message": "resource deleted"}`))
		c := mockClient(string(jsonResponse), respCode)
		_, err := c.DeleteServerBiosConfigSet(ctx, uuid.New())

		return err
	})
}

func TestBiosConfigSetGet(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Record: &BiosConfigSetTest})
		assert.NoError(t, err)

		c := mockClient(string(jsonResponse), respCode)

		resp, err := c.GetServerBiosConfigSet(ctx, uuid.New())
		if expectError {
			assert.Error(t, err)

			return err
		}

		assert.NoError(t, err)
		require.NotNil(t, resp) // stop testing if set is nil
		assert.Equal(t, BiosConfigSetTest, *resp.Record.(*fleetdbapi.BiosConfigSet))

		return err
	})
}

func TestConfgSetList(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Records: &[]fleetdbapi.BiosConfigSet{BiosConfigSetTest}})
		assert.NoError(t, err)

		c := mockClient(string(jsonResponse), respCode)

		testBiosConfigSetQueryParams := fleetdbapi.BiosConfigSetListParams{
			Params: []fleetdbapi.BiosConfigSetQueryParams{
				{
					Set: fleetdbapi.BiosConfigSetQuery{
						Components: []fleetdbapi.BiosConfigComponentQuery{
							{
								Name: "RTX",
							},
						},
					},
					LogicalOperator:    fleetdbapi.OperatorLogicalOR,
					ComparitorOperator: fleetdbapi.OperatorComparitorLike,
				},
				{
					Set: fleetdbapi.BiosConfigSetQuery{
						Components: []fleetdbapi.BiosConfigComponentQuery{
							{
								Settings: []fleetdbapi.BiosConfigSettingQuery{
									{
										Key:   "PCIE Lanes",
										Value: "x16",
									},
								},
							},
						},
					},
					LogicalOperator:    fleetdbapi.OperatorLogicalAND,
					ComparitorOperator: fleetdbapi.OperatorComparitorEqual,
				},
			},
			Pagination: fleetdbapi.PaginationParams{
				Limit:   50,
				Page:    4,
				Cursor:  "cursor",
				Preload: false,
				OrderBy: "nothing",
			},
		}

		resp, err := c.ListServerBiosConfigSet(ctx, &testBiosConfigSetQueryParams)
		if expectError {
			assert.Error(t, err)

			return err
		}

		assert.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, []fleetdbapi.BiosConfigSet{BiosConfigSetTest}, *resp.Records.(*[]fleetdbapi.BiosConfigSet))

		return err
	})
}

func TestConfgSetUpdate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource updated", Slug: BiosConfigSetTest.ID})
		require.NoError(t, err)

		id, err := uuid.Parse(BiosConfigSetTest.ID)
		require.NoError(t, err)

		c := mockClient(string(jsonResponse), respCode)
		resp, err := c.UpdateServerBiosConfigSet(ctx, id, BiosConfigSetTest)

		if expectError {
			assert.Error(t, err)

			return err
		}

		require.NotNil(t, resp) // stop testing if resp is nil
		assert.Equal(t, BiosConfigSetTest.ID, resp.Slug)

		return err
	})
}
