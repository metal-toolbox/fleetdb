package fleetdbapi_test

import (
	"bytes"
	"context"
	"encoding/json"

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

func TestIntegrationServerConfigSetCreate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, _ bool) error {
		s.Client.SetToken(authToken)

		configSetTemp := configSetTest
		configSetTemp.ID = ""
		resp, err := s.Client.CreateServerConfigSet(realClientTestCtx, configSetTest)
		if err != nil {
			return err
		}

		_, err = uuid.Parse(resp.Slug)
		assert.NoError(t, err)

		configSetTest.ID = resp.Slug

		return nil
	})

	var testCases = []struct {
		testName      string
		configSetName string
		configSetID   string
		expectedError bool
		msgs          []string
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
			[]string{"invalid payload: ConfigSetCreate{}", "Field validation for 'Name' failed on the 'required' tag", "400"},
		},
		{
			"config set router: config set create; duplicate config",
			"Integration Test Config Set Fail 2",
			dbtools.FixtureConfigSet.ID,
			true,
			[]string{fmt.Sprintf("unable to insert into %s", models.TableNames.ConfigSets), "duplicate key value violates unique constraint", "400"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			configSetTemp := configSetTest
			configSetTemp.Name = tc.configSetName
			configSetTemp.ID = tc.configSetID

			resp, err := s.Client.CreateServerConfigSet(context.TODO(), configSetTemp)
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

				if tc.configSetID != "" {
					assert.Equal(t, tc.configSetID, resp.Slug)
				} else {
					_, err = uuid.Parse(resp.Slug)
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestIntegrationServerConfigSetGet(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, _ bool) error {
		s.Client.SetToken(authToken)

		parsedID, err := uuid.Parse(dbtools.FixtureConfigSet.ID)
		require.NoError(t, err)

		resp, err := s.Client.GetServerConfigSet(realClientTestCtx, parsedID)
		if err != nil {
			return err
		}

		assertEntireConfigSetEqual(t,
			dbtools.FixtureConfigSet,
			dbtools.FixtureConfigComponents,
			dbtools.FixtureConfigComponentSettings,
			resp.Record.(*fleetdbapi.ConfigSet))

		return nil
	})

	var testCases = []struct {
		testName         string
		configSetID      string
		expectedError    bool
		expectedResponse string
		msg              string
	}{
		{
			"config set router: config set get; success",
			dbtools.FixtureConfigSet.ID,
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
			id, err := uuid.Parse(tc.configSetID)
			assert.NoError(t, err)

			resp, err := s.Client.GetServerConfigSet(context.TODO(), id)

			if tc.expectedError {
				assert.Nil(t, resp)
				assert.Nil(t, resp)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.msg)
				assert.Contains(t, err.Error(), tc.expectedResponse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Contains(t, resp.Message, tc.msg)

				var set *fleetdbapi.ConfigSet = resp.Record.(*fleetdbapi.ConfigSet)

				// Fixtures are stored as models.ConfigSet, while the API returns fleetdbapi.ConfigSet, so we must manually compare the values
				assertEntireConfigSetEqual(t,
					dbtools.FixtureConfigSet,
					dbtools.FixtureConfigComponents,
					dbtools.FixtureConfigComponentSettings,
					set)
			}
		})
	}
}

