package fleetdbapi

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// OperatorComparitorType is used to control what kind of search is performed for any query. (preferrable option being UINTS)
type OperatorComparitorType string

// OperatorLogicalType is used to define how to group a query with other queries. Ex: (query expression) AND (query expression) OR (query expression)
type OperatorLogicalType string

// TODO; Should these just be ints?
const (
	// OperatorLogicalOR informs the SQL Builder to use OR when adding the param to the SQL Query. Making it Inclusive
	OperatorLogicalOR OperatorLogicalType = "or"
	// OperatorLogicalAND informs the SQL Builder to use AND when adding the param to the SQL Query. Making it Explicitly Inclusive
	OperatorLogicalAND = "and"
)

// TODO; We really should just make these the same values as qm.operators, or they should just be ints (preferrable option being UINTS)
const (
	// OperatorComparitorEqual means the value has to match the keys exactly
	OperatorComparitorEqual OperatorComparitorType = "eq"
	// OperatorComparitorNotEqual means the value has to match the keys exactly, but then exclude those matches
	OperatorComparitorNotEqual = "!="
	// OperatorComparitorLike allows you to pass in a value with % in it and match anything like it. If your string has no % in it one will be added to the end automatically
	OperatorComparitorLike = "like"
	// OperatorComparitorGreaterThan will convert the value at the given key to an int and return results that are greater than Value
	OperatorComparitorGreaterThan = "gt"
	// OperatorComparitorLessThan will convert the value at the given key to an int and return results that are less than Value
	OperatorComparitorLessThan = "lt"
)

// appendOperatorQueryMod is a helper function to build qm.QueryMods.
func appendOperatorQueryMod(mods []qm.QueryMod, comparitor OperatorComparitorType, name, value string) []qm.QueryMod {
	if value != "" {
		switch comparitor {
		case OperatorComparitorLike:
			if !strings.Contains(value, "%") {
				value = fmt.Sprintf("%s%%", value)
			}

			mod := qm.Where(fmt.Sprintf("%s LIKE ?", name), value)
			mods = append(mods, mod)
		case OperatorComparitorNotEqual:
			mod := qm.Where(fmt.Sprintf("%s != ?", name), value)
			mods = append(mods, mod)
		case OperatorComparitorEqual:
			fallthrough
		default:
			mod := qm.Where(fmt.Sprintf("%s = ?", name), value)
			mods = append(mods, mod)
		}
	}

	return mods
}

// OperatorURLQueryEncoder will be passed to a urlquery encoder to escape Operator types from query strings
// TODO; If we swap OperatorComparitorType and OperatorLogicalType to ints, this function will not be needed.
// NOTE: reflect.Set() doesnt convert custom string types back and forth like it is able to with custom int types
// So we have to escape and workaround these custom string types by escaping them when we parse the query.
func OperatorURLQueryEncoder(rv reflect.Value) string {
	switch rv.Type() {
	case reflect.TypeOf(OperatorComparitorType("")):
		return "oct_" + rv.String()
	case reflect.TypeOf(OperatorLogicalType("")):
		return "olt_" + rv.String()
	default:
		return rv.String()
	}
}

// OperatorURLQueryDecoder will be passed to a urlquery decoder to escape Operator types from query strings
// TODO; If we swap OperatorComparitorType and OperatorLogicalType to ints, this function will not be needed.
// reflect.Set() doesnt convert custom string types back and forth like it is able to with custom int types
// So we have escaped them and to get the values back we parse out the escape value and see if the strings are operator constants
func OperatorURLQueryDecoder(s string) (reflect.Value, error) {
	if strings.HasPrefix(s, "oct_") {
		s = strings.Replace(s, "oct_", "", 1)
	} else if strings.HasPrefix(s, "olt_") {
		s = strings.Replace(s, "olt_", "", 1)
	}

	switch s {
	case string(OperatorComparitorEqual):
		fallthrough
	case string(OperatorComparitorNotEqual):
		fallthrough
	case string(OperatorComparitorLike):
		fallthrough
	case string(OperatorComparitorGreaterThan):
		fallthrough
	case string(OperatorComparitorLessThan):
		return reflect.ValueOf(OperatorComparitorType(s)), nil
	case string(OperatorLogicalOR):
		fallthrough
	case string(OperatorLogicalAND):
		return reflect.ValueOf(OperatorLogicalType(s)), nil
	default:
		return reflect.ValueOf(s), nil
	}
}
