package fleetdbapi

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/hetiansu5/urlquery"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// ServerSkuQuery defines values you can query ServerSkus with. Empty strings are ignored. Empty Arrays are ignored.
type ServerSkuQuery struct {
	Name             string           `query:"name"`
	Version          string           `query:"version"`
	Vendor           string           `query:"vendor"`
	Chassis          string           `query:"chassis"`
	BMCModel         string           `query:"bmc_model"`
	MotherboardModel string           `query:"motherboard_model"`
	CPUVendor        string           `query:"cpu_vendor"`
	CPUModel         string           `query:"cpu_model"`
	CPUCores         []int64          `query:"cpu_cores"`
	CPUHertz         []int64          `query:"cpu_hertz"`
	CPUCount         []int64          `query:"cpu_count"`
	AuxDevices       []AuxDeviceQuery `query:"aux_devices"`
	Disks            []DiskQuery      `query:"disks"`
	Memory           []MemoryQuery    `query:"memory"`
	Nics             []NicQuery       `query:"nics"`
}

// ServerSkuQueryParams defines a ServerSkuQuery struct and operators you can use to query ServerSkus with. If LogicalOperator is an empty string, it will default to OperatorLogicalAND. If ComparitorOperator is an empty string, it will default to OperatorComparitorEqual
type ServerSkuQueryParams struct {
	Sku                ServerSkuQuery         `query:"sku"`
	LogicalOperator    OperatorLogicalType    `query:"logical"`
	ComparitorOperator OperatorComparitorType `query:"comparitor"`
}

// ServerSkuListParams params is an array of potential expressions when querying.
// Each one will have a Sku. This Sku will define values you want to search on, empty strings will be ignored, empty arrays will be ignored.
// The ComparitorOperator will define how you want to compare those values.
// All values within a single ServerSkuQueryParams item will be grouped together and "AND"'ed.
// The LogicalOperator will define how that ServerSkuQueryParams item will be grouped with other ServerSkuQueryParams items.
// Note: You must set PaginationParams.Preload to load AuxDevices, Disks, Memory, and Nics.
type ServerSkuListParams struct {
	Params     []ServerSkuQueryParams `query:"params"`
	Pagination PaginationParams       `query:"page"`
}

// setQuery implements the queryParams interface
func (p *ServerSkuListParams) setQuery(q url.Values) {
	if p == nil {
		return
	}

	encoder := urlquery.NewEncoder()
	encoder.RegisterEncodeFunc(reflect.String, OperatorURLQueryEncoder)

	bytes, err := encoder.Marshal(p)
	if err != nil {
		q.Set("error", err.Error())
	} else {
		newValues, err := url.ParseQuery(string(bytes))
		if err != nil {
			q.Set("error", err.Error())
		} else {
			for key, values := range newValues {
				for _, value := range values {
					q.Add(key, value)
				}
			}
		}
	}
}

// parseServerSkuListParams converts the queryURL to a ServerSkuListParams object
func parseServerSkuListParams(c *gin.Context) (*ServerSkuListParams, error) {
	params := ServerSkuListParams{}
	bytes := c.Request.URL.RawQuery

	parser := urlquery.NewParser()
	parser.RegisterDecodeFunc(reflect.String, OperatorURLQueryDecoder)

	err := parser.Unmarshal([]byte(bytes), &params)
	if err != nil {
		return nil, err
	}

	return &params, nil
}

// queryMods converts the list params into sql conditions that can be added to sql queries
func (p *ServerSkuListParams) queryMods() []qm.QueryMod {
	mods := []qm.QueryMod{qm.Distinct(fmt.Sprintf("\"%s\".*", models.TableNames.ServerSku))}

	// Only INNER JOIN if we have query params for it
	haveAuxDevice := false
	haveDisk := false
	haveMemory := false
	haveNic := false

	for i := range p.Params {
		whereMods := []qm.QueryMod{}

		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.ServerSkuTableColumns.Name, p.Params[i].Sku.Name)
		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.ServerSkuTableColumns.Version, p.Params[i].Sku.Version)
		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.ServerSkuTableColumns.Vendor, p.Params[i].Sku.Vendor)
		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.ServerSkuTableColumns.Chassis, p.Params[i].Sku.Chassis)
		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.ServerSkuTableColumns.BMCModel, p.Params[i].Sku.BMCModel)
		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.ServerSkuTableColumns.MotherboardModel, p.Params[i].Sku.MotherboardModel)
		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.ServerSkuTableColumns.CPUVendor, p.Params[i].Sku.CPUVendor)
		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.ServerSkuTableColumns.CPUModel, p.Params[i].Sku.CPUModel)
		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.ServerSkuTableColumns.CPUCores, p.Params[i].Sku.CPUCores)
		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.ServerSkuTableColumns.CPUHertz, p.Params[i].Sku.CPUHertz)
		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.ServerSkuTableColumns.CPUCount, p.Params[i].Sku.CPUCount)

		for j := range p.Params[i].Sku.AuxDevices {
			haveAuxDevice = true

			whereMods = append(whereMods, p.Params[i].Sku.AuxDevices[j].queryMods(p.Params[i].ComparitorOperator)...)
		}

		for j := range p.Params[i].Sku.Disks {
			haveDisk = true

			whereMods = append(whereMods, p.Params[i].Sku.Disks[j].queryMods(p.Params[i].ComparitorOperator)...)
		}

		for j := range p.Params[i].Sku.Memory {
			haveMemory = true

			whereMods = append(whereMods, p.Params[i].Sku.Memory[j].queryMods(p.Params[i].ComparitorOperator)...)
		}

		for j := range p.Params[i].Sku.Nics {
			haveNic = true

			whereMods = append(whereMods, p.Params[i].Sku.Nics[j].queryMods(p.Params[i].ComparitorOperator)...)
		}

		if len(whereMods) > 0 {
			if p.Params[i].LogicalOperator == OperatorLogicalOR {
				whereMods = []qm.QueryMod{
					qm.Or2(qm.Expr(whereMods...)),
				}
			} else {
				whereMods = []qm.QueryMod{
					qm.Expr(whereMods...),
				}
			}

			// We do these in separate chunks since qm.Expr() can only be run on qm.WhereQueryMod{}.
			// And each loop will then be a bunch of groups of qm.exprMod{} and not qm.WhereQueryMod{}.
			mods = append(mods, whereMods...)
		}
	}

	// Join AuxDevice table
	if haveAuxDevice {
		mods = append(mods, qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
			models.TableNames.ServerSkuAuxDevice,
			models.ServerSkuTableColumns.ID,
			models.ServerSkuAuxDeviceTableColumns.SkuID)))
	}

	// Join Disk table
	if haveDisk {
		mods = append(mods, qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
			models.TableNames.ServerSkuDisk,
			models.ServerSkuTableColumns.ID,
			models.ServerSkuDiskTableColumns.SkuID)))
	}

	// Join Memory table
	if haveMemory {
		mods = append(mods, qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
			models.TableNames.ServerSkuMemory,
			models.ServerSkuTableColumns.ID,
			models.ServerSkuMemoryTableColumns.SkuID)))
	}

	// Join Nic table
	if haveNic {
		mods = append(mods, qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
			models.TableNames.ServerSkuNic,
			models.ServerSkuTableColumns.ID,
			models.ServerSkuNicTableColumns.SkuID)))
	}

	mods = append(mods, p.Pagination.queryMods()...)

	return mods
}
