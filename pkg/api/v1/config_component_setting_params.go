package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ConfigComponentSettingQuery defines values you can query ConfigComponentSettings with. Empty strings are ignored.
type ConfigComponentSettingQuery struct {
	Key   string `query:"key"`
	Value string `query:"value"`
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (ccs *ConfigComponentSettingQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.ConfigComponentSettingTableColumns.SettingsKey, ccs.Key)
	mods = appendOperatorQueryMod(mods, comparitor, models.ConfigComponentSettingTableColumns.SettingsValue, ccs.Value)

	return mods
}
