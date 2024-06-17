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

func TestBiosConfigSetQuery(t *testing.T) {
	testBiosConfigSetQueryParams := BiosConfigSetListParams{
		Params: []BiosConfigSetQueryParams{
			{
				Set: BiosConfigSetQuery{ // Look for all RTX cards
					Components: []BiosConfigComponentQuery{
						{
							Name: "RTX",
						},
					},
				},
				LogicalOperator:    OperatorLogicalOR,
				ComparitorOperator: OperatorComparitorLike,
			},
			{
				Set: BiosConfigSetQuery{ // Look for any components with PCIE that are using 16 PCIE lanes
					Components: []BiosConfigComponentQuery{
						{
							Settings: []BiosConfigSettingQuery{
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
			Preload: true,
			OrderBy: "nothing",
		},
	}

	values, err := url.ParseQuery("params%5B0%5D%5Bset%5D%5Bcomponents%5D%5B0%5D%5Bname%5D=RTX&params%5B0%5D%5Blogical%5D=olt_or&params%5B0%5D%5Bcomparitor%5D=oct_like&params%5B1%5D%5Bset%5D%5Bcomponents%5D%5B0%5D%5Bsettings%5D%5B0%5D%5Bkey%5D=PCIE+Lanes&params%5B1%5D%5Bset%5D%5Bcomponents%5D%5B0%5D%5Bsettings%5D%5B0%5D%5Bvalue%5D=x16&params%5B1%5D%5Blogical%5D=olt_and&params%5B1%5D%5Bcomparitor%5D=oct_eq&page%5BLimit%5D=50&page%5BPage%5D=4&page%5BPreload%5D=1&page%5BCursor%5D=cursor&page%5BOrderBy%5D=nothing")
	require.NoError(t, err)

	testCases := []struct {
		testName string
		params   BiosConfigSetListParams
		expected url.Values
	}{
		{
			testName: "config set params: set query test 1",
			params:   testBiosConfigSetQueryParams,
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

func TestBiosConfigSetQueryMods(t *testing.T) {
	testBiosConfigSetQueryParams := BiosConfigSetListParams{
		Params: []BiosConfigSetQueryParams{
			{
				Set: BiosConfigSetQuery{ // Look for all RTX cards
					Components: []BiosConfigComponentQuery{
						{
							Name: "RTX",
						},
					},
				},
				LogicalOperator:    OperatorLogicalOR,
				ComparitorOperator: OperatorComparitorLike,
			},
			{
				Set: BiosConfigSetQuery{ // Look for any components with PCIE that are using 16 PCIE lanes
					Components: []BiosConfigComponentQuery{
						{
							Settings: []BiosConfigSettingQuery{
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
			Preload: true,
			OrderBy: "nothing",
		},
	}

	mods := []qm.QueryMod{}
	whereMods := []qm.QueryMod{}
	whereMods = appendOperatorQueryMod(whereMods, OperatorComparitorLike, models.BiosConfigComponentTableColumns.Name, testBiosConfigSetQueryParams.Params[0].Set.Components[0].Name)
	whereMods = []qm.QueryMod{
		qm.Or2(qm.Expr(whereMods...)),
	}
	mods = append(mods, whereMods...)

	whereMods = []qm.QueryMod{}
	whereMods = appendOperatorQueryMod(whereMods, OperatorComparitorEqual, models.BiosConfigSettingTableColumns.SettingsKey, testBiosConfigSetQueryParams.Params[1].Set.Components[0].Settings[0].Key)
	whereMods = appendOperatorQueryMod(whereMods, OperatorComparitorEqual, models.BiosConfigSettingTableColumns.SettingsValue, testBiosConfigSetQueryParams.Params[1].Set.Components[0].Settings[0].Value)
	whereMods = []qm.QueryMod{
		qm.Expr(whereMods...),
	}
	mods = append(mods, whereMods...)

	mods = append(mods, qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
		models.TableNames.BiosConfigComponents,
		models.BiosConfigSetTableColumns.ID,
		models.BiosConfigComponentTableColumns.FKBiosConfigSetID)))
	mods = append(mods, qm.InnerJoin(fmt.Sprintf("%s on %s = %s",
		models.TableNames.BiosConfigSettings,
		models.BiosConfigComponentTableColumns.ID,
		models.BiosConfigSettingTableColumns.FKBiosConfigComponentID)))

	mods = append(mods, testBiosConfigSetQueryParams.Pagination.queryMods()...)

	testCases := []struct {
		testName string
		params   BiosConfigSetListParams
		mods     []qm.QueryMod
	}{
		{
			testName: "config set params: set query mod test 1",
			params:   testBiosConfigSetQueryParams,
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
