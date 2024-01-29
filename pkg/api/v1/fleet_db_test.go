package fleetdb_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	fleetDBApi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

func TestFleetDBCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		srv := fleetDBApi.Server{UUID: uuid.New(), FacilityCode: "Test1"}
		jsonResponse := json.RawMessage([]byte(`{"message": "resource created", "slug":"00000000-0000-0000-0000-000000001234"}`))

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.Create(ctx, srv)
		if !expectError {
			assert.Equal(t, "00000000-0000-0000-0000-000000001234", res.String())
		}

		return err
	})
}

func TestFleetDBDelete(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse := json.RawMessage([]byte(`{"message": "resource deleted"}`))
		c := mockClient(string(jsonResponse), respCode)
		_, err := c.Delete(ctx, fleetDBApi.Server{UUID: uuid.New()})

		return err
	})
}
func TestFleetDBGet(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		srv := fleetDBApi.Server{UUID: uuid.New(), FacilityCode: "Test1"}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Record: srv})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.Get(ctx, srv.UUID)
		if !expectError {
			assert.Equal(t, srv.UUID, res.UUID)
			assert.Equal(t, srv.FacilityCode, res.FacilityCode)
		}

		return err
	})
}

func TestFleetDBList(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		srv := []fleetDBApi.Server{{UUID: uuid.New(), FacilityCode: "Test1"}}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Records: srv})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.List(ctx, nil)
		if !expectError {
			assert.ElementsMatch(t, srv, res)
		}

		return err
	})
}

func TestFleetDBUpdate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Message: "resource updated"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		_, err = c.Update(ctx, uuid.UUID{}, fleetDBApi.Server{Name: "new-name"})

		return err
	})
}

func TestFleetDBCreateAttributes(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		attr := fleetDBApi.Attributes{Namespace: "unit-test", Data: json.RawMessage([]byte(`{"test":"unit"}`))}
		jsonResponse := json.RawMessage([]byte(`{"message": "resource created"}`))

		c := mockClient(string(jsonResponse), respCode)
		_, err := c.CreateAttributes(ctx, uuid.New(), attr)

		return err
	})
}
func TestFleetDBGetAttributes(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		attr := &fleetDBApi.Attributes{Namespace: "unit-test", Data: json.RawMessage([]byte(`{"test":"unit"}`))}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Record: attr})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.GetAttributes(ctx, uuid.UUID{}, "unit-test")
		if !expectError {
			assert.Equal(t, attr, res)
		}

		return err
	})
}

func TestFleetDBDeleteAttributes(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Message: "resource deleted"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		_, err = c.DeleteAttributes(ctx, uuid.UUID{}, "unit-test")

		return err
	})
}

func TestFleetDBListAttributes(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		attrs := []fleetDBApi.Attributes{{Namespace: "unit-test", Data: json.RawMessage([]byte(`{"test":"unit"}`))}}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Records: attrs})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.ListAttributes(ctx, uuid.UUID{}, nil)
		if !expectError {
			assert.ElementsMatch(t, attrs, res)
		}

		return err
	})
}

func TestFleetDBUpdateAttributes(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Message: "resource updated"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		_, err = c.UpdateAttributes(ctx, uuid.UUID{}, "unit-test", json.RawMessage([]byte(`{"test":"unit"}`)))

		return err
	})
}

func TestFleetDBComponentsGet(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		sc := []fleetDBApi.ServerComponent{{Name: "unit-test", Serial: "1234"}}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Records: sc})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.GetComponents(ctx, uuid.UUID{}, nil)
		if !expectError {
			assert.ElementsMatch(t, sc, res)
		}

		return err
	})
}

func TestFleetDBComponentsList(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		sc := []fleetDBApi.ServerComponent{{Name: "unit-test", Serial: "1234"}}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Records: sc})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.ListComponents(ctx, &fleetDBApi.ServerComponentListParams{Name: "unit-test", Serial: "1234"})
		if !expectError {
			assert.ElementsMatch(t, sc, res)
		}

		return err
	})
}

func TestFleetDBComponentsCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Message: "resource created"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, err := c.CreateComponents(ctx, uuid.New(), fleetDBApi.ServerComponentSlice{{Name: "unit-test"}})
		if !expectError {
			assert.Contains(t, res.Message, "resource created")
		}

		return err
	})
}

func TestFleetDBComponentsUpdate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Message: "resource updated"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, err := c.UpdateComponents(ctx, uuid.New(), fleetDBApi.ServerComponentSlice{{Name: "unit-test"}})
		if !expectError {
			assert.Contains(t, res.Message, "resource updated")
		}

		return err
	})
}

func TestFleetDBVersionedAttributeCreate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		va := fleetDBApi.VersionedAttributes{Namespace: "unit-test", Data: json.RawMessage([]byte(`{"test":"unit"}`))}
		jsonResponse := json.RawMessage([]byte(`{"message": "resource created", "slug":"the-namespace"}`))

		c := mockClient(string(jsonResponse), respCode)
		resp, err := c.CreateVersionedAttributes(ctx, uuid.New(), va)
		if !expectError {
			assert.Equal(t, "the-namespace", resp.Slug)
		}

		return err
	})
}

