package fleetdbapi_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/models"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

// zero values that change for each test run to enable object comparison
func zeroUUIDValues(sc *fleetdbapi.ServerComponent) {
	sc.ServerUUID = uuid.UUID{}
	sc.UUID = uuid.UUID{}
	sc.ComponentTypeID = ""
}

func zeroTimeValues(sc *fleetdbapi.ServerComponent) {
	sc.CreatedAt = time.Time{}
	sc.UpdatedAt = time.Time{}

	for i := 0; i < len(sc.VersionedAttributes); i++ {
		sc.VersionedAttributes[i].CreatedAt = time.Time{}
		sc.VersionedAttributes[i].LastReportedAt = time.Time{}
	}
}

func componentByNameVendorModelSerial(expected fleetdbapi.ServerComponent, actualSlice fleetdbapi.ServerComponentSlice) *fleetdbapi.ServerComponent {
	for _, c := range actualSlice {
		if c.Name == expected.Name && c.Vendor == expected.Vendor && c.Model == expected.Model && c.Serial == expected.Serial {
			return &c
		}
	}

	return nil
}

func TestIntegrationServerListComponents(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, _ int, expectError bool) error {
		s.Client.SetToken(authToken)

		components, _, err := s.Client.ListComponents(ctx, nil)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, components, 7)
		}

		return err
	})

	testCases := []struct {
		testName string
		params   *fleetdbapi.ServerComponentListParams
		expected fleetdbapi.ServerComponentSlice
		errorMsg string
	}{
		// TODO(joel): tests for unhappy paths
		{
			"by model",
			&fleetdbapi.ServerComponentListParams{Model: "Belly"},
			fleetdbapi.ServerComponentSlice{
				{
					Model:             "Belly",
					Serial:            "Up",
					ComponentTypeName: "Fins",
					ComponentTypeSlug: "fins",
				},
			},
			"",
		},
		{
			"by model, versioned attributes",
			&fleetdbapi.ServerComponentListParams{
				Model: "Normal Fin",
				VersionedAttributeListParams: []fleetdbapi.AttributeListParams{
					{
						Namespace: "hollow.versioned",
						Keys:      []string{"something"},
						Operator:  "eq",
						Value:     "cool",
					},
				},
			},
			fleetdbapi.ServerComponentSlice{
				{
					Model:             "Normal Fin",
					Serial:            "Left",
					Name:              "Normal Fin",
					ComponentTypeName: "Fins",
					ComponentTypeSlug: "fins",
					VersionedAttributes: []fleetdbapi.VersionedAttributes{
						{
							Namespace: "hollow.versioned",
							Data:      json.RawMessage(`{"something":"cool"}`),
						},
					},
				},
			},
			"",
		},
		{
			"pagination Limit",
			&fleetdbapi.ServerComponentListParams{
				Pagination: &fleetdbapi.PaginationParams{
					Limit: 1,
				},
				Model: "Belly",
			},
			fleetdbapi.ServerComponentSlice{
				{
					Model:             "Belly",
					Serial:            "Up",
					ComponentTypeName: "Fins",
					ComponentTypeSlug: "fins",
				},
			},
			"",
		},
		{
			"pagination Limit, Offset",
			&fleetdbapi.ServerComponentListParams{
				Pagination: &fleetdbapi.PaginationParams{
					Limit:   1,
					Page:    2,
					OrderBy: models.ServerComponentTableColumns.CreatedAt + " DESC",
				},
			},
			fleetdbapi.ServerComponentSlice{
				{
					Name:              "Normal Fin",
					Serial:            "Right",
					ComponentTypeName: "Fins",
					ComponentTypeSlug: "fins",
				},
			},
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			got, res, err := s.Client.ListComponents(context.TODO(), tc.params)
			if tc.errorMsg != "" {
				assert.Contains(t, err.Error(), tc.errorMsg)
				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, res)

			// zero timestamp, uuid values for object comparison
			for i := 0; i < len(got); i++ {
				zeroUUIDValues(&got[i])
				zeroTimeValues(&got[i])
			}

			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestIntegrationServerGetComponents(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, _ int, expectError bool) error {
		s.Client.SetToken(authToken)

		components, _, err := s.Client.GetComponents(ctx, uuid.MustParse(dbtools.FixtureNemo.ID), nil)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, components, 2)
		}

		return err
	})

	// init fixture data

	// 1. get list of servers
	listParams := &fleetdbapi.ServerListParams{
		PaginationParams: &fleetdbapi.PaginationParams{
			Preload: true,
			OrderBy: models.ServerColumns.CreatedAt + " DESC",
		},
	}
	servers, _, err := s.Client.List(context.Background(), listParams)
	if err != nil {
		t.Fatal(err)
	}

	// expect atleast 1 server for test to proceed
	assert.GreaterOrEqual(t, len(servers), 1)

	// 2. get component type slice
	componentTypeSlice, _, err := s.Client.ListServerComponentTypes(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// expect atleast 1 component type to proceed
	assert.GreaterOrEqual(t, len(componentTypeSlice), 1)

	// fixture to create a server components
	csFixtureCreate := fleetdbapi.ServerComponentSlice{
		{
			ServerUUID:        servers[1].UUID,
			Name:              "My Lucky Fin",
			Vendor:            "barracuda",
			Model:             "a lucky fin",
			Serial:            "right",
			ComponentTypeID:   componentTypeSlice.ByName("Fins").ID,
			ComponentTypeName: componentTypeSlice.ByName("Fins").Name,
			ComponentTypeSlug: componentTypeSlice.ByName("Fins").Slug,
			VersionedAttributes: []fleetdbapi.VersionedAttributes{
				{
					Namespace: dbtools.FixtureNamespaceVersioned,
					Data:      json.RawMessage(`{"version":"1.0"}`),
				},
				{
					Namespace: dbtools.FixtureNamespaceVersioned,
					Data:      json.RawMessage(`{"version":"2.0"}`),
				},
			},
		},
	}

	// create server component
	_, err = s.Client.CreateComponents(context.TODO(), servers[1].UUID, csFixtureCreate)
	if err != nil {
		t.Fatal(err)
	}

	var testCases = []struct {
		testName        string
		srvUUID         uuid.UUID
		expectedCount   int
		expectedInSlice fleetdbapi.ServerComponent
		errorMsg        string
	}{
		{
			"returns not found on missing server uuid",
			uuid.New(),
			0,
			fleetdbapi.ServerComponent{},
			"response code: 404",
		},
		{
			"component Versioned Attributes is returned as expected",
			servers[1].UUID,
			3,
			fleetdbapi.ServerComponent{
				ServerUUID:        servers[1].UUID,
				Name:              "My Lucky Fin",
				Vendor:            "barracuda",
				Model:             "a lucky fin",
				Serial:            "right",
				ComponentTypeID:   componentTypeSlice.ByName("Fins").ID,
				ComponentTypeName: componentTypeSlice.ByName("Fins").Name,
				ComponentTypeSlug: componentTypeSlice.ByName("Fins").Slug,
				VersionedAttributes: []fleetdbapi.VersionedAttributes{
					{
						Namespace: dbtools.FixtureNamespaceVersioned,
						Data:      json.RawMessage(`{"version":"2.0"}`),
					},
				},
			},
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			got, _, err := s.Client.GetComponents(context.TODO(), tt.srvUUID, nil)
			if tt.errorMsg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)

				return
			}

			assert.Nil(t, err)

			assert.Equal(t, tt.expectedCount, len(got))
			gotc := componentByNameVendorModelSerial(
				tt.expectedInSlice,
				got,
			)

			if gotc != nil {
				// zero variable values before comparison
				gotc.UUID = uuid.Nil
				zeroTimeValues(gotc)

				assert.Equal(t, tt.expectedInSlice, *gotc)
			} else {
				t.Fatal("expected component, got nil")
			}
		})
	}
}

