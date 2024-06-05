package fleetdbapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

func TestIntegrationFirmwareList(t *testing.T) {
	s := serverTest(t)

	scopes := []string{"read:server-component-firmwares", "write:server-component-firmwares"}
	scopedRealClientTests(t, scopes, func(ctx context.Context, authToken string, _ int, expectError bool) error {
		s.Client.SetToken(authToken)

		params := fleetdbapi.ComponentFirmwareVersionListParams{
			Vendor:    "",
			Model:     nil,
			Version:   "",
			Component: "",
		}

		r, resp, err := s.Client.ListServerComponentFirmware(ctx, &params)
		if !expectError {
			require.NoError(t, err)
			assert.Len(t, r, 7)
			assert.EqualValues(t, 7, resp.PageCount)
			assert.EqualValues(t, 1, resp.TotalPages)
			assert.EqualValues(t, 7, resp.TotalRecordCount)
			// We returned everything, so we shouldnt have a next page info
			assert.Nil(t, resp.Links.Next)
			assert.Nil(t, resp.Links.Previous)
		}

		return err
	})

	var testCases = []struct {
		testName      string
		params        *fleetdbapi.ComponentFirmwareVersionListParams
		expectedUUIDs []string
		expectError   bool
		errorMsg      string
	}{
		{
			"search by vendor",
			&fleetdbapi.ComponentFirmwareVersionListParams{
				Vendor: "Dell",
			},
			[]string{
				dbtools.FixtureDellR640BMC.ID,
				dbtools.FixtureDellR640BIOS.ID,
				dbtools.FixtureDellR6515BMC.ID,
				dbtools.FixtureDellR6515BIOS.ID,
				dbtools.FixtureDellR640CPLD.ID,
			},
			false,
			"",
		},
		{
			"search by model",
			&fleetdbapi.ComponentFirmwareVersionListParams{
				Model: []string{"X11DPH-T"},
			},
			[]string{dbtools.FixtureSuperMicroX11DPHTBMC.ID},
			false,
			"",
		},
		{
			"search by version",
			&fleetdbapi.ComponentFirmwareVersionListParams{
				Version: "2.6.6",
			},
			[]string{dbtools.FixtureDellR6515BIOS.ID},
			false,
			"",
		},
		{
			"search by component",
			&fleetdbapi.ComponentFirmwareVersionListParams{
				Component: "bios",
			},
			[]string{dbtools.FixtureDellR6515BIOS.ID, dbtools.FixtureDellR640BIOS.ID},
			false,
			"",
		},
		{
			"search by non-exist component",
			&fleetdbapi.ComponentFirmwareVersionListParams{
				Component: "non-exist component",
			},
			[]string{},
			false,
			"",
		},
		{
			"search by filename",
			&fleetdbapi.ComponentFirmwareVersionListParams{
				Filename: "BIOS_C4FT0_WN64_2.6.6.EXE",
			},
			[]string{dbtools.FixtureDellR6515BIOS.ID},
			false,
			"",
		},
		{
			"search by checksum",
			&fleetdbapi.ComponentFirmwareVersionListParams{
				Checksum: "1ddcb3c3d0fc5925ef03a3dde768e9e245c579039dd958fc0f3a9c6368b6c5f4",
			},
			[]string{dbtools.FixtureDellR6515BIOS.ID},
			false,
			"",
		},
		{
			"limit results",
			&fleetdbapi.ComponentFirmwareVersionListParams{
				Vendor:     "Dell",
				Pagination: &fleetdbapi.PaginationParams{Limit: 1},
			},
			[]string{dbtools.FixtureDellR640BIOS.ID},
			false,
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			r, _, err := s.Client.ListServerComponentFirmware(context.TODO(), tt.params)
			if tt.expectError {
				assert.NoError(t, err)
				return
			}

			var actual []string

			for _, srv := range r {
				actual = append(actual, srv.UUID.String())
			}

			assert.ElementsMatch(t, tt.expectedUUIDs, actual)
		})
	}
}

