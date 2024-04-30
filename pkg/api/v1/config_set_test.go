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

var configSetTest fleetdbapi.ConfigSet = fleetdbapi.ConfigSet{
	ID:      uuid.NewString(),
	Name:    "Test",
	Version: "version",
	Components: []fleetdbapi.ConfigComponent{
		{
			Name:   "SM Motherboard",
			Vendor: "SUPERMICRO",
			Serial: "BIOS",
			Model:  "ATX",
			Settings: []fleetdbapi.ConfigComponentSetting{
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
			Settings: []fleetdbapi.ConfigComponentSetting{
				{
					Key:    "PXEEnable",
					Value:  "true",
					Custom: []byte(`{}`),
				},
				{
					Key:   "SRIOVEnable",
					Value: "false",
				},
				{
					Key:    "position",
					Value:  "1",
					Custom: []byte(`{ "lanes": 8 }`),
				},
			},
		},
	},
}

func TestConfigSetCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource created", Slug: configSetTest.ID})
		require.NoError(t, err)

		c := mockClient(string(jsonResponse), respCode)
		resp, err := c.CreateServerConfigSet(ctx, configSetTest)

		if expectError {
			assert.Error(t, err)

			return err
		}

		require.NotNil(t, resp) // stop testing if resp is nil
		assert.Equal(t, configSetTest.ID, resp.Slug)

		return err
	})
}

func TestConfigSetDelete(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, _ bool) error {
		jsonResponse := json.RawMessage([]byte(`{"message": "resource deleted"}`))
		c := mockClient(string(jsonResponse), respCode)
		_, err := c.DeleteServerConfigSet(ctx, uuid.New())

		return err
	})
}

func TestConfigSetGet(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Record: &configSetTest})
		assert.NoError(t, err)

		c := mockClient(string(jsonResponse), respCode)

		resp, err := c.GetServerConfigSet(ctx, uuid.New())
		if expectError {
			assert.Error(t, err)

			return err
		}

		assert.NoError(t, err)
		require.NotNil(t, resp) // stop testing if set is nil
		assert.Equal(t, configSetTest, *resp.Record.(*fleetdbapi.ConfigSet))

		return err
	})
}

func TestConfgSetList(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Records: &[]fleetdbapi.ConfigSet{configSetTest}})
		assert.NoError(t, err)

		c := mockClient(string(jsonResponse), respCode)

		testConfigSetQueryParams := fleetdbapi.ConfigSetListParams{
			Params: []fleetdbapi.ConfigSetQueryParams{
				{
					Set: fleetdbapi.ConfigSetQuery{
						Components: []fleetdbapi.ConfigComponentQuery{
							{
								Name: "RTX",
							},
						},
					},
					LogicalOperator:    fleetdbapi.OperatorLogicalOR,
					ComparitorOperator: fleetdbapi.OperatorComparitorLike,
				},
				{
					Set: fleetdbapi.ConfigSetQuery{
						Components: []fleetdbapi.ConfigComponentQuery{
							{
								Settings: []fleetdbapi.ConfigComponentSettingQuery{
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

		resp, err := c.ListServerConfigSet(ctx, &testConfigSetQueryParams)
		if expectError {
			assert.Error(t, err)

			return err
		}

		assert.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, []fleetdbapi.ConfigSet{configSetTest}, *resp.Records.(*[]fleetdbapi.ConfigSet))

		return err
	})
}

func TestConfgSetUpdate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource updated", Slug: configSetTest.ID})
		require.NoError(t, err)

		id, err := uuid.Parse(configSetTest.ID)
		require.NoError(t, err)

		c := mockClient(string(jsonResponse), respCode)
		resp, err := c.UpdateServerConfigSet(ctx, id, configSetTest)

		if expectError {
			assert.Error(t, err)

			return err
		}

		require.NotNil(t, resp) // stop testing if resp is nil
		assert.Equal(t, configSetTest.ID, resp.Slug)

		return err
	})
}