func TestIntegrationServerCreateComponents(t *testing.T) {
	s := serverTest(t)

	// fixture objects
	var servers []fleetdbapi.Server

	var componentTypeSlice fleetdbapi.ServerComponentTypeSlice

	// run default client tests
	realClientTests(t, func(ctx context.Context, authToken string, _ int, expectError bool) error {
		s.Client.SetToken(authToken)

		var sc fleetdbapi.ServerComponentSlice

		if !expectError {
			var err error
			// 2. retrieve list of servers, expect the test db is populated with one or more test servers
			servers, _, err = s.Client.List(context.Background(), nil)
			if err != nil {
				t.Fatal(err)
			}

			componentTypeSlice, _, err = s.Client.ListServerComponentTypes(context.Background(), nil)
			if err != nil {
				t.Fatal(err)
			}

			sc = fleetdbapi.ServerComponentSlice{
				{
					ServerUUID:        servers[0].UUID,
					ComponentTypeID:   componentTypeSlice[0].ID,
					ComponentTypeName: componentTypeSlice[0].Name,
					ComponentTypeSlug: componentTypeSlice[0].Slug,
					Name:              "Fin A",
					Model:             "Normal Fin",
					Serial:            "Left Upper",
				},
			}
		}

		_, err := s.Client.CreateComponents(ctx, uuid.MustParse(dbtools.FixtureNemo.ID), sc)
		if !expectError {
			require.NoError(t, err)
		}

		return err
	})

	// make sure all fixtures are available
	assert.GreaterOrEqual(t, len(servers), 1)
	assert.GreaterOrEqual(t, len(componentTypeSlice), 1)

	var testCases = []struct {
		testName    string
		serverUUID  uuid.UUID
		components  fleetdbapi.ServerComponentSlice
		responseMsg string
		errorMsg    string
	}{
		{
			"unknown server query returns 404",
			uuid.New(),
			nil,
			"",
			"hollow client received a server error - response code: 404, message: resource not found",
		},
		{
			"create component and list by Name works",
			servers[0].UUID,
			fleetdbapi.ServerComponentSlice{
				{
					ServerUUID:        servers[0].UUID,
					ComponentTypeID:   componentTypeSlice[0].ID,
					ComponentTypeName: componentTypeSlice[0].Name,
					ComponentTypeSlug: componentTypeSlice[0].Slug,
					Name:              "Fin B",
					Model:             "Normal Fin",
					Serial:            "Left Lower",
				},
			},
			"resource created",
			"",
		},
		{
			"create component which violates unique constraint on ServerID, ComponentTypeID, Serial should return error",
			servers[0].UUID,
			fleetdbapi.ServerComponentSlice{
				{
					ServerUUID:        servers[0].UUID,
					ComponentTypeID:   componentTypeSlice[0].ID,
					ComponentTypeName: componentTypeSlice[0].Name,
					ComponentTypeSlug: componentTypeSlice[0].Slug,
					Name:              "Fin B",
					Model:             "Normal Fin",
					Serial:            "Left Lower",
				},
				{
					ServerUUID:        servers[0].UUID,
					ComponentTypeID:   componentTypeSlice[0].ID,
					ComponentTypeName: componentTypeSlice[0].Name,
					ComponentTypeSlug: componentTypeSlice[0].Slug,
					Name:              "Fin B",
					Model:             "Normal Fin",
					Serial:            "Left Lower",
				},
			},
			"",
			"duplicate key value violates unique constraint",
		},
		{
			"create component with unknown server UUID returns error",
			uuid.New(),
			fleetdbapi.ServerComponentSlice{
				{
					ServerUUID:        uuid.New(),
					ComponentTypeID:   componentTypeSlice[0].ID,
					ComponentTypeName: componentTypeSlice[0].Name,
					ComponentTypeSlug: componentTypeSlice[0].Slug,
					Name:              "Fin B",
					Model:             "Normal Fin",
					Serial:            "Left Lower 2",
				},
			},
			"",
			"resource not found",
		},
		{
			"create component validates field constraints",
			servers[0].UUID,
			fleetdbapi.ServerComponentSlice{
				{
					ServerUUID:      servers[0].UUID,
					ComponentTypeID: "lala",
					Name:            "Fin B",
					Model:           "Normal Fin",
					Serial:          "Left Lower 2",
				},
			},
			"",
			"error in server component payload",
		},
		{
			"create component with empty slice returns error",
			servers[0].UUID,
			fleetdbapi.ServerComponentSlice{},
			"",
			"error in server component payload",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			// create
			res, err := s.Client.CreateComponents(context.TODO(), tt.serverUUID, tt.components)
			if tt.errorMsg != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)

				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, res)
			assert.Contains(t, res.Message, tt.responseMsg)

			params := &fleetdbapi.ServerComponentListParams{Name: tt.components[0].Name, Pagination: &fleetdbapi.PaginationParams{Preload: true}}

			got, _, err := s.Client.ListComponents(context.TODO(), params)
			if err != nil {
				t.Error(err)
			}

			// zero timestamp values for object comparison
			for i := 0; i < len(got); i++ {
				zeroTimeValues(&got[i])

				got[i].UUID = uuid.Nil
			}

			assert.Equal(t, tt.components, got)
		})
	}
}

