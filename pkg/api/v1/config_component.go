package fleetdbapi

import (
	"time"

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
		Vendor: cc.Vendor,
		Model:  cc.Model,
		Serial: cc.Serial,
	}

	return dbcc
}

// fromDBModelConfigComponent converts a models.ConfigComponent (created by sqlboiler) into a ConfigComponent
func (cc *ConfigComponent) fromDBModelConfigComponent(component *models.ConfigComponent) {
	cc.ID = component.ID
	cc.Name = component.Name
	cc.Vendor = component.Vendor
	cc.Model = component.Model
	cc.Serial = component.Serial
	cc.CreatedAt = component.CreatedAt.Time
	cc.UpdatedAt = component.CreatedAt.Time

	if component.R != nil {
		cc.Settings = make([]ConfigComponentSetting, len(component.R.FKComponentConfigComponentSettings))
		for i, dbSetting := range component.R.FKComponentConfigComponentSettings {
			cc.Settings[i].fromDBModelConfigComponentSetting(dbSetting)
		}
	}
}
