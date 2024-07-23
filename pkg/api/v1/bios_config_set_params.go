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

// BiosConfigSetQuery defines values you can query BiosConfigSets with. Empty strings are ignored.
type BiosConfigSetQuery struct {
	Name       string                     `query:"name"`
	Version    string                     `query:"version"`
	Components []BiosConfigComponentQuery `query:"components"`
}

// BiosConfigSetQueryParams defines a BiosConfigSetQuery struct and operators you can use to query BiosConfigSets with. If LogicalOperator is an empty string, it will default to OperatorLogicalAND. If ComparitorOperator is an empty string, it will default to OperatorComparitorEqual
type BiosConfigSetQueryParams struct {
	Set                BiosConfigSetQuery     `query:"set"`
	LogicalOperator    OperatorLogicalType    `query:"logical"`
	ComparitorOperator OperatorComparitorType `query:"comparitor"`
}

// BiosConfigSetListParams params is an array of potential expressions when querying.
// Each one will have a Set. This Set will define values you want to search on, empty strings will be ignored.
// The ComparitorOperator will define how you want to compare those values.
// All values within a single BiosConfigSetQueryParams item will be grouped together and "AND"'ed.
// The LogicalOperator will define how that BiosConfigSetQueryParams item will be grouped with other BiosConfigSetQueryParams items.
// Note: You must set PaginationParams.Preload to load BiosConfigComponents and BiosConfigSettings.
type BiosConfigSetListParams struct {
	Params     []BiosConfigSetQueryParams `query:"params"`
	Pagination PaginationParams           `query:"page"`
}

// setQuery implements queryParams.
func (p *BiosConfigSetListParams) setQuery(q url.Values) {
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

func parseBiosConfigSetListParams(c *gin.Context) (*BiosConfigSetListParams, error) {
	params := BiosConfigSetListParams{}
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
func (p *BiosConfigSetListParams) queryMods() []qm.QueryMod {
	mods := []qm.QueryMod{qm.Distinct(fmt.Sprintf("\"%s\".*", models.TableNames.BiosConfigSets))}

	// Only INNER JOIN if we have query params for settings or components
	haveComponents := false
	haveSettings := false

	for i := range p.Params {
		whereMods := []qm.QueryMod{}

		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.BiosConfigSetTableColumns.Name, p.Params[i].Set.Name)
		whereMods = appendOperatorQueryMod(whereMods, p.Params[i].ComparitorOperator, models.BiosConfigSetTableColumns.Version, p.Params[i].Set.Version)

		for j := range p.Params[i].Set.Components {
			haveComponents = true
			if len(p.Params[i].Set.Components[j].Settings) > 0 {
				haveSettings = true
			}

			whereMods = append(whereMods, p.Params[i].Set.Components[j].queryMods(p.Params[i].ComparitorOperator)...)
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

	// Join Components table
	if haveComponents {
		mods = append(mods, qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
			models.TableNames.BiosConfigComponents,
			models.BiosConfigSetTableColumns.ID,
			models.BiosConfigComponentTableColumns.FKBiosConfigSetID)))
	}

	// Join Settings into Components
	if haveSettings {
		mods = append(mods, qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
			models.TableNames.BiosConfigSettings,
			models.BiosConfigComponentTableColumns.ID,
			models.BiosConfigSettingTableColumns.FKBiosConfigComponentID)))
	}

	mods = append(mods, p.Pagination.queryMods()...)

	return mods
}