func TestFleetDBGetVersionedAttributess(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		va := []fleetDBApi.VersionedAttributes{{Namespace: "test", Data: json.RawMessage([]byte(`{}`))}}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Records: va})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.GetVersionedAttributes(ctx, uuid.New(), "namespace")
		if !expectError {
			assert.ElementsMatch(t, va, res)
		}

		return err
	})
}

func TestFleetDBListVersionedAttributess(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		va := []fleetDBApi.VersionedAttributes{{Namespace: "test", Data: json.RawMessage([]byte(`{}`))}}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Records: va})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.ListVersionedAttributes(ctx, uuid.New())
		if !expectError {
			assert.ElementsMatch(t, va, res)
		}

		return err
	})
}

func TestFleetDBCreateServerComponentFirmware(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		firmware := fleetDBApi.ComponentFirmwareVersion{
			UUID:    uuid.New(),
			Vendor:  "Dell",
			Model:   []string{"R615"},
			Version: "21.07.00",
		}
		jsonResponse := json.RawMessage([]byte(`{"message": "resource created", "slug":"00000000-0000-0000-0000-000000001234"}`))

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.CreateServerComponentFirmware(ctx, firmware)
		if !expectError {
			assert.Equal(t, "00000000-0000-0000-0000-000000001234", res.String())
		}

		return err
	})
}

func TestFleetDBServerComponentFirmwareDelete(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse := json.RawMessage([]byte(`{"message": "resource deleted"}`))
		c := mockClient(string(jsonResponse), respCode)
		_, err := c.DeleteServerComponentFirmware(ctx, fleetDBApi.ComponentFirmwareVersion{UUID: uuid.New()})

		return err
	})
}
func TestFleetDBServerComponentFirmwareGet(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		firmware := fleetDBApi.ComponentFirmwareVersion{
			UUID:    uuid.New(),
			Vendor:  "Dell",
			Model:   []string{"R615"},
			Version: "21.07.00",
		}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Record: firmware})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.GetServerComponentFirmware(ctx, firmware.UUID)
		if !expectError {
			assert.Equal(t, firmware.UUID, res.UUID)
			assert.Equal(t, firmware.Vendor, res.Vendor)
			assert.Equal(t, firmware.Model, res.Model)
			assert.Equal(t, firmware.Version, res.Version)
		}

		return err
	})
}

func TestFleetDBServerComponentFirmwareList(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		firmware := []fleetDBApi.ComponentFirmwareVersion{{
			UUID:    uuid.New(),
			Vendor:  "Dell",
			Model:   []string{"R615"},
			Version: "21.07.00",
		}}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Records: firmware})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, _, err := c.ListServerComponentFirmware(ctx, nil)
		if !expectError {
			assert.ElementsMatch(t, firmware, res)
		}

		return err
	})
}

func TestFleetDBServerComponentFirmwareUpdate(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Message: "resource updated"})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		_, err = c.UpdateServerComponentFirmware(ctx, uuid.UUID{}, fleetDBApi.ComponentFirmwareVersion{UUID: uuid.New()})

		return err
	})
}

func TestBillOfMaterialsBatchUpload(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		bom := []fleetDBApi.Bom{{SerialNum: "fakeSerialNum1", AocMacAddress: "fakeAocMacAddress1", BmcMacAddress: "fakeBmcMacAddress1"}}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Record: bom})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		res, err := c.BillOfMaterialsBatchUpload(ctx, bom)
		if !expectError {
			assert.Equal(t, []interface{}([]interface{}{
				map[string]interface{}{
					"aoc_mac_address": "fakeAocMacAddress1",
					"bmc_mac_address": "fakeBmcMacAddress1",
					"metro":           "",
					"num_def_pwd":     "",
					"num_defi_pmi":    "",
					"serial_num":      "fakeSerialNum1"}}), res.Record)
		}

		return err
	})
}

func TestGetBomInfoByAOCMacAddr(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		bom := fleetDBApi.Bom{SerialNum: "fakeSerialNum1", AocMacAddress: "fakeAocMacAddress1", BmcMacAddress: "fakeBmcMacAddress"}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Record: bom})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		respBom, _, err := c.GetBomInfoByAOCMacAddr(ctx, "fakeAocMacAddress1")
		if !expectError {
			assert.Equal(t, &bom, respBom)
		}

		return err
	})
}

func TestGetBomInfoByBMCMacAddr(t *testing.T) {
	mockClientTests(t, func(ctx context.Context, respCode int, expectError bool) error {
		bom := fleetDBApi.Bom{SerialNum: "fakeSerialNum1", AocMacAddress: "fakeAocMacAddress1", BmcMacAddress: "fakeBmcMacAddress1"}
		jsonResponse, err := json.Marshal(fleetDBApi.ServerResponse{Record: bom})
		require.Nil(t, err)

		c := mockClient(string(jsonResponse), respCode)
		respBom, _, err := c.GetBomInfoByBMCMacAddr(ctx, "fakeBmcMacAddress1")
		if !expectError {
			assert.Equal(t, &bom, respBom)
		}

		return err
	})
}
