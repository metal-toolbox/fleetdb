package fleetdbapi

import (
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ConfigComponentSetting represents a Configuration Component Setting
type ConfigComponentSetting struct {
	ID        string    `json:"id"`
	Key       string    `json:"key" binding:"required"`
	Value     string    `json:"value" binding:"required"`
	Custom    []byte    `json:"custom,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ccs *ConfigComponentSetting) toDBModelConfigComponentSetting() *models.ConfigComponentSetting {
	dbccs := &models.ConfigComponentSetting{
		SettingsKey:   ccs.Key,
		SettingsValue: ccs.Value,
		Custom:        null.JSONFrom(ccs.Custom),
	}

	return dbccs
}

func (ccs *ConfigComponentSetting) fromDBModelConfigComponentSetting(setting *models.ConfigComponentSetting) {
	ccs.ID = setting.ID
	ccs.Key = setting.SettingsKey
	ccs.Value = setting.SettingsValue
	ccs.CreatedAt = setting.CreatedAt.Time
	ccs.UpdatedAt = setting.UpdatedAt.Time

	if !setting.Custom.IsZero() {
		ccs.Custom = setting.Custom.JSON
	}
}
