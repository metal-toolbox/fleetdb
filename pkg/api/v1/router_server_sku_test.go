package fleetdbapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/models"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationServerSkuCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, _ bool) error {
		s.Client.SetToken(authToken)

		resp, err := s.Client.CreateServerSku(realClientTestCtx, ServerSkuTest)
		if err != nil {
			return err
		}

		_, err = uuid.Parse(resp.Slug)
		assert.NoError(t, err)

		return nil
	})

	var testCases = []struct {
		testName      string
		Name          string
		Version       string
		expectedError bool
	}{
		{
			"server sku: create; success",
			"DreamMachine2",
			"Version 1 million",
			false,
		},
		{
			"server sku: create; failure; invalide config",
			"",
			"",
			true,
		},
		{
			"server sku: create; duplicate sku",
			dbtools.FixtureServerSku.Name,
			dbtools.FixtureServerSku.Version,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ServerSkuTemp := ServerSkuTest
			ServerSkuTemp.Name = tc.Name
			ServerSkuTemp.Version = tc.Version

			resp, err := s.Client.CreateServerSku(context.TODO(), ServerSkuTemp)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}

		})
	}
}

func TestIntegrationServerSkuGet(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, _ bool) error {
		s.Client.SetToken(authToken)

		parsedID, err := uuid.Parse(dbtools.FixtureServerSku.ID)
		require.NoError(t, err)

		resp, err := s.Client.GetServerSku(realClientTestCtx, parsedID)
		if err != nil {
			return err
		}

		require.NotNil(t, resp)

		assertEntireServerSkuModelEqual(t,
			dbtools.FixtureServerSku,
			dbtools.FixtureServerSkuAuxDevices,
			dbtools.FixtureServerSkuDisks,
			dbtools.FixtureServerSkuMemory,
			dbtools.FixtureServerSkuNics,
			resp.Record.(*fleetdbapi.ServerSku))

		return nil
	})

	var testCases = []struct {
		testName      string
		id            string
		expectedError bool
	}{
		{
			"server sku: get; success",
			dbtools.FixtureServerSku.ID,
			false,
		},
		{
			"server sku: get; unknown sku",
			uuid.NewString(),
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			id, err := uuid.Parse(tc.id)
			require.NoError(t, err)

			resp, err := s.Client.GetServerSku(context.TODO(), id)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				assertEntireServerSkuModelEqual(t,
					dbtools.FixtureServerSku,
					dbtools.FixtureServerSkuAuxDevices,
					dbtools.FixtureServerSkuDisks,
					dbtools.FixtureServerSkuMemory,
					dbtools.FixtureServerSkuNics,
					resp.Record.(*fleetdbapi.ServerSku))
			}
		})
	}
}

