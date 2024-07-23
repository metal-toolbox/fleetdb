package fleetdbapi_test

import (
	"context"

	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/models"

	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

func TestIntegrationServerBiosConfigSetCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, _ bool) error {
		s.Client.SetToken(authToken)

		BiosConfigSetTemp := BiosConfigSetTest
		BiosConfigSetTemp.ID = ""
		resp, err := s.Client.CreateServerBiosConfigSet(realClientTestCtx, BiosConfigSetTest)
		if err != nil {
			return err
		}

		_, err = uuid.Parse(resp.Slug)
		assert.NoError(t, err)

		BiosConfigSetTest.ID = resp.Slug

		return nil
	})

	var testCases = []struct {
		testName          string
		BiosConfigSetName string
		BiosConfigSetID   string
		expectedError     bool
		msgs              []string
	}{
		{
			"config set router: config set create; success",
			"Integration Test Config Set Success 1",
			uuid.NewString(),
			false,
			[]string{"resource created"},
		},
		{
			"config set router: config set create; success empty ID",
			"Integration Test Config Set Success 2",
			"",
			false,
			[]string{"resource created"},
		},
		{
			"config set router: config set create; invalid config",
			"",
			"",
			true,
			[]string{"invalid payload: BiosConfigSetCreate{}", "Field validation for 'Name' failed on the 'required' tag", "400"},
		},
		{
			"config set router: config set create; duplicate config",
			"Integration Test Config Set Fail 2",
			dbtools.FixtureBiosConfigSet.ID,
			true,
			[]string{fmt.Sprintf("unable to insert into %s", models.TableNames.BiosConfigSets), "duplicate key value violates unique constraint", "400"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			BiosConfigSetTemp := BiosConfigSetTest
			BiosConfigSetTemp.Name = tc.BiosConfigSetName
			BiosConfigSetTemp.ID = tc.BiosConfigSetID

			resp, err := s.Client.CreateServerBiosConfigSet(context.TODO(), BiosConfigSetTemp)
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				for _, msg := range tc.msgs {
					assert.Contains(t, err.Error(), msg)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				for _, msg := range tc.msgs {
					assert.Contains(t, resp.Message, msg)
				}

				if tc.BiosConfigSetID != "" {
					assert.Equal(t, tc.BiosConfigSetID, resp.Slug)
				} else {
					_, err = uuid.Parse(resp.Slug)
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestIntegrationServerBiosConfigSetGet(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, _ bool) error {
		s.Client.SetToken(authToken)

		parsedID, err := uuid.Parse(dbtools.FixtureBiosConfigSet.ID)
		require.NoError(t, err)

		resp, err := s.Client.GetServerBiosConfigSet(realClientTestCtx, parsedID)
		if err != nil {
			return err
		}

		require.NotNil(t, resp)

		assertEntireBiosConfigSetEqual(t,
			dbtools.FixtureBiosConfigSet,
			dbtools.FixtureBiosConfigComponents,
			dbtools.FixtureBiosConfigSettings,
			resp.Record.(*fleetdbapi.BiosConfigSet))

		return nil
	})

	var testCases = []struct {
		testName         string
		BiosConfigSetID  string
		expectedError    bool
		expectedResponse string
		msg              string
	}{
		{
			"config set router: config set get; success",
			dbtools.FixtureBiosConfigSet.ID,
			false,
			"200",
			"resource retrieved",
		},
		{
			"config set router: config set get; unknown config",
			uuid.NewString(),
			true,
			"404",
			"resource not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			id, err := uuid.Parse(tc.BiosConfigSetID)
			require.NoError(t, err)

			resp, err := s.Client.GetServerBiosConfigSet(context.TODO(), id)

			if tc.expectedError {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.msg)
				assert.Contains(t, err.Error(), tc.expectedResponse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Contains(t, resp.Message, tc.msg)

				// Fixtures are stored as models.BiosConfigSet, while the API returns fleetdbapi.BiosConfigSet, so we must manually compare the values
				assertEntireBiosConfigSetEqual(t,
					dbtools.FixtureBiosConfigSet,
					dbtools.FixtureBiosConfigComponents,
					dbtools.FixtureBiosConfigSettings,
					resp.Record.(*fleetdbapi.BiosConfigSet))
			}
		})
	}
}

func TestIntegrationServerBiosConfigSetDelete(t *testing.T) {
	s := serverTest(t)

	BiosConfigSetTemp := BiosConfigSetTest

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, expectedError bool) error {
		s.Client.SetToken(authToken)

		var parsedID uuid.UUID
		if expectedError {
			var err error
			parsedID, err = uuid.NewUUID()
			require.NoError(t, err)
		} else {
			BiosConfigSetTemp.Name = "Integration Test Config Set Delete"
			BiosConfigSetTemp.ID = ""
			resp, err := s.Client.CreateServerBiosConfigSet(realClientTestCtx, BiosConfigSetTemp)
			require.NoError(t, err)

			parsedID, err = uuid.Parse(resp.Slug)
			require.NoError(t, err)
		}

		_, err := s.Client.DeleteServerBiosConfigSet(realClientTestCtx, parsedID)
		if err != nil {
			return err
		}

		return nil
	})

	var testCases = []struct {
		testName         string
		BiosConfigSetID  string
		expectedError    bool
		expectedResponse string
		msg              string
	}{
		{
			"config set router: config set delete; success",
			"",
			false,
			"200",
			"resource deleted",
		},
		{
			"config set router: config set delete; unknown config",
			uuid.NewString(),
			true,
			"404",
			"resource not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			BiosConfigSetTemp.Name = tc.testName
			BiosConfigSetTemp.ID = ""
			resp, err := s.Client.CreateServerBiosConfigSet(context.TODO(), BiosConfigSetTemp)
			require.NoError(t, err)
			require.NotNil(t, resp)

			if tc.BiosConfigSetID == "" {
				tc.BiosConfigSetID = resp.Slug
			}

			parsedID, err := uuid.Parse(tc.BiosConfigSetID)
			require.NoError(t, err)
			require.NotNil(t, resp)

			resp, err = s.Client.DeleteServerBiosConfigSet(context.TODO(), parsedID)

			if tc.expectedError {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.msg)
				assert.Contains(t, err.Error(), tc.expectedResponse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Contains(t, resp.Message, tc.msg)
			}
		})
	}
}