func TestIntegrationServerConfigSetDelete(t *testing.T) {
	s := serverTest(t)

	configSetTemp := configSetTest

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, expectedError bool) error {
		s.Client.SetToken(authToken)

		var parsedID uuid.UUID
		if expectedError {
			var err error
			parsedID, err = uuid.NewUUID()
			require.NoError(t, err)
		} else {
			configSetTemp.Name = "Integration Test Config Set Delete"
			configSetTemp.ID = ""
			resp, err := s.Client.CreateServerConfigSet(realClientTestCtx, configSetTemp)
			require.NoError(t, err)

			parsedID, err = uuid.Parse(resp.Slug)
			require.NoError(t, err)
		}

		_, err := s.Client.DeleteServerConfigSet(realClientTestCtx, parsedID)
		if err != nil {
			return err
		}

		return nil
	})

	var testCases = []struct {
		testName         string
		configSetID      string
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
			configSetTemp.Name = tc.testName
			configSetTemp.ID = ""
			resp, err := s.Client.CreateServerConfigSet(context.TODO(), configSetTemp)
			require.NoError(t, err)
			require.NotNil(t, resp)

			if tc.configSetID == "" {
				tc.configSetID = resp.Slug
			}

			parsedID, err := uuid.Parse(tc.configSetID)
			require.NoError(t, err)
			require.NotNil(t, resp)

			resp, err = s.Client.DeleteServerConfigSet(context.TODO(), parsedID)

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

func TestIntegrationServerConfigSetUpdate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, expectedError bool) error {
		s.Client.SetToken(authToken)

		configSetTemp := configSetTest
		var parsedID uuid.UUID
		var err error

		if expectedError {
			parsedID, err = uuid.Parse(configSetTest.ID)
			require.NoError(t, err)
		} else {
			configSetTemp.ID = ""
			configSetTemp.Name = "Integration Test Config Set Update"
			configSetTemp.Version = "oldVersion"
			resp, err := s.Client.CreateServerConfigSet(realClientTestCtx, configSetTemp)
			require.NoError(t, err)
			require.NotNil(t, resp)

			configSetTemp.ID = resp.Slug

			parsedID, err = uuid.Parse(resp.Slug)
			require.NoError(t, err)

			configSetTemp.Version = "newVersion"
		}

		_, err = s.Client.UpdateServerConfigSet(realClientTestCtx, parsedID, configSetTemp)
		if err != nil {
			return err
		}

		if !expectedError {
			resp, err := s.Client.GetServerConfigSet(realClientTestCtx, parsedID)
			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, configSetTemp.Version, resp.Record.(*fleetdbapi.ConfigSet).Version)
		}

		return nil
	})

	var testCases = []struct {
		testName         string
		configSetID      string
		expectedError    bool
		expectedResponse string
		msg              string
	}{
		{
			"config set router: config set update; success",
			configSetTest.ID,
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
			configSetTemp := configSetTest
			configSetTemp.Version = "superNewVersion"

			parsedID, err := uuid.Parse(tc.configSetID)
			require.NoError(t, err)

			if !tc.expectedError {
				configSetTemp.ID = ""
				configSetTemp.Name = "Integration Test Config Set Update 2"
				configSetTemp.Version = "oldVersion"
				resp, err := s.Client.CreateServerConfigSet(context.TODO(), configSetTemp)
				require.NoError(t, err)
				require.NotNil(t, resp)

				configSetTemp.ID = resp.Slug

				parsedID, err = uuid.Parse(resp.Slug)
				require.NoError(t, err)

				configSetTemp.Version = "newVersion"
			}

			resp, err := s.Client.UpdateServerConfigSet(context.TODO(), parsedID, configSetTemp)

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

				resp, err := s.Client.GetServerConfigSet(context.TODO(), parsedID)
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, resp.Record.(*fleetdbapi.ConfigSet).Version, configSetTemp.Version)
			}
		})
	}
}

