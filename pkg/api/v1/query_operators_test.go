package fleetdbapi

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func TestQueryMod(t *testing.T) {
	testCases := []struct {
		testName   string
		comparitor OperatorComparitorType
		name       string
		value      string
		expected   qm.QueryMod
	}{
		{
			testName:   "query operators: query mod test 1",
			comparitor: OperatorComparitorLike,
			name:       "column1",
			value:      "value1",
			expected:   qm.Where("column1 LIKE ?", "value1%"),
		},
		{
			testName:   "query operators: query mod test 2",
			comparitor: OperatorComparitorLike,
			name:       "column2",
			value:      "value2%",
			expected:   qm.Where("column2 LIKE ?", "value2%"),
		},
		{
			testName:   "query operators: query mod test 3",
			comparitor: OperatorComparitorLike,
			name:       "column3",
			value:      "%value3",
			expected:   qm.Where("column3 LIKE ?", "%value3"),
		},
		{
			testName:   "query operators: query mod test 4",
			comparitor: OperatorComparitorNotEqual,
			name:       "column4",
			value:      "value4",
			expected:   qm.Where("column4 != ?", "value4"),
		},
		{
			testName:   "query operators: query mod test 5",
			comparitor: OperatorComparitorEqual,
			name:       "column5",
			value:      "value5",
			expected:   qm.Where("column5 = ?", "value5"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mods := []qm.QueryMod{}
			mods = appendOperatorQueryMod(mods, tc.comparitor, tc.name, tc.value)
			assert.Equal(t, tc.expected, mods[0])
		})
	}
}

func TestQueryURLEncoder(t *testing.T) {
	testCases := []struct {
		testName string
		value    reflect.Value
		expected string
	}{
		{
			testName: "query operators: query encoder test 1",
			value:    reflect.ValueOf(OperatorLogicalOR),
			expected: "olt_or",
		},
		{
			testName: "query operators: query encoder test 2",
			value:    reflect.ValueOf("olt_" + OperatorLogicalOR),
			expected: "olt_olt_or",
		},
		{
			testName: "query operators: query encoder test 3",
			value:    reflect.ValueOf(OperatorComparitorEqual),
			expected: "oct_eq",
		},
		{
			testName: "query operators: query encoder test 3",
			value:    reflect.ValueOf("oct_" + OperatorComparitorEqual),
			expected: "oct_oct_eq",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			s := OperatorURLQueryEncoder(tc.value)
			assert.Equal(t, tc.expected, s)
		})
	}
}

func TestQueryURLDecoder(t *testing.T) {
	testCases := []struct {
		testName string
		value    string
		expected reflect.Value
	}{
		{
			testName: "query operators: query decoder test 1",
			value:    "olt_or",
			expected: reflect.ValueOf(OperatorLogicalOR),
		},
		{
			testName: "query operators: query decoder test 2",
			value:    "olt_olt_or",
			expected: reflect.ValueOf("olt_or"),
		},
		{
			testName: "query operators: query decoder test 3",
			value:    "oct_eq",
			expected: reflect.ValueOf(OperatorComparitorEqual),
		},
		{
			testName: "query operators: query decoder test 4",
			value:    "oct_oct_eq",
			expected: reflect.ValueOf("oct_eq"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			s, err := OperatorURLQueryDecoder(tc.value)
			require.NoError(t, err)
			assert.Equal(t, tc.expected.Type(), s.Type())
			assert.Equal(t, tc.expected.String(), s.String())
		})
	}
}
