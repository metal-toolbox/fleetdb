package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerSkuDiskQuery defines values you can query BiosConfigComponents with. Empty strings are ignored.
type ServerSkuDiskQuery struct {
	Bytes    []int64 `query:"bytes"`
	Protocol string  `query:"protocol"`
	Count    []int64 `query:"count"`
}

// QueryMods converts the list params into sql conditions that can be added to sql queries
func (ssku *ServerSkuDiskQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuDiskTableColumns.Bytes, ssku.Bytes)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuDiskTableColumns.Protocol, ssku.Protocol)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuDiskTableColumns.Count, ssku.Count)

	return mods
}
