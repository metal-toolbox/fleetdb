package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// BiosConfigComponentQuery defines values you can query BiosConfigComponents with. Empty strings are ignored.
type BiosConfigComponentQuery struct {
	Name     string                   `query:"name"`
	Vendor   string                   `query:"vendor"`
	Model    string                   `query:"model"`
	Settings []BiosConfigSettingQuery `query:"settings"`
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (cc *BiosConfigComponentQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.BiosConfigComponentTableColumns.Name, cc.Name)
	mods = appendOperatorQueryMod(mods, comparitor, models.BiosConfigComponentTableColumns.Vendor, cc.Vendor)
	mods = appendOperatorQueryMod(mods, comparitor, models.BiosConfigComponentTableColumns.Model, cc.Model)
	mods = appendOperatorQueryMod(mods, comparitor, models.BiosConfigComponentTableColumns.Vendor, cc.Vendor)

	for i := range cc.Settings {
		mods = append(mods, cc.Settings[i].queryMods(comparitor)...)
	}

	return mods
}
