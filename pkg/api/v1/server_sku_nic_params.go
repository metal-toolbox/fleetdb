package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerSkuNicQuery defines values you can query BiosConfigComponents with. Empty strings are ignored.
type ServerSkuNicQuery struct {
	PortBandwidth []int64 `query:"port_bandwidth"`
	PortCount     []int64 `query:"port_count"`
	Count         []int64 `query:"count"`
}

// QueryMods converts the list params into sql conditions that can be added to sql queries
func (ssku *ServerSkuNicQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuNicTableColumns.PortBandwidth, ssku.PortBandwidth)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuNicTableColumns.PortCount, ssku.PortCount)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuNicTableColumns.Count, ssku.Count)

	return mods
}