func TestIntegrationServerSkuUpdate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, expectedError bool) error {
		s.Client.SetToken(authToken)

		ServerSkuTemp := ServerSkuTest
		var parsedID uuid.UUID
		var err error

		if expectedError {
			parsedID, err = uuid.NewUUID()
			require.NoError(t, err)
		} else {
			ServerSkuTemp.Name = "Integration Test Server Sku Update"
			ServerSkuTemp.Version = "Test Version"
			resp, err := s.Client.CreateServerSku(realClientTestCtx, ServerSkuTemp)
			require.NoError(t, err)
			require.NotNil(t, resp)

			parsedID, err = uuid.Parse(resp.Slug)
			require.NoError(t, err)

			resp, err = s.Client.GetServerSku(realClientTestCtx, parsedID)
			require.NoError(t, err)
			require.NotNil(t, resp)

			ServerSkuTemp = *resp.Record.(*fleetdbapi.ServerSku)
		}

		ServerSkuTemp.Version = "Test Version 2"
		ServerSkuTemp.AuxDevices[0].Vendor = "AMDX"
		ServerSkuTemp.Disks[0].Bytes = 50
		ServerSkuTemp.Memory[0].Bytes = 50
		ServerSkuTemp.Nics[0].PortCount = 99
		_, err = s.Client.UpdateServerSku(realClientTestCtx, parsedID, ServerSkuTemp)
		if err != nil {
			return err
		}

		if !expectedError {
			resp, err := s.Client.GetServerSku(realClientTestCtx, parsedID)
			require.NoError(t, err)
			require.NotNil(t, resp)

			sku := *resp.Record.(*fleetdbapi.ServerSku)

			assertServerSkuEqual(t, &ServerSkuTemp, &sku)
		}

		return nil
	})

	var testCases = []struct {
		testName      string
		id            string
		expectedError bool
	}{
		{
			"server sku: update; success",
			dbtools.FixtureServerSku.ID,
			false,
		},
		{
			"server sku: update; invalid uuid",
			uuid.NewString(),
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			ServerSkuTemp := fleetdbapi.ServerSku{}

			parsedID, err := uuid.Parse(tc.id)
			require.NoError(t, err)

			if !tc.expectedError {
				resp, err := s.Client.GetServerSku(context.TODO(), parsedID)
				require.NoError(t, err)
				require.NotNil(t, resp)

				ServerSkuTemp = *resp.Record.(*fleetdbapi.ServerSku)
				ServerSkuTemp.Version = "Test Version 2"
				ServerSkuTemp.AuxDevices[0].Vendor = "AMDX"
				ServerSkuTemp.Disks[0].Bytes = 50
				ServerSkuTemp.Memory[0].Bytes = 50
				ServerSkuTemp.Nics[0].PortCount = 99
			}

			resp, err := s.Client.UpdateServerSku(context.TODO(), parsedID, ServerSkuTemp)

			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				resp, err := s.Client.GetServerSku(context.TODO(), parsedID)
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				sku := *resp.Record.(*fleetdbapi.ServerSku)

				assertServerSkuEqual(t, &ServerSkuTemp, &sku)
			}
		})
	}
}

func TestIntegrationServerSkuDelete(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, expectedError bool) error {
		s.Client.SetToken(authToken)

		ServerSkuTemp := ServerSkuTest
		ServerSkuTemp.Name = "Integration Test Server Sku Delete"
		ServerSkuTemp.Version = "Test Version"
		resp, err := s.Client.CreateServerSku(realClientTestCtx, ServerSkuTemp)
		if expectedError && err != nil {
			return err
		}
		require.NoError(t, err)
		require.NotNil(t, resp)

		parsedID, err := uuid.Parse(resp.Slug)
		require.NoError(t, err)

		resp, err = s.Client.DeleteServerSku(realClientTestCtx, parsedID)
		if err != nil {
			return err
		}

		assert.Equal(t, "1", resp.Slug)

		resp, err = s.Client.GetServerSku(realClientTestCtx, parsedID)
		assert.Error(t, err)
		assert.Nil(t, resp)

		return nil
	})

	var testCases = []struct {
		testName      string
		expectedError bool
	}{
		{
			"server sku: delete; success",
			false,
		},
		{
			"server sku: delete; invalide uuid",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			var parsedID uuid.UUID
			var err error

			if tc.expectedError {
				parsedID, err = uuid.NewUUID()
				require.NoError(t, err)
			} else {
				ServerSkuTemp := ServerSkuTest
				ServerSkuTemp.Name = "Integration Test Server Sku Delete"
				ServerSkuTemp.Version = "Test Version"
				resp, err := s.Client.CreateServerSku(context.TODO(), ServerSkuTemp)
				require.NoError(t, err)
				require.NotNil(t, resp)

				parsedID, err = uuid.Parse(resp.Slug)
				require.NoError(t, err)
			}

			resp, err := s.Client.DeleteServerSku(context.TODO(), parsedID)
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)

				assert.Equal(t, "1", resp.Slug)

				resp, err = s.Client.GetServerSku(context.TODO(), parsedID)
				assert.Error(t, err)
				assert.Nil(t, resp)
			}
		})
	}
}

