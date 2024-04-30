package fleetdbapi

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

func TestConfigSetQuery(t *testing.T) {
	testConfigSetQueryParams := ConfigSetListParams{
		Params: []ConfigSetQueryParams{
			{
				Set: ConfigSetQuery{ // Look for all RTX cards
					Components: []ConfigComponentQuery{
						{
							Name: "RTX",
						},
					},
				},
				LogicalOperator:    OperatorLogicalOR,
				ComparitorOperator: OperatorComparitorLike,
			},
			{
				Set: ConfigSetQuery{ // Look for any components with PCIE that are using 16 PCIE lanes
					Components: []ConfigComponentQuery{
						{
							Settings: []ConfigComponentSettingQuery{
								{
									Key:   "PCIE Lanes",
									Value: "x16",
								},
							},
						},
					},
				},
				LogicalOperator:    OperatorLogicalAND,
				ComparitorOperator: OperatorComparitorEqual,
			},
		},
		Pagination: PaginationParams{
			Limit:   50,
			Page:    4,
			Cursor:  "cursor",
			Preload: false,
			OrderBy: "nothing",
		},
	}

	values, err := url.ParseQuery("params%5B0%5D%5Bset%5D%5Bcomponents%5D%5B0%5D%5Bname%5D=RTX&params%5B0%5D%5Blogical%5D=olt_or&params%5B0%5D%5Bcomparitor%5D=oct_like&params%5B1%5D%5Bset%5D%5Bcomponents%5D%5B0%5D%5Bsettings%5D%5B0%5D%5Bkey%5D=PCIE+Lanes&params%5B1%5D%5Bset%5D%5Bcomponents%5D%5B0%5D%5Bsettings%5D%5B0%5D%5Bvalue%5D=x16&params%5B1%5D%5Blogical%5D=olt_and&params%5B1%5D%5Bcomparitor%5D=oct_eq&pagination%5BLimit%5D=50&pagination%5BPage%5D=4&pagination%5BCursor%5D=cursor&pagination%5BOrderBy%5D=nothing")
	require.NoError(t, err)

	testCases := []struct {
		testName string
		params   ConfigSetListParams
		expected url.Values
	}{
		{
			testName: "config set params: set query test 1",
			params:   testConfigSetQueryParams,
			expected: values,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			newValues := make(url.Values)
			tc.params.setQuery(newValues)
			assert.Equal(t, tc.expected, newValues)
		})
	}
}

func TestConfigSetQueryMods(t *testing.T) {
	testConfigSetQueryParams := ConfigSetListParams{
		Params: []ConfigSetQueryParams{
			{
				Set: ConfigSetQuery{ // Look for all RTX cards
					Components: []ConfigComponentQuery{
						{
							Name: "RTX",
						},
					},
				},
				LogicalOperator:    OperatorLogicalOR,
				ComparitorOperator: OperatorComparitorLike,
			},
			{
				Set: ConfigSetQuery{ // Look for any components with PCIE that are using 16 PCIE lanes
					Components: []ConfigComponentQuery{
						{
							Settings: []ConfigComponentSettingQuery{
								{
									Key:   "PCIE Lanes",
									Value: "x16",
								},
							},
						},
					},
				},
				LogicalOperator:    OperatorLogicalAND,
				ComparitorOperator: OperatorComparitorEqual,
			},
		},
		Pagination: PaginationParams{
			Limit:   50,
			Page:    4,
			Cursor:  "cursor",
			Preload: false,
			OrderBy: "nothing",
		},
	}

	mods := []qm.QueryMod{}
	whereMods := []qm.QueryMod{}
	whereMods = appendOperatorQueryMod(whereMods, OperatorComparitorLike, models.ConfigComponentTableColumns.Name, testConfigSetQueryParams.Params[0].Set.Components[0].Name)
	whereMods = []qm.QueryMod{
		qm.Or2(qm.Expr(whereMods...)),
	}
	mods = append(mods, whereMods...)

	whereMods = []qm.QueryMod{}
	whereMods = appendOperatorQueryMod(whereMods, OperatorComparitorEqual, models.ConfigComponentSettingTableColumns.SettingsKey, testConfigSetQueryParams.Params[1].Set.Components[0].Settings[0].Key)
	whereMods = appendOperatorQueryMod(whereMods, OperatorComparitorEqual, models.ConfigComponentSettingTableColumns.SettingsValue, testConfigSetQueryParams.Params[1].Set.Components[0].Settings[0].Value)
	whereMods = []qm.QueryMod{
		qm.Expr(whereMods...),
	}
	mods = append(mods, whereMods...)

	mods = append(mods, qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
		models.TableNames.ConfigComponents,
		models.ConfigSetTableColumns.ID,
		models.ConfigComponentTableColumns.FKConfigSetID)))
	mods = append(mods, qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
		models.TableNames.ConfigComponentSettings,
		models.ConfigComponentTableColumns.ID,
		models.ConfigComponentSettingTableColumns.FKComponentID)))

	mods = append(mods, testConfigSetQueryParams.Pagination.queryMods()...)

	testCases := []struct {
		testName string
		params   ConfigSetListParams
		mods     []qm.QueryMod
	}{
		{
			testName: "config set params: set query mod test 1",
			params:   testConfigSetQueryParams,
			mods:     mods,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			newMods := tc.params.queryMods()
			assert.Equal(t, tc.mods, newMods)
		})
	}
}
