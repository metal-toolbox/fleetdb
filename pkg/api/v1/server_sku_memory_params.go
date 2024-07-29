package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// MemoryQuery defines values you can query Memories with. Empty strings are ignored.
type MemoryQuery struct {
	Vendor string  `query:"vendor"`
	Model  string  `query:"model"`
	Bytes  []int64 `query:"bytes"`
	Count  []int64 `query:"count"`
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (ssku *MemoryQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuMemoryTableColumns.Vendor, ssku.Vendor)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuMemoryTableColumns.Model, ssku.Model)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuMemoryTableColumns.Bytes, ssku.Bytes)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuMemoryTableColumns.Count, ssku.Count)

	return mods
}
