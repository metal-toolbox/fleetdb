package fleetdbapi

import (
	"time"

	"github.com/google/uuid"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ComponentFirmwareVersion represents a firmware file
type ComponentFirmwareVersion struct {
	UUID          uuid.UUID `json:"uuid"`
	Vendor        string    `json:"vendor" binding:"required,lowercase"`
	Model         []string  `json:"model" binding:"required"`
	Filename      string    `json:"filename" binding:"required"`
	Version       string    `json:"version" binding:"required"`
	Component     string    `json:"component" binding:"required,lowercase"`
	Checksum      string    `json:"checksum" binding:"required,lowercase"`
	UpstreamURL   string    `json:"upstream_url" binding:"required"`
	RepositoryURL string    `json:"repository_url" binding:"required"`
	// The client has to always explicitly set this to true or false
	// for this to work with the validator, it needs to be a bool.
	InstallInband *bool     `json:"install_inband" binding:"required"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (f *ComponentFirmwareVersion) fromDBModel(dbF *models.ComponentFirmwareVersion) error {
	var err error

	f.UUID, err = uuid.Parse(dbF.ID)
	if err != nil {
		return err
	}

	f.Component = dbF.Component
	f.Vendor = dbF.Vendor
	f.Model = dbF.Model
	f.Filename = dbF.Filename
	f.Version = dbF.Version
	f.Checksum = dbF.Checksum
	f.UpstreamURL = dbF.UpstreamURL
	f.RepositoryURL = dbF.RepositoryURL
	f.CreatedAt = dbF.CreatedAt.Time
	f.UpdatedAt = dbF.UpdatedAt.Time
	f.InstallInband = &dbF.InstallInband

	return nil
}

func (f *ComponentFirmwareVersion) toDBModel() (*models.ComponentFirmwareVersion, error) {
	var installInband bool
	if f.InstallInband != nil {
		installInband = *f.InstallInband
	}

	dbF := &models.ComponentFirmwareVersion{
		Component:     f.Component,
		Vendor:        f.Vendor,
		Model:         f.Model,
		Filename:      f.Filename,
		Version:       f.Version,
		Checksum:      f.Checksum,
		UpstreamURL:   f.UpstreamURL,
		RepositoryURL: f.RepositoryURL,
		InstallInband: installInband,
	}

	if f.UUID.String() != uuid.Nil.String() {
		dbF.ID = f.UUID.String()
	}

	return dbF, nil
}