func TestIntegrationFirmwareGet(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, _ int, expectError bool) error {
		s.Client.SetToken(authToken)
		fw, _, err := s.Client.GetServerComponentFirmware(ctx, uuid.MustParse(dbtools.FixtureDellR640BMC.ID))

		if !expectError {
			require.NoError(t, err)
			assert.Equal(t, fw.UUID, uuid.MustParse(dbtools.FixtureDellR640BMC.ID))
		}

		return err
	})
}

func TestIntegrationServerComponentFirmwareCreate(t *testing.T) {
	s := serverTest(t)

	var inbandFalse, oemFalse bool
	inbandTrue := true
	oemTrue := true
	realClientTests(t, func(ctx context.Context, authToken string, _ int, expectError bool) error {
		s.Client.SetToken(authToken)
		testFirmware := fleetdbapi.ComponentFirmwareVersion{
			UUID:          uuid.New(),
			Vendor:        "dell",
			Model:         []string{"r615"},
			Filename:      "foobar",
			Version:       "21.07.00",
			Component:     "system",
			Checksum:      "foobar",
			UpstreamURL:   "https://vendor.com/firmwares/DSU_21.07.00/",
			RepositoryURL: "http://example-firmware-bucket.s3.amazonaws.com/firmware/dell/DSU_21.07.00/",
			InstallInband: &inbandFalse,
			OEM:           &oemFalse,
		}

		id, resp, err := s.Client.CreateServerComponentFirmware(ctx, testFirmware)
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, id)
			assert.Equal(t, testFirmware.UUID.String(), id.String())
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/server-component-firmwares/%s", id), resp.Links.Self.Href)
		}

		return err
	})

	var testCases = []struct {
		testName         string
		firmware         *fleetdbapi.ComponentFirmwareVersion
		expectedError    bool
		expectedResponse string
		errorMsg         string
	}{
		{
			"empty required parameters",
			&fleetdbapi.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "dell",
				Model:         nil,
				Filename:      "foobar",
				Version:       "12345",
				Component:     "bios",
				Checksum:      "foobar",
				UpstreamURL:   "https://vendor.com/firmware-file",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/foobar",
				InstallInband: &inbandFalse,
			},
			true,
			"400",
			"Error:Field validation for 'Model' failed on the 'required' tag",
		},
		{
			"required lowercase parameters",
			&fleetdbapi.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "DELL",
				Model:         []string{"r615"},
				Filename:      "foobar",
				Version:       "12345",
				Component:     "bios",
				Checksum:      "foobar",
				UpstreamURL:   "https://vendor.com/firmware-file",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/foobar",
				InstallInband: &inbandFalse,
			},
			true,
			"400",
			"Error:Field validation for 'Vendor' failed on the 'lowercase' tag",
		},
		{
			"required installInband parameter",
			&fleetdbapi.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "DELL",
				Model:         []string{"r615"},
				Filename:      "foobar",
				Version:       "12345",
				Component:     "bios",
				Checksum:      "foobar",
				UpstreamURL:   "https://vendor.com/firmware-file",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/foobar",
			},
			true,
			"400",
			"Error:Field validation for 'InstallInband' failed on the 'required' tag",
		},
		{
			"required OEM parameter",
			&fleetdbapi.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "DELL",
				Model:         []string{"r615"},
				Filename:      "foobar",
				Version:       "12345",
				Component:     "bios",
				Checksum:      "foobar",
				UpstreamURL:   "https://vendor.com/firmware-file",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/foobar",
				InstallInband: &inbandFalse,
			},
			true,
			"400",
			"Error:Field validation for 'OEM' failed on the 'required' tag",
		},
		{
			"filename allowed to be mixed case",
			&fleetdbapi.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "dell",
				Model:         []string{"r615"},
				Filename:      "fooBAR",
				Version:       "12345",
				Component:     "bios",
				Checksum:      "foobar",
				UpstreamURL:   "https://vendor.com/firmware-file",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/foobar",
				InstallInband: &inbandFalse,
				OEM:           &oemFalse,
			},
			false,
			"200",
			"",
		},
		{
			"install inband can be set to true",
			&fleetdbapi.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "intel",
				Model:         []string{"e751"},
				Filename:      "fooBAR",
				Version:       "001",
				Component:     "nic",
				Checksum:      "blah",
				UpstreamURL:   "https://vendor.com/blob",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/blob",
				InstallInband: &inbandTrue,
				OEM:           &oemFalse,
			},
			false,
			"200",
			"",
		},
		{
			"oem can be set to true",
			&fleetdbapi.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "intel",
				Model:         []string{"e751"},
				Filename:      "fooBAR",
				Version:       "002",
				Component:     "nic",
				Checksum:      "blah",
				UpstreamURL:   "https://vendor.com/blob",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/blob",
				InstallInband: &inbandTrue,
				OEM:           &oemTrue,
			},
			false,
			"200",
			"",
		},
		{
			"duplicate vendor/component/version/filename not allowed",
			&fleetdbapi.ComponentFirmwareVersion{
				UUID:          uuid.New(),
				Vendor:        "dell",
				Model:         []string{"r615"},
				Filename:      "fooBAR",
				Version:       "12345",
				Component:     "bios",
				Checksum:      "foobar",
				UpstreamURL:   "https://vendor.com/firmware-file",
				RepositoryURL: "https://example-bucket.s3.awsamazon.com/foobar",
				InstallInband: &inbandFalse,
				OEM:           &oemFalse,
			},
			true,
			"400",
			"unable to insert into component_firmware_version: pq: duplicate key value violates unique constraint \"vendor_component_version_filename_unique\"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			fwUUID, _, err := s.Client.CreateServerComponentFirmware(context.TODO(), *tt.firmware)
			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Contains(t, err.Error(), tt.expectedResponse)

				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.firmware.UUID.String(), fwUUID.String())
		})
	}
}

