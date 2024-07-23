package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerSkuAuxDeviceQuery defines values you can query BiosConfigComponents with. Empty strings are ignored.
type ServerSkuAuxDeviceQuery struct {
	Vendor     string `query:"vendor"`
	Model      string `query:"model"`
	DeviceType string `query:"device_type"`
}

// QueryMods converts the list params into sql conditions that can be added to sql queries
func (ssku *ServerSkuAuxDeviceQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuAuxDeviceTableColumns.Vendor, ssku.Vendor)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuAuxDeviceTableColumns.Model, ssku.Model)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuAuxDeviceTableColumns.DeviceType, ssku.DeviceType)

	return mods
}
