package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// NicQuery defines values you can query Nic structs with. Empty strings are ignored.
type NicQuery struct {
	Vendor        string  `query:"vendor"`
	Model         string  `query:"model"`
	PortBandwidth []int64 `query:"port_bandwidth"`
	PortCount     []int64 `query:"port_count"`
	Count         []int64 `query:"count"`
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (ssku *NicQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuNicTableColumns.Vendor, ssku.Vendor)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuNicTableColumns.Model, ssku.Model)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuNicTableColumns.PortBandwidth, ssku.PortBandwidth)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuNicTableColumns.PortCount, ssku.PortCount)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuNicTableColumns.Count, ssku.Count)

	return mods
}
