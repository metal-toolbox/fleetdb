package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// DiskQuery defines values you can query Disks with. Empty strings are ignored.
type DiskQuery struct {
	Vendor   string  `query:"vendor"`
	Model    string  `query:"model"`
	Bytes    []int64 `query:"bytes"`
	Protocol string  `query:"protocol"`
	Count    []int64 `query:"count"`
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (ssku *DiskQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuDiskTableColumns.Vendor, ssku.Vendor)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuDiskTableColumns.Model, ssku.Model)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuDiskTableColumns.Bytes, ssku.Bytes)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuDiskTableColumns.Protocol, ssku.Protocol)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuDiskTableColumns.Count, ssku.Count)

	return mods
}
