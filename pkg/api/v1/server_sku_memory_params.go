package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerSkuMemoryQuery defines values you can query BiosConfigComponents with. Empty strings are ignored.
type ServerSkuMemoryQuery struct {
	Bytes []int64 `query:"bytes"`
	Count []int64 `query:"count"`
}

// QueryMods converts the list params into sql conditions that can be added to sql queries
func (ssku *ServerSkuMemoryQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuMemoryTableColumns.Bytes, ssku.Bytes)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuMemoryTableColumns.Count, ssku.Count)

	return mods
}