func TestIntegrationServerBiosConfigSetUpdate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, expectedError bool) error {
		s.Client.SetToken(authToken)

		BiosConfigSetTemp := BiosConfigSetTest
		var parsedID uuid.UUID
		var err error

		if expectedError {
			parsedID, err = uuid.Parse(BiosConfigSetTest.ID)
			require.NoError(t, err)
		} else {
			BiosConfigSetTemp.ID = ""
			BiosConfigSetTemp.Name = "Integration Test Config Set Update"
			BiosConfigSetTemp.Version = "oldVersion"
			resp, err := s.Client.CreateServerBiosConfigSet(realClientTestCtx, BiosConfigSetTemp)
			require.NoError(t, err)
			require.NotNil(t, resp)

			BiosConfigSetTemp.ID = resp.Slug

			parsedID, err = uuid.Parse(resp.Slug)
			require.NoError(t, err)

			BiosConfigSetTemp.Version = "newVersion"
		}

		_, err = s.Client.UpdateServerBiosConfigSet(realClientTestCtx, parsedID, BiosConfigSetTemp)
		if err != nil {
			return err
		}

		if !expectedError {
			resp, err := s.Client.GetServerBiosConfigSet(realClientTestCtx, parsedID)
			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, BiosConfigSetTemp.Version, resp.Record.(*fleetdbapi.BiosConfigSet).Version)
		}

		return nil
	})

	var testCases = []struct {
		testName         string
		BiosConfigSetID  string
		expectedError    bool
		expectedResponse string
		msg              string
	}{
		{
			"config set router: config set update; success",
			BiosConfigSetTest.ID,
			false,
			"200",
			"resource updated",
		},
		{
			"config set router: config set update; invalid uuid",
			uuid.NewString(),
			true,
			"404",
			"resource not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			BiosConfigSetTemp := BiosConfigSetTest
			BiosConfigSetTemp.Version = "superNewVersion"

			parsedID, err := uuid.Parse(tc.BiosConfigSetID)
			require.NoError(t, err)

			if !tc.expectedError {
				BiosConfigSetTemp.ID = ""
				BiosConfigSetTemp.Name = "Integration Test Config Set Update 2"
				BiosConfigSetTemp.Version = "oldVersion"
				resp, err := s.Client.CreateServerBiosConfigSet(context.TODO(), BiosConfigSetTemp)
				require.NoError(t, err)
				require.NotNil(t, resp)

				BiosConfigSetTemp.ID = resp.Slug

				parsedID, err = uuid.Parse(resp.Slug)
				require.NoError(t, err)

				BiosConfigSetTemp.Version = "newVersion"
			}

			resp, err := s.Client.UpdateServerBiosConfigSet(context.TODO(), parsedID, BiosConfigSetTemp)

			if tc.expectedError {
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.msg)
				assert.Contains(t, err.Error(), tc.expectedResponse)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Contains(t, resp.Message, tc.msg)

				parsedID, err := uuid.Parse(resp.Slug)
				require.NoError(t, err)

				resp, err := s.Client.GetServerBiosConfigSet(context.TODO(), parsedID)
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, resp.Record.(*fleetdbapi.BiosConfigSet).Version, BiosConfigSetTemp.Version)
			}
		})
	}
}

