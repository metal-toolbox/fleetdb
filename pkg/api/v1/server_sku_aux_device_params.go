package fleetdbapi

import (
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// AuxDeviceQuery defines values you can query Aux Devices with. Empty strings are ignored.
type AuxDeviceQuery struct {
	Vendor     string `query:"vendor"`
	Model      string `query:"model"`
	DeviceType string `query:"device_type"`
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (ssku *AuxDeviceQuery) queryMods(comparitor OperatorComparitorType) []qm.QueryMod {
	mods := []qm.QueryMod{}

	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuAuxDeviceTableColumns.Vendor, ssku.Vendor)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuAuxDeviceTableColumns.Model, ssku.Model)
	mods = appendOperatorQueryMod(mods, comparitor, models.ServerSkuAuxDeviceTableColumns.DeviceType, ssku.DeviceType)

	return mods
}