func TestIntegrationServerUpdateComponents(t *testing.T) {
	s := serverTest(t)
	// fixture objects
	var servers []fleetdbapi.Server

	// run default client tests
	realClientTests(t, func(ctx context.Context, authToken string, _ int, expectError bool) error {
		s.Client.SetToken(authToken)

		var sc fleetdbapi.ServerComponentSlice

		if !expectError {
			var err error

			// In order to get the components, we need to preload them
			params := &fleetdbapi.ServerListParams{
				PaginationParams: &fleetdbapi.PaginationParams{
					Preload: true,
				},
			}

			// 2. retrieve list of servers, expect the test db is populated with one or more test servers
			servers, _, err = s.Client.List(context.Background(), params)
			if err != nil {
				t.Fatal(err)
			}

			var nemo *fleetdbapi.Server
			for _, s := range servers {
				if s.Name == dbtools.FixtureNemo.Name.String {
					nemo = &s
				}
			}
			require.NotNil(t, nemo, "couldn't find nemo")

			// update serial attribute for update to work
			require.NotZero(t, len(nemo.Components))
			sc = fleetdbapi.ServerComponentSlice{nemo.Components[0]}
			sc[0].Serial = "lala"
		}

		_, err := s.Client.UpdateComponents(ctx, uuid.MustParse(dbtools.FixtureNemo.ID), sc)
		if !expectError {
			require.NoError(t, err)
		}

		return err
	})

	// fixtures given to test cases below
	var componentFixture []fleetdbapi.ServerComponent

	var serverFixture fleetdbapi.Server

	// The component fixture targeted in test cases below
	fixtureComponentName := "My Lucky Fin"
	fixtureComponentVendor := "Barracuda"
	fixtureComponentSerial := "Right"

	// identify component and server fixture for test
	for _, server := range servers {
		for _, c := range server.Components {
			if c.Name == fixtureComponentName && c.Vendor == fixtureComponentVendor && c.Serial == fixtureComponentSerial {
				componentFixture = append(componentFixture, c)
				serverFixture = server
			}
		}
	}

	// helper method to return fixture copy
	componentFixtureCopy := func() []fleetdbapi.ServerComponent {
		var c []fleetdbapi.ServerComponent

		c = append(c, componentFixture...)

		return c
	}

	// expect test fixture to be present
	assert.NotEmpty(t, componentFixture[0].UUID)
	assert.NotEmpty(t, serverFixture.UUID)

	// change are changes to be applied to the fixture object included in each test case
	type change struct {
		versionedAttributes json.RawMessage
		attributes          json.RawMessage
		// unsetFlags purges component attributes on the first component in the components slice
		// included in the test case.
		// bool value index:
		// 0 - unsetUUID
		// 1 - unset componentSerial
		// 2 - unset componentServer ID
		// 3 - unset component type UUID
		unsetFlags []bool
	}

	var testCases = []struct {
		testName    string
		serverUUID  uuid.UUID
		components  fleetdbapi.ServerComponentSlice
		change      change
		responseMsg string
		errorMsg    string
	}{
		{
			"component update for unknown server return error",
			uuid.New(),
			nil,
			change{},
			"",
			"resource not found",
		},
		{
			"component update with empty component slice returns error",
			serverFixture.UUID,
			nil,
			change{},
			"",
			"ServerComponentSlice is empty",
		},
		{
			"component update validation for non-nil UUID returns error",
			serverFixture.UUID,
			componentFixtureCopy(),
			// unset component uuid
			change{unsetFlags: []bool{true, false, false, false}},
			"",
			"component update requires a non-nil UUID",
		},
		{
			"component update validation for empty serial returns error",
			serverFixture.UUID,
			componentFixtureCopy(),
			// unset component serial
			change{unsetFlags: []bool{false, true, false, false}},
			"",
			"Field validation for 'Serial' failed on the 'required' tag",
		},
		{
			"component update validation for empty server ID returns error",
			serverFixture.UUID,
			componentFixtureCopy(),
			// unset component server ID
			change{unsetFlags: []bool{false, false, true, false}},
			"",
			"Field validation for 'ServerUUID' failed on the 'required' tag",
		},
		{
			"component update validation for empty component type ID returns error",
			serverFixture.UUID,
			componentFixtureCopy(),
			// unset component type UUID
			change{unsetFlags: []bool{false, false, false, true}},
			"",
			"Field validation for 'ComponentTypeID' failed on the 'required' tag",
		},
		{
			"component update on versioned attributes",
			serverFixture.UUID,
			componentFixtureCopy(),
			change{versionedAttributes: []byte(`{"version":"2.12345"}`)},
			"",
			"",
		},
		{
			"component update on attributes - attribute gets updated",
			serverFixture.UUID,
			componentFixtureCopy(),
			change{attributes: []byte(`{"twitches":"false"}`)},
			"",
			"",
		},
		{
			"component update on attributes - attribute gets added",
			serverFixture.UUID,
			componentFixtureCopy(),
			change{attributes: []byte(`{"wavy":"true"}`)},
			"",
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			// check flags to unset components attributes
			if len(tt.change.unsetFlags) == 4 {
				if tt.change.unsetFlags[0] {
					tt.components[0].UUID = uuid.Nil
				}

				if tt.change.unsetFlags[1] {
					tt.components[0].Serial = ""
				}

				if tt.change.unsetFlags[2] {
					tt.components[0].ServerUUID = uuid.Nil
				}

				if tt.change.unsetFlags[3] {
					tt.components[0].ComponentTypeID = ""
				}
			}

			var listParams *fleetdbapi.ServerComponentListParams

			// test case updates versioned attributes
			if len(tt.change.versionedAttributes) > 0 {
				tt.components[0].VersionedAttributes = []fleetdbapi.VersionedAttributes{
					{
						Namespace: "hollow.metadata",
						Data:      tt.change.versionedAttributes,
					},
				}

				model := "testUpdatedVersionedAttributes" + time.Now().String()
				tt.components[0].Model = model

				listParams = &fleetdbapi.ServerComponentListParams{
					Name:       fixtureComponentName,
					Serial:     fixtureComponentSerial,
					Vendor:     fixtureComponentVendor,
					Model:      model,
					Pagination: &fleetdbapi.PaginationParams{Preload: true},
				}
			}

			// test case updates attributes
			if len(tt.change.attributes) > 0 {
				tt.components[0].Attributes = []fleetdbapi.Attributes{
					{
						Namespace: dbtools.FixtureNamespaceOtherdata,
						Data:      tt.change.attributes,
					},
				}

				model := "testUpdatedAttributes" + time.Now().String()
				tt.components[0].Model = model

				listParams = &fleetdbapi.ServerComponentListParams{
					Name:       fixtureComponentName,
					Serial:     fixtureComponentSerial,
					Vendor:     fixtureComponentVendor,
					Model:      model,
					Pagination: &fleetdbapi.PaginationParams{Preload: true},
				}
			}

			//	update component
			_, err := s.Client.UpdateComponents(context.TODO(), tt.serverUUID, tt.components)
			// assert any expected errors
			if tt.errorMsg != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)

				return
			}

			// asert no errors
			assert.Nil(t, err)

			// list component updated
			got, _, err := s.Client.ListComponents(context.TODO(), listParams)
			if err != nil {
				t.Error(err)
			}

			// assert versioned attributes change
			if len(tt.change.versionedAttributes) > 0 {
				assert.Len(t, got, 1)
				assert.Equal(t, tt.change.versionedAttributes, got[0].VersionedAttributes[0].Data)
			}

			// assert attributes change
			if len(tt.change.attributes) > 0 {
				assert.Len(t, got, 1)
				assert.Equal(t, tt.change.attributes, got[0].Attributes[0].Data)
			}
		})
	}
}

func TestIntegrationServerComponentDelete(t *testing.T) {
	s := serverTest(t)

	var serverID uuid.UUID

	realClientTests(t, func(ctx context.Context, authToken string, _ int, expectError bool) error {
		s.Client.SetToken(authToken)

		var err error

		_, err = s.Client.DeleteServerComponents(ctx, serverID)

		if !expectError {
			return nil
		}

		return err
	})

	serverID, err := uuid.Parse(dbtools.FixtureNemo.ID)
	if err != nil {
		t.Fatal(err)
	}

	var testCases = []struct {
		testName         string
		serverID         uuid.UUID
		expectedError    bool
		errorMsg         string
		expectedResponse string
	}{
		{
			"unknown server UUID returns not found",
			uuid.New(),
			true,
			"",
			"resource not found",
		},
		{
			"server components removed",
			serverID,
			false,
			"",
			"resource deleted",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			resp, err := s.Client.DeleteServerComponents(context.TODO(), tt.serverID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Contains(t, err.Error(), tt.expectedResponse)

				return
			}

			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Contains(t, tt.expectedResponse, resp.Message)
		})
	}
}
