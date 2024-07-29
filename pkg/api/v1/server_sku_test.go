package fleetdbapi_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ServerSkuTest = fleetdbapi.ServerSku{
	Name:             "DreamMachine",
	Version:          "version",
	Vendor:           "AMD",
	Chassis:          "4U",
	BMCModel:         "1",
	MotherboardModel: "ATX",
	CPUVendor:        "AMD",
	CPUModel:         "7995WX",
	CPUCores:         96,
	CPUHertz:         2500000000,
	CPUCount:         1,
	AuxDevices: []fleetdbapi.AuxDevice{
		{
			Vendor:     "AMD",
			Model:      "W7900",
			DeviceType: "Dedicated GPU",
			Details:    []byte(`{"slot": 2,"stream-processors": 6144,"compute-units": 96,"tflops": 122.64,"memory-type": "GDDR6","memory-bytes": 48000000000,"TBP": 299}`),
		},
		{
			Vendor:     "AMD",
			Model:      "W7900",
			DeviceType: "Dedicated GPU",
			Details:    []byte(`{"slot": 3,"stream-processors": 6144,"compute-units": 96,"tflops": 122.64,"memory-type": "GDDR6","memory-bytes": 48000000000,"TBP": 299}`),
		},
	},
	Disks: []fleetdbapi.Disk{
		{
			Bytes:    8000000000000,
			Protocol: "SATA",
			Count:    1,
		},
		{
			Bytes:    4000000000000,
			Protocol: "NVME",
			Count:    1,
		},
	},
	Memory: []fleetdbapi.Memory{
		{
			Bytes: 8000000000,
			Count: 2,
		},
		{
			Bytes: 8000000000,
			Count: 2,
		},
	},
	Nics: []fleetdbapi.Nic{
		{
			PortBandwidth: 1000000000,
			PortCount:     2,
		},
		{
			PortBandwidth: 10000000000,
			PortCount:     1,
		},
	},
}

func TestServerSkuCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		id := uuid.NewString()
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource created", Slug: id})
		require.NoError(t, err)

		c := mockClient(string(jsonResponse), respCode)
		resp, err := c.CreateServerSku(ctx, ServerSkuTest)

		if expectError {
			assert.Error(t, err)

			return err
		}

		require.NotNil(t, resp) // stop testing if resp is nil
		assert.Equal(t, id, resp.Slug)

		return err
	})
}

func TestServerSkuDelete(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, _ bool) error {
		jsonResponse := json.RawMessage([]byte(`{"message": "resource deleted"}`))
		c := mockClient(string(jsonResponse), respCode)
		_, err := c.DeleteServerSku(ctx, uuid.New())

		return err
	})
}

func TestServerSkuGet(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Record: &ServerSkuTest})
		assert.NoError(t, err)

		c := mockClient(string(jsonResponse), respCode)

		resp, err := c.GetServerSku(ctx, uuid.New())
		if expectError {
			assert.Error(t, err)

			return err
		}

		assert.NoError(t, err)
		require.NotNil(t, resp) // stop testing if set is nil
		assert.Equal(t, ServerSkuTest, *resp.Record.(*fleetdbapi.ServerSku))

		return err
	})
}

func TestServerSkuList(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Records: &[]fleetdbapi.ServerSku{ServerSkuTest}})
		assert.NoError(t, err)

		c := mockClient(string(jsonResponse), respCode)

		testServerSkuQueryParams := fleetdbapi.ServerSkuListParams{
			Params: []fleetdbapi.ServerSkuQueryParams{
				{
					Sku: fleetdbapi.ServerSkuQuery{
						Name: "DreamMachine",
					},
					LogicalOperator:    fleetdbapi.OperatorLogicalOR,
					ComparitorOperator: fleetdbapi.OperatorComparitorLike,
				},
			},
			Pagination: fleetdbapi.PaginationParams{
				Limit:   50,
				Page:    4,
				Cursor:  "cursor",
				Preload: true,
				OrderBy: "nothing",
			},
		}

		resp, err := c.ListServerSku(ctx, &testServerSkuQueryParams)
		if expectError {
			assert.Error(t, err)

			return err
		}

		assert.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, []fleetdbapi.ServerSku{ServerSkuTest}, *resp.Records.(*[]fleetdbapi.ServerSku))

		return err
	})
}

func TestServerSkuUpdate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		idStr := uuid.NewString()
		jsonResponse, err := json.Marshal(fleetdbapi.ServerResponse{Message: "resource updated", Slug: idStr})
		require.NoError(t, err)

		id, _ := uuid.Parse(idStr)

		c := mockClient(string(jsonResponse), respCode)
		resp, err := c.UpdateServerSku(ctx, id, ServerSkuTest)

		if expectError {
			assert.Error(t, err)

			return err
		}

		require.NotNil(t, resp) // stop testing if resp is nil
		assert.Equal(t, idStr, resp.Slug)

		return err
	})
}
