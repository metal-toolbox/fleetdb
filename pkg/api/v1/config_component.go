package fleetdbapi

import (
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ConfigComponent represents a Configuration Component
type ConfigComponent struct {
	ID        string                   `json:"id"`
	Name      string                   `json:"name" binding:"required"`
	Vendor    string                   `json:"vendor"`
	Model     string                   `json:"model"`
	Serial    string                   `json:"serial"`
	Settings  []ConfigComponentSetting `json:"settings" binding:"required"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

func (cc *ConfigComponent) toDBModelConfigComponent() *models.ConfigComponent {
	dbcc := &models.ConfigComponent{
		Name:   cc.Name,
		Vendor: null.StringFrom(cc.Vendor),
		Model:  null.StringFrom(cc.Model),
		Serial: null.StringFrom(cc.Serial),
	}

	return dbcc
}

// fromDBModelConfigComponent converts a models.ConfigComponent (created by sqlboiler) into a ConfigComponent
func (cc *ConfigComponent) fromDBModelConfigComponent(component *models.ConfigComponent) {
	cc.ID = component.ID
	cc.Name = component.Name
	cc.Vendor = component.Vendor.String
	cc.Model = component.Model.String
	cc.Serial = component.Serial.String
	cc.CreatedAt = component.CreatedAt.Time
	cc.UpdatedAt = component.CreatedAt.Time

	if component.R != nil {
		cc.Settings = make([]ConfigComponentSetting, len(component.R.FKComponentConfigComponentSettings))
		for i, dbSetting := range component.R.FKComponentConfigComponentSettings {
			cc.Settings[i].fromDBModelConfigComponentSetting(dbSetting)
		}
	}
}
