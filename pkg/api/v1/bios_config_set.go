package fleetdbapi

import (
	// "bytes"
	"time"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// BiosConfigSet represents a BIOS Configuration Set
type BiosConfigSet struct {
	ID         string                `json:"id"`
	Name       string                `json:"name" binding:"required"`
	Version    string                `json:"version" binding:"required"`
	Components []BiosConfigComponent `json:"components" binding:"required"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

// toDBModelBiosConfigSet converts a BiosConfigSet into a models.BiosConfigSet
func (cs *BiosConfigSet) toDBModelBiosConfigSet() *models.BiosConfigSet {
	dbcs := &models.BiosConfigSet{
		Name:    cs.Name,
		Version: cs.Version,
		ID:      cs.ID,
	}

	return dbcs
}

// toDBModelBiosConfigSetDeep converts a BiosConfigSet into a models.BiosConfigSet. It also includes all relations, doing a deep copy
func (cs *BiosConfigSet) toDBModelBiosConfigSetDeep() *models.BiosConfigSet {
	dbcs := cs.toDBModelBiosConfigSet()

	if len(cs.Components) > 0 {
		dbcs.R = dbcs.R.NewStruct()

		for _, component := range cs.Components {
			dbcs.R.FKBiosConfigSetBiosConfigComponents = append(dbcs.R.FKBiosConfigSetBiosConfigComponents, component.toDBModelBiosConfigComponentDeep(dbcs))
		}
	}

	return dbcs
}

// fromDBModelBiosConfigSet converts a models.BiosConfigSet into a BiosConfigSet
func (cs *BiosConfigSet) fromDBModelBiosConfigSet(set *models.BiosConfigSet) {
	cs.ID = set.ID
	cs.Name = set.Name
	cs.Version = set.Version
	cs.CreatedAt = set.CreatedAt.Time
	cs.UpdatedAt = set.CreatedAt.Time

	if set.R != nil {
		cs.Components = make([]BiosConfigComponent, len(set.R.FKBiosConfigSetBiosConfigComponents))
		for i, dbComponent := range set.R.FKBiosConfigSetBiosConfigComponents {
			cs.Components[i].fromDBModelBiosConfigComponent(dbComponent)
		}
	}
}