func TestIntegrationServerBiosConfigSetList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, _ bool) error {
		s.Client.SetToken(authToken)

		testBiosConfigSetQueryParams := fleetdbapi.BiosConfigSetListParams{
			Params: []fleetdbapi.BiosConfigSetQueryParams{
				{
					Set: fleetdbapi.BiosConfigSetQuery{
						Name: "Fixture Test Config Set",
					},
					LogicalOperator:    fleetdbapi.OperatorLogicalOR,
					ComparitorOperator: fleetdbapi.OperatorComparitorEqual,
				},
			},
			Pagination: fleetdbapi.PaginationParams{
				Preload: true,
			},
		}

		_, err := s.Client.ListServerBiosConfigSet(realClientTestCtx, &testBiosConfigSetQueryParams)
		if err != nil {
			return err
		}

		return nil
	})

	// Setup for queries
	BiosConfigSetSetup := BiosConfigSetTest

	// Item 1
	BiosConfigSetSetup.ID = ""
	BiosConfigSetSetup.Name = "List Test 1"
	resp, err := s.Client.CreateServerBiosConfigSet(context.TODO(), BiosConfigSetSetup)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// Item 2
	BiosConfigSetSetup.ID = ""
	BiosConfigSetSetup.Name = "List Test 2"
	var tempComponentsArray []fleetdbapi.BiosConfigComponent
	for i := range BiosConfigSetSetup.Components { // Remove PCIE devices from components
		if BiosConfigSetSetup.Components[i].Model != "PCIE" {
			tempComponentsArray = append(tempComponentsArray, BiosConfigSetSetup.Components[i])
		}
	}
	backupComponents := BiosConfigSetSetup.Components
	BiosConfigSetSetup.Components = tempComponentsArray

	resp, err = s.Client.CreateServerBiosConfigSet(context.TODO(), BiosConfigSetSetup)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Item 3
	BiosConfigSetSetup.ID = ""
	BiosConfigSetSetup.Name = "List Test 3"
	newComponent := fleetdbapi.BiosConfigComponent{
		Name:   "TP-LinkNetwork Adapter",
		Vendor: "TP-Link",
		Model:  "PCIE",
		Settings: []fleetdbapi.BiosConfigSetting{
			{
				Key:   "PXEEnable",
				Value: "true",
				Raw:   []byte(`{}`),
			},
			{
				Key:   "SRIOVEnable",
				Value: "false",
			},
			{
				Key:   "position",
				Value: "2",
				Raw:   []byte(`{ "lanes": 8 }`),
			},
		},
	}
	BiosConfigSetSetup.Components = append(backupComponents, newComponent) // Add Additional PCIE device for param test 5

	for i := range BiosConfigSetSetup.Components { // Change SM Motherboard to Dell
		if BiosConfigSetSetup.Components[i].Name == "SM Motherboard" {
			BiosConfigSetSetup.Components[i].Name = "Dell Motherboard"
		}
	}

	resp, err = s.Client.CreateServerBiosConfigSet(context.TODO(), BiosConfigSetSetup)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Find none
	listTestParams0 := fleetdbapi.BiosConfigSetListParams{
		Params: []fleetdbapi.BiosConfigSetQueryParams{
			{
				Set: fleetdbapi.BiosConfigSetQuery{
					Name: "Not a good name",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorEqual,
			},
		},
		Pagination: fleetdbapi.PaginationParams{
			Preload: false,
		},
	}

	// Get all 3
	listTestParams1 := fleetdbapi.BiosConfigSetListParams{
		Params: []fleetdbapi.BiosConfigSetQueryParams{
			{
				Set: fleetdbapi.BiosConfigSetQuery{
					Name: "List Test",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
		},
		Pagination: fleetdbapi.PaginationParams{
			Preload: false,
		},
	}

	// Get all but "List Test 3"
	listTestParams2 := fleetdbapi.BiosConfigSetListParams{
		Params: []fleetdbapi.BiosConfigSetQueryParams{
			{
				Set: fleetdbapi.BiosConfigSetQuery{
					Name: "List Test",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
			{
				Set: fleetdbapi.BiosConfigSetQuery{
					Name: "List Test 3",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorNotEqual,
			},
		},
		Pagination: fleetdbapi.PaginationParams{
			Preload: false,
		},
	}

	// Get all based on components "%Motherboard"
	listTestParams3 := fleetdbapi.BiosConfigSetListParams{
		Params: []fleetdbapi.BiosConfigSetQueryParams{
			{
				Set: fleetdbapi.BiosConfigSetQuery{
					Name: "List Test",
					Components: []fleetdbapi.BiosConfigComponentQuery{
						{
							Name: "%Motherboard",
						},
					},
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
		},
		Pagination: fleetdbapi.PaginationParams{
			Preload: false,
		},
	}

	// Get all but the 3rd based on components "SM Motherboard"
	listTestParams4 := fleetdbapi.BiosConfigSetListParams{
		Params: []fleetdbapi.BiosConfigSetQueryParams{
			{
				Set: fleetdbapi.BiosConfigSetQuery{
					Name: "List Test",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
			{
				Set: fleetdbapi.BiosConfigSetQuery{
					Components: []fleetdbapi.BiosConfigComponentQuery{
						{
							Name: "SM Motherboard",
						},
					},
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorEqual,
			},
		},
		Pagination: fleetdbapi.PaginationParams{
			Preload: false,
		},
	}

	// Get all but second based on PCIE devices. This is to test "DISTINCT", since Item 3 has two PCIE components
	listTestParams5 := fleetdbapi.BiosConfigSetListParams{
		Params: []fleetdbapi.BiosConfigSetQueryParams{
			{
				Set: fleetdbapi.BiosConfigSetQuery{
					Name: "List Test",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
			{
				Set: fleetdbapi.BiosConfigSetQuery{
					Components: []fleetdbapi.BiosConfigComponentQuery{
						{
							Model: "PCIE",
						},
					},
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorEqual,
			},
		},
		Pagination: fleetdbapi.PaginationParams{
			Preload: true,
		},
	}

	var testCases = []struct {
		testName      string
		params        fleetdbapi.BiosConfigSetListParams
		expectedCount int
		expectedError bool
	}{
		{
			"config set router: config set list; find none",
			listTestParams0,
			0,
			false,
		},
		{
			"config set router: config set list; find all",
			listTestParams1,
			3,
			false,
		},
		{
			"config set router: config set list; NOT test",
			listTestParams2,
			2,
			false,
		},
		{
			"config set router: config set list; wildcard test",
			listTestParams3,
			3,
			false,
		},
		{
			"config set router: config set list; component query test",
			listTestParams4,
			2,
			false,
		},
		{
			"config set router: config set list; DISTINCT test",
			listTestParams5,
			2,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			resp, err := s.Client.ListServerBiosConfigSet(context.TODO(), &tc.params)

			if tc.expectedError {
				assert.Nil(t, resp)
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)

				assert.Equal(t, tc.expectedCount, len(*resp.Records.(*[]fleetdbapi.BiosConfigSet)))
			}
		})
	}
}

func assertEntireBiosConfigSetEqual(t *testing.T, expectedBiosConfigSet *models.BiosConfigSet, expectedComponents []*models.BiosConfigComponent, expectedSettings [][]*models.BiosConfigSetting, actual *fleetdbapi.BiosConfigSet) {
	assertBiosConfigSetEqual(t, expectedBiosConfigSet, actual)
	require.Equal(t, len(expectedComponents), len(actual.Components))
	require.Equal(t, len(expectedSettings), len(actual.Components))

	// Results can come back out of order, so we have to independently find components and settings
	for i, expectedComponent := range expectedComponents {
		foundComponent := false

		for _, actualComponent := range actual.Components {
			if actualComponent.ID == expectedComponent.ID {
				assertBiosConfigComponentEqual(t, expectedComponent, &actualComponent)

				for _, expectedSetting := range expectedSettings[i] {
					foundSetting := false

					for _, actualSetting := range actualComponent.Settings {
						if actualSetting.ID == expectedSetting.ID {
							assertBiosConfigSettingEqual(t, expectedSetting, &actualSetting)

							foundSetting = true

							break
						}
					}

					if !foundSetting {
						assert.Fail(t, fmt.Sprintf("expected to find component `%s` with setting `%s`", expectedComponent.Name, expectedSetting.SettingsKey))
					}
				}

				foundComponent = true

				break
			}
		}

		if !foundComponent {
			assert.Fail(t, fmt.Sprintf("expected to find component `%s`", expectedComponent.Name))
		}
	}
}

func assertBiosConfigSetEqual(t *testing.T, expected *models.BiosConfigSet, actual *fleetdbapi.BiosConfigSet) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Version, actual.Version)

	assert.WithinDuration(t, expected.CreatedAt.Time, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt.Time, actual.UpdatedAt, time.Second)
}

func assertBiosConfigComponentEqual(t *testing.T, expected *models.BiosConfigComponent, actual *fleetdbapi.BiosConfigComponent) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Vendor, actual.Vendor)
	assert.Equal(t, expected.Model, actual.Model)

	assert.WithinDuration(t, expected.CreatedAt.Time, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt.Time, actual.UpdatedAt, time.Second)
}

func assertBiosConfigSettingEqual(t *testing.T, expected *models.BiosConfigSetting, actual *fleetdbapi.BiosConfigSetting) {
	assert.Equal(t, expected.SettingsKey, actual.Key)
	assert.Equal(t, expected.SettingsValue, actual.Value)

	assert.WithinDuration(t, expected.CreatedAt.Time, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt.Time, actual.UpdatedAt, time.Second)

	if expected.Raw.IsZero() {
		assert.Nil(t, actual.Raw)
	} else {
		assert.JSONEq(t, string(expected.Raw.JSON), string(actual.Raw))
	}
}
