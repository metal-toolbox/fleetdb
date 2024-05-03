package fleetdbapi

import (
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// BiosConfigSetting represents a BIOS Configuration Component Setting
type BiosConfigSetting struct {
	ID        string    `json:"id"`
	Key       string    `json:"key" binding:"required"`
	Value     string    `json:"value" binding:"required"`
	Raw       []byte    `json:"raw,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ccs *BiosConfigSetting) toDBModelBiosConfigSetting() *models.BiosConfigSetting {
	dbccs := &models.BiosConfigSetting{
		SettingsKey:   ccs.Key,
		SettingsValue: ccs.Value,
		Raw:           null.JSONFrom(ccs.Raw),
	}

	return dbccs
}

func (ccs *BiosConfigSetting) fromDBModelBiosConfigSetting(setting *models.BiosConfigSetting) {
	ccs.ID = setting.ID
	ccs.Key = setting.SettingsKey
	ccs.Value = setting.SettingsValue
	ccs.CreatedAt = setting.CreatedAt.Time
	ccs.UpdatedAt = setting.UpdatedAt.Time

	if !setting.Raw.IsZero() {
		ccs.Raw = setting.Raw.JSON
	}
}
