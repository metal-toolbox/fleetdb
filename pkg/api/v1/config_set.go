package fleetdbapi

import (
	// "bytes"
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ConfigSet represents a Configuration Set
type ConfigSet struct {
	ID         string            `json:"id"`
	Name       string            `json:"name" binding:"required"`
	Version    string            `json:"version" binding:"required"`
	Components []ConfigComponent `json:"components" binding:"required"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

// toDBModelConfigSet converts a ConfigSet into a models.ConfigSet (created by sqlboiler)
func (cs *ConfigSet) toDBModelConfigSet() *models.ConfigSet {
	dbcs := &models.ConfigSet{
		Name:    cs.Name,
		Version: null.StringFrom(cs.Version),
		ID:      cs.ID,
	}

	return dbcs
}

// fromDBModelConfigSet converts a models.ConfigSet (created by sqlboiler) into a ConfigSet
func (cs *ConfigSet) fromDBModelConfigSet(set *models.ConfigSet) error {
	cs.ID = set.ID
	cs.Name = set.Name
	cs.Version = set.Version.String
	cs.CreatedAt = set.CreatedAt.Time
	cs.UpdatedAt = set.CreatedAt.Time

	if set.R != nil {
		cs.Components = make([]ConfigComponent, len(set.R.FKConfigSetConfigComponents))
		for i, dbComponent := range set.R.FKConfigSetConfigComponents {
			cs.Components[i].fromDBModelConfigComponent(dbComponent)
		}
	}

	return nil
}