func TestIntegrationServerConfigSetList(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(realClientTestCtx context.Context, authToken string, _ int, _ bool) error {
		s.Client.SetToken(authToken)

		testConfigSetQueryParams := fleetdbapi.ConfigSetListParams{
			Params: []fleetdbapi.ConfigSetQueryParams{
				{
					Set: fleetdbapi.ConfigSetQuery{
						Name: "Fixture Test Config Set",
					},
					LogicalOperator:    fleetdbapi.OperatorLogicalOR,
					ComparitorOperator: fleetdbapi.OperatorComparitorEqual,
				},
			},
			Pagination: fleetdbapi.PaginationParams{},
		}

		_, err := s.Client.ListServerConfigSet(realClientTestCtx, &testConfigSetQueryParams)
		if err != nil {
			return err
		}

		return nil
	})

	// Setup for queries
	configSetSetup := configSetTest
	configSetSetup.ID = ""
	configSetSetup.Name = "List Test 1"
	resp, err := s.Client.CreateServerConfigSet(context.TODO(), configSetSetup)
	require.NoError(t, err)
	require.NotNil(t, resp)

	configSetSetup.ID = ""
	configSetSetup.Name = "List Test 2"
	resp, err = s.Client.CreateServerConfigSet(context.TODO(), configSetSetup)
	require.NoError(t, err)
	require.NotNil(t, resp)

	configSetSetup.ID = ""
	configSetSetup.Name = "List Test 3"
	configSetSetup.Components[0].Name = "Dell Motherboard"
	resp, err = s.Client.CreateServerConfigSet(context.TODO(), configSetSetup)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Find none
	listTestParams0 := fleetdbapi.ConfigSetListParams{
		Params: []fleetdbapi.ConfigSetQueryParams{
			{
				Set: fleetdbapi.ConfigSetQuery{
					Name: "Not a good name",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorEqual,
			},
		},
		Pagination: fleetdbapi.PaginationParams{},
	}

	// Get all 3
	listTestParams1 := fleetdbapi.ConfigSetListParams{
		Params: []fleetdbapi.ConfigSetQueryParams{
			{
				Set: fleetdbapi.ConfigSetQuery{
					Name: "List Test",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
		},
		Pagination: fleetdbapi.PaginationParams{},
	}

	// Get all but "List Test 3"
	listTestParams2 := fleetdbapi.ConfigSetListParams{
		Params: []fleetdbapi.ConfigSetQueryParams{
			{
				Set: fleetdbapi.ConfigSetQuery{
					Name: "List Test",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
			{
				Set: fleetdbapi.ConfigSetQuery{
					Name: "List Test 3",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorNotEqual,
			},
		},
		Pagination: fleetdbapi.PaginationParams{},
	}

	// Get all based on components "%Motherboard"
	listTestParams3 := fleetdbapi.ConfigSetListParams{
		Params: []fleetdbapi.ConfigSetQueryParams{
			{
				Set: fleetdbapi.ConfigSetQuery{
					Name: "List Test",
					Components: []fleetdbapi.ConfigComponentQuery{
						{
							Name: "%Motherboard",
						},
					},
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
		},
		Pagination: fleetdbapi.PaginationParams{},
	}
	// Get all but the 3rd based on components "SM Motherboard"
	listTestParams4 := fleetdbapi.ConfigSetListParams{
		Params: []fleetdbapi.ConfigSetQueryParams{
			{
				Set: fleetdbapi.ConfigSetQuery{
					Name: "List Test",
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorLike,
			},
			{
				Set: fleetdbapi.ConfigSetQuery{
					Components: []fleetdbapi.ConfigComponentQuery{
						{
							Name: "SM Motherboard",
						},
					},
				},
				LogicalOperator:    fleetdbapi.OperatorLogicalAND,
				ComparitorOperator: fleetdbapi.OperatorComparitorEqual,
			},
		},
		Pagination: fleetdbapi.PaginationParams{},
	}

	var testCases = []struct {
		testName      string
		params        fleetdbapi.ConfigSetListParams
		expectedCount int
		expectedError bool
	}{
		{
			"config set router: config set list; success 0",
			listTestParams0,
			0,
			false,
		},
		{
			"config set router: config set list; success 1",
			listTestParams1,
			3,
			false,
		},
		{
			"config set router: config set list; success 2",
			listTestParams2,
			2,
			false,
		},
		{
			"config set router: config set list; success 3",
			listTestParams3,
			3,
			false,
		},
		{
			"config set router: config set list; success 4",
			listTestParams4,
			2,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			resp, err := s.Client.ListServerConfigSet(context.TODO(), &tc.params)

			if tc.expectedError {
				assert.Nil(t, resp)
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)

				assert.Equal(t, tc.expectedCount, len(*resp.Records.(*[]fleetdbapi.ConfigSet)))
			}
		})
	}
}

func assertEntireConfigSetEqual(t *testing.T, expectedConfigSet *models.ConfigSet, expectedComponents []*models.ConfigComponent, expectedSettings [][]*models.ConfigComponentSetting, actual *fleetdbapi.ConfigSet) {
	assertConfigSetEqual(t, expectedConfigSet, actual)
	require.Equal(t, len(expectedComponents), len(actual.Components))
	require.Equal(t, len(expectedSettings), len(actual.Components))

	// Results can come back out of order, so we have to independently find components and settings
	for i, expectedComponent := range expectedComponents {
		foundComponent := false

		for _, actualComponent := range actual.Components {
			if actualComponent.ID == expectedComponent.ID {
				assertConfigComponentEqual(t, expectedComponent, &actualComponent)

				for _, expectedSetting := range expectedSettings[i] {
					foundSetting := false

					for _, actualSetting := range actualComponent.Settings {
						if actualSetting.ID == expectedSetting.ID {
							assertConfigComponentSettingEqual(t, expectedSetting, &actualSetting)

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

func assertConfigSetEqual(t *testing.T, expected *models.ConfigSet, actual *fleetdbapi.ConfigSet) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Version.String, actual.Version)

	assert.WithinDuration(t, expected.CreatedAt.Time, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt.Time, actual.UpdatedAt, time.Second)

}

func assertConfigComponentEqual(t *testing.T, expected *models.ConfigComponent, actual *fleetdbapi.ConfigComponent) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Vendor.String, actual.Vendor)
	assert.Equal(t, expected.Model.String, actual.Model)
	assert.Equal(t, expected.Serial.String, actual.Serial)

	assert.WithinDuration(t, expected.CreatedAt.Time, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt.Time, actual.UpdatedAt, time.Second)
}

func assertConfigComponentSettingEqual(t *testing.T, expected *models.ConfigComponentSetting, actual *fleetdbapi.ConfigComponentSetting) {
	assert.Equal(t, expected.SettingsKey, actual.Key)
	assert.Equal(t, expected.SettingsValue, actual.Value)

	assert.WithinDuration(t, expected.CreatedAt.Time, actual.CreatedAt, time.Second)
	assert.WithinDuration(t, expected.UpdatedAt.Time, actual.UpdatedAt, time.Second)

	if expected.Custom.IsZero() {
		assert.Nil(t, actual.Custom)
	} else { // JSON []byte can have excess namespace, so lets remove that when comparing
		var expectedBuffer bytes.Buffer
		err := json.Compact(&expectedBuffer, expected.Custom.JSON)
		require.NoError(t, err)

		var actualBuffer bytes.Buffer
		err = json.Compact(&actualBuffer, actual.Custom)
		require.NoError(t, err)

		assert.Equal(t, expectedBuffer.Bytes(), actualBuffer.Bytes())
	}
}