func TestIntegrationServerComponentFirmwareDelete(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, _ int, _ bool) error {
		s.Client.SetToken(authToken)
		_, err := s.Client.DeleteServerComponentFirmware(ctx, fleetdbapi.ComponentFirmwareVersion{UUID: uuid.MustParse(dbtools.FixtureDellR640CPLD.ID)})

		return err
	})

	_, err := s.Client.DeleteServerComponentFirmware(context.TODO(), fleetdbapi.ComponentFirmwareVersion{UUID: uuid.MustParse(dbtools.FixtureDellR640BMC.ID)})
	assert.Contains(t, err.Error(), "violates foreign key constraint \"fk_firmware_id_ref_component_firmware_version\"")
}

func TestIntegrationServerComponentFirmwareUpdate(t *testing.T) {
	s := serverTest(t)

	realClientTests(t, func(ctx context.Context, authToken string, _ int, expectError bool) error {
		s.Client.SetToken(authToken)

		inband := true
		oem := true
		fw := fleetdbapi.ComponentFirmwareVersion{
			UUID:          uuid.MustParse(dbtools.FixtureDellR640BMC.ID),
			Vendor:        "dell",
			Model:         []string{"r615"},
			Filename:      "foobarino",
			Version:       "21.07.00",
			Component:     "bios",
			Checksum:      "foobar",
			UpstreamURL:   "https://vendor.com/firmware-file",
			RepositoryURL: "https://example-firmware-bucket.s3.amazonaws.com/firmware/dell/r615/bios/filename.ext",
			InstallInband: &inband,
			OEM:           &oem,
		}

		resp, err := s.Client.UpdateServerComponentFirmware(ctx, uuid.MustParse(dbtools.FixtureDellR640BMC.ID), fw)
		if !expectError {
			require.NoError(t, err)
			assert.NotNil(t, resp.Links.Self)
			assert.Equal(t, fmt.Sprintf("http://test.hollow.com/api/v1/server-component-firmwares/%s", dbtools.FixtureDellR640BMC.ID), resp.Links.Self.Href)
			fw, _, _ := s.Client.GetServerComponentFirmware(ctx, uuid.MustParse(dbtools.FixtureDellR640BMC.ID))
			assert.Equal(t, "foobarino", fw.Filename)
		}

		return err
	})
}
