package fleetdbapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

func TestIntegrationServerCreateComponentsGetInventory(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, _ int, expectError bool) error {
		s.Client.SetToken(authToken)

		attrs, _, err := s.Client.GetComponents(ctx, uuid.MustParse(dbtools.FixtureNemo.ID), nil)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, attrs, 2)
		}

		return err
	})

	// init fixture data
	// 1. get list of servers
	servers, _, err := s.Client.List(context.Background(), nil)
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

	expectedAttribute1 := `{"id":"Enclosure.Internal.0-1","chassis_type":"StorageEnclosure","description":"PCIe SSD Backplane 1"}`
	expectedAttribute2 := `{"id":"NIC.Embedded.2-1-1","physid":"2"}`
	expectedVAStatus := `{"Health":"OK","State":"Enabled"}`

	metadataData := fmt.Sprintf("[%v,%v]", expectedAttribute1, expectedAttribute2)
	// `[{"chassis_type":"StorageEnclosure","description":"PCIe SSD Backplane 1","id":"Enclosure.Internal.0-1"},{"id":"NIC.Embedded.2-1-1","physid":"2"}]`
	vaStatus := fmt.Sprintf(`[
		{"status":%v},{"nic_port_status":{"Health":"OK","State":"Enabled","active_link_technology":
		 "Ethernet","id":"NIC.Slot.3-2-1","link_status":"Up","macaddress":"40:A6:B7:5C:E5:61"}},{"nic_port_status":{"Health":"OK","State":"Enabled","active_link_technology":"Ethernet","id":"NIC.Slot.3-1-1","link_status":"Up","macaddress":"40:A6:B7:5C:E5:60"}}
		]`, expectedVAStatus)

	// fixture to create a server components
	csFixtureCreate := fleetdbapi.ServerComponentSlice{
		{
			ServerUUID:        servers[1].UUID,
			Name:              "fakeName",
			Vendor:            "fakeVendor",
			Model:             "fakeModel",
			Serial:            "fakeSerial",
			ComponentTypeID:   componentTypeSlice.ByName("Fins").ID,
			ComponentTypeName: componentTypeSlice.ByName("Fins").Name,
			ComponentTypeSlug: componentTypeSlice.ByName("Fins").Slug,
			Attributes: []fleetdbapi.Attributes{
				{
					Namespace: "sh.hollow.alloy.outofband.metadata",
					Data:      json.RawMessage([]byte(metadataData)),
				},
			},
			VersionedAttributes: []fleetdbapi.VersionedAttributes{
				{
					Namespace: "sh.hollow.alloy.outofband.status",
					Data:      json.RawMessage([]byte(vaStatus)),
				},
			},
		},
	}

	// create server component
	_, err = s.Client.CreateComponents(context.TODO(), servers[1].UUID, csFixtureCreate)
	if err != nil {
		t.Fatal(err)
	}

	var gotAttribute1 bool
	var gotAttribute2 bool
	got, _, err := s.Client.GetServerInventory(context.TODO(), servers[1].UUID, false)
	for _, c := range got.Components {
		if c.Name == "fakeName" && c.Vendor == "fakeVendor" && c.Model == "fakeModel" && c.Serial == "fakeSerial" {
			gotVAStatus, _ := json.Marshal(c.Status)
			if string(gotVAStatus) != expectedVAStatus {
				t.Errorf("got VA status %v, expect %v", string(gotVAStatus), expectedVAStatus)
			}

			gotMetadataData, _ := json.Marshal(c.Attributes)
			fmt.Printf("string(gotMetadataData) = %v\n", string(gotMetadataData))
			if string(gotMetadataData) == expectedAttribute1 {
				gotAttribute1 = true
			}
			if string(gotMetadataData) == expectedAttribute2 {
				gotAttribute2 = true
			}
		}
	}
	if !gotAttribute1 {
		t.Errorf("failed to receive metadata %v", expectedAttribute1)
	}
	if !gotAttribute2 {
		t.Errorf("failed to receive metadata %v", expectedAttribute2)
	}
}