func TestIntegrationServerSkuList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, respCode int, expectedError bool) error {
		s.Client.SetToken(authToken)

		params := fleetdbapi.ServerSkuListParams{
			Params: []fleetdbapi.ServerSkuQueryParams{
				{
					Sku: fleetdbapi.ServerSkuQuery{
						Name: "DreamMachine",
					},
					LogicalOperator:    fleetdbapi.OperatorLogicalOR,
					ComparitorOperator: fleetdbapi.OperatorComparitorEqual,
				},
			},
		}

		_, err := s.Client.ListServerSku(realClientTestCtx, &params)
		if err != nil {
			return err
		}

		return nil
	})

	ServerSkuTempSetup := ServerSkuTest

	ServerSkuTempSetup.Name = "List Test 1"
	resp, err := s.Client.CreateServerSku(context.TODO(), ServerSkuTempSetup)
	require.NoError(t, err)
	require.NotNil(t, resp)

	ServerSkuTempSetup.Name = "List Test 2"
	resp, err = s.Client.CreateServerSku(context.TODO(), ServerSkuTempSetup)
	require.NoError(t, err)
	require.NotNil(t, resp)

	ServerSkuTempSetup.Name = "List Test 3"
	for i := range ServerSkuTempSetup.Disks { // remove NVME
		ServerSkuTempSetup.Disks[i].Protocol = "SATA"
		ServerSkuTempSetup.Disks[i].Bytes = 10
	}
	resp, err = s.Client.CreateServerSku(context.TODO(), ServerSkuTempSetup)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Get All 3
	paramTest0 := fleetdbapi.ServerSkuListParams{
		Params: []fleetdbapi.ServerSkuQueryParams{
			{
				Sku: fleetdbapi.ServerSkuQuery{
					Name: "List Test",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
		},
	}

	// Get all based on Disk Protocol
	paramTest1 := fleetdbapi.ServerSkuListParams{
		Params: []fleetdbapi.ServerSkuQueryParams{
			{
				Sku: fleetdbapi.ServerSkuQuery{
					Name: "List Test",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
			{
				Sku: fleetdbapi.ServerSkuQuery{
					Disks: []fleetdbapi.DiskQuery{
						{
							Protocol: "NVME",
						},
					},
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorEqual,
			},
		},
	}

	// Get all with a disk spaces less than 11
	paramTest2 := fleetdbapi.ServerSkuListParams{
		Params: []fleetdbapi.ServerSkuQueryParams{
			{
				Sku: fleetdbapi.ServerSkuQuery{
					Name: "List Test",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
			{
				Sku: fleetdbapi.ServerSkuQuery{
					Disks: []fleetdbapi.DiskQuery{
						{
							Bytes: []int64{11},
						},
					},
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLessThan,
			},
		},
	}

	var testCases = []struct {
		testName      string
		params        fleetdbapi.ServerSkuListParams
		expectedCount int
	}{
		{
			"server sku: list; get all three",
			paramTest0,
			3,
		},
		{
			"server sku: list; get all based on disk protocol",
			paramTest1,
			2,
		},
		{
			"server sku: list; get all based on disk space",
			paramTest2,
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			resp, err := s.Client.ListServerSku(context.TODO(), &tc.params)

			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, tc.expectedCount, len(*resp.Records.(*[]fleetdbapi.ServerSku)))
		})
	}
}

func assertServerSkuEqual(t *testing.T, expected *fleetdbapi.ServerSku, actual *fleetdbapi.ServerSku) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Version, actual.Version)
	assert.Equal(t, expected.Vendor, actual.Vendor)
	assert.Equal(t, expected.Chassis, actual.Chassis)
	assert.Equal(t, expected.BMCModel, actual.BMCModel)
	assert.Equal(t, expected.MotherboardModel, actual.MotherboardModel)
	assert.Equal(t, expected.CPUVendor, actual.CPUVendor)
	assert.Equal(t, expected.CPUModel, actual.CPUModel)
	assert.Equal(t, expected.CPUCores, actual.CPUCores)
	assert.Equal(t, expected.CPUHertz, actual.CPUHertz)
	assert.Equal(t, expected.CPUCount, actual.CPUCount)
	assert.ElementsMatch(t, expected.AuxDevices, actual.AuxDevices)
	assert.ElementsMatch(t, expected.Disks, actual.Disks)
	assert.ElementsMatch(t, expected.Memory, actual.Memory)
	assert.ElementsMatch(t, expected.Nics, actual.Nics)
}

func assertEntireServerSkuModelEqual(t *testing.T,
	expectedServerSku *models.ServerSku,
	expectedAuxDevices []*models.ServerSkuAuxDevice,
	expectedDisks []*models.ServerSkuDisk,
	expectedMemories []*models.ServerSkuMemory,
	expectedNics []*models.ServerSkuNic,
	actual *fleetdbapi.ServerSku,
) {
	assertServerSkuModelEqual(t, expectedServerSku, actual)

	require.Equal(t, len(expectedAuxDevices), len(actual.AuxDevices))
	require.Equal(t, len(expectedDisks), len(actual.Disks))
	require.Equal(t, len(expectedMemories), len(actual.Memory))
	require.Equal(t, len(expectedNics), len(actual.Nics))

	for _, expectedAuxDevice := range expectedAuxDevices {
		foundAuxDevice := false

		for _, actualAuxDevice := range actual.AuxDevices {
			if checkAuxDeviceModelEqual(t, expectedAuxDevice, &actualAuxDevice) {
				assertServerSkuAuxDeviceModelEqual(t, expectedAuxDevice, &actualAuxDevice)

				foundAuxDevice = true
				break
			}
		}

		if !foundAuxDevice {
			assert.Fail(t, fmt.Sprintf("expected to find disk: `%+v` in list: %+v", expectedAuxDevice, actual.AuxDevices))
		}
	}

	for _, expectedDisk := range expectedDisks {
		foundDisk := false

		for _, actualDisk := range actual.Disks {
			if checkDiskModelEqual(expectedDisk, &actualDisk) {
				assertServerSkuDiskModelEqual(t, expectedDisk, &actualDisk)

				foundDisk = true
				break
			}
		}

		if !foundDisk {
			assert.Fail(t, fmt.Sprintf("expected to find disk: `%+v` in list: %+v", expectedDisk, actual.Disks))
		}
	}

	for _, expectedMemory := range expectedMemories {
		foundMemory := false

		for _, actualMemory := range actual.Memory {
			if checkMemoryModelEqual(expectedMemory, &actualMemory) {
				assertServerSkuMemoryModelEqual(t, expectedMemory, &actualMemory)

				foundMemory = true
				break
			}
		}

		if !foundMemory {
			assert.Fail(t, fmt.Sprintf("expected to find memory: `%+v` in list: %+v", expectedMemory, actual.Memory))
		}
	}

	for _, expectedNic := range expectedNics {
		foundNic := false

		for _, actualNic := range actual.Nics {
			if checkNicModelEqual(expectedNic, &actualNic) {
				assertServerSkuNicModelEqual(t, expectedNic, &actualNic)

				foundNic = true
				break
			}
		}

		if !foundNic {
			assert.Fail(t, fmt.Sprintf("expected to find nic: `%+v` in list: %+v", expectedNic, actual.Nics))
		}
	}
}

func assertServerSkuModelEqual(t *testing.T,
	expected *models.ServerSku,
	actual *fleetdbapi.ServerSku,
) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Version, actual.Version)
	assert.Equal(t, expected.Chassis, actual.Chassis)
	assert.Equal(t, expected.BMCModel, actual.BMCModel)
	assert.Equal(t, expected.MotherboardModel, actual.MotherboardModel)
	assert.Equal(t, expected.CPUVendor, actual.CPUVendor)
	assert.Equal(t, expected.CPUModel, actual.CPUModel)
	assert.Equal(t, expected.CPUHertz, actual.CPUHertz)
	assert.Equal(t, expected.CPUCount, actual.CPUCount)
}

func assertServerSkuAuxDeviceModelEqual(t *testing.T,
	expected *models.ServerSkuAuxDevice,
	actual *fleetdbapi.AuxDevice,
) {
	assert.Equal(t, expected.Vendor, actual.Vendor)
	assert.Equal(t, expected.Model, actual.Model)
	assert.Equal(t, expected.DeviceType, actual.DeviceType)
	assert.JSONEq(t, string(expected.Details), string(actual.Details))
}

func assertServerSkuDiskModelEqual(t *testing.T,
	expected *models.ServerSkuDisk,
	actual *fleetdbapi.Disk,
) {
	assert.Equal(t, expected.Vendor, actual.Vendor)
	assert.Equal(t, expected.Model, actual.Model)
	assert.Equal(t, expected.Bytes, actual.Bytes)
	assert.Equal(t, expected.Protocol, actual.Protocol)
	assert.Equal(t, expected.Count, actual.Count)
}

func assertServerSkuMemoryModelEqual(t *testing.T,
	expected *models.ServerSkuMemory,
	actual *fleetdbapi.Memory,
) {
	assert.Equal(t, expected.Vendor, actual.Vendor)
	assert.Equal(t, expected.Model, actual.Model)
	assert.Equal(t, expected.Bytes, actual.Bytes)
	assert.Equal(t, expected.Count, actual.Count)
}

func assertServerSkuNicModelEqual(t *testing.T,
	expected *models.ServerSkuNic,
	actual *fleetdbapi.Nic,
) {
	assert.Equal(t, expected.Vendor, actual.Vendor)
	assert.Equal(t, expected.Model, actual.Model)
	assert.Equal(t, expected.PortBandwidth, actual.PortBandwidth)
	assert.Equal(t, expected.PortCount, actual.PortCount)
	assert.Equal(t, expected.Count, actual.Count)
}

func checkAuxDeviceModelEqual(t *testing.T, expected *models.ServerSkuAuxDevice, actual *fleetdbapi.AuxDevice) bool {
	j1, j2 := map[string]interface{}{}, map[string]interface{}{}

	err := json.Unmarshal(expected.Details, &j1)
	require.NoError(t, err)

	err = json.Unmarshal(actual.Details, &j2)
	require.NoError(t, err)

	if !reflect.DeepEqual(j2, j1) {
		return false
	}

	return expected.Vendor == actual.Vendor &&
		expected.Model == actual.Model &&
		expected.DeviceType == actual.DeviceType
}

func checkDiskModelEqual(expected *models.ServerSkuDisk, actual *fleetdbapi.Disk) bool {
	return expected.Vendor == actual.Vendor &&
		expected.Model == actual.Model &&
		expected.Bytes == actual.Bytes &&
		expected.Protocol == actual.Protocol &&
		expected.Count == actual.Count
}

func checkMemoryModelEqual(expected *models.ServerSkuMemory, actual *fleetdbapi.Memory) bool {
	return expected.Vendor == actual.Vendor &&
		expected.Model == actual.Model &&
		expected.Bytes == actual.Bytes &&
		expected.Count == actual.Count
}

func checkNicModelEqual(expected *models.ServerSkuNic, actual *fleetdbapi.Nic) bool {
	return expected.Vendor == actual.Vendor &&
		expected.Model == actual.Model &&
		expected.PortBandwidth == actual.PortBandwidth &&
		expected.PortCount == actual.PortCount &&
		expected.Count == actual.Count
}
