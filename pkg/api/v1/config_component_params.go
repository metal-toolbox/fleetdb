package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ConfigComponentQuery defines values you can query ConfigComponents with. Empty strings are ignored.
type ConfigComponentQuery struct {
	Name     string                        `query:"name"`
	Vendor   string                        `query:"vendor"`
	Model    string                        `query:"model"`
	Serial   string                        `query:"serial"`
	Settings []ConfigComponentSettingQuery `query:"settings"`
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (cc *ConfigComponentQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.ConfigComponentTableColumns.Name, cc.Name)
	mods = appendOperatorQueryMod(mods, comparitor, models.ConfigComponentTableColumns.Vendor, cc.Vendor)
	mods = appendOperatorQueryMod(mods, comparitor, models.ConfigComponentTableColumns.Model, cc.Model)
	mods = appendOperatorQueryMod(mods, comparitor, models.ConfigComponentTableColumns.Serial, cc.Serial)
	mods = appendOperatorQueryMod(mods, comparitor, models.ConfigComponentTableColumns.Vendor, cc.Vendor)

	for i := range cc.Settings {
		mods = append(mods, cc.Settings[i].queryMods(comparitor)...)
	}

	return mods
}
