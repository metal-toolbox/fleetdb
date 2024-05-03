package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// BiosConfigSettingQuery defines values you can query BiosConfigSettings with. Empty strings are ignored.
type BiosConfigSettingQuery struct {
	Key   string `query:"key"`
	Value string `query:"value"`
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (ccs *BiosConfigSettingQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.BiosConfigSettingTableColumns.SettingsKey, ccs.Key)
	mods = appendOperatorQueryMod(mods, comparitor, models.BiosConfigSettingTableColumns.SettingsValue, ccs.Value)

	return mods
}
