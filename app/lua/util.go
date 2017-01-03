package lua

import (
	"reflect"
	"strconv"

	"github.com/yuin/gopher-lua"
	"net/url"
	"strings"
	"time"
)

// GetStructVariables loads all the global variables
// from a lua file into a struct using reflect
func GetStructVariables(src interface{}, L *lua.LState) error {
	v := reflect.ValueOf(src).Elem()

	// Loop all struct fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldTag := v.Type().Field(i)

		// If field contains the tag lua
		if t, ok := fieldTag.Tag.Lookup("lua"); ok {
			if t == "" {
				continue
			}

			// Get variable from the lua stack
			variable := L.GetGlobal(t)
			if variable.Type() == lua.LTNil {
				continue
			}

			// Determine what type of variable is and
			// set the field
			switch variable.Type() {

			// Variable is integer
			case lua.LTNumber:
				n, err := strconv.ParseInt(variable.String(), 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(n)

			// Variable is boolean
			case lua.LTBool:
				field.SetBool(lua.LVAsBool(variable))

			// Variable is string
			case lua.LTString:
				field.SetString(variable.String())
			}

		}
	}
	return nil
}

// MapToTable converts a Go map to a lua table
func MapToTable(m map[string]interface{}) *lua.LTable {
	// Main table pointer
	resultTable := &lua.LTable{}

	// Loop map
	for key, element := range m {

		switch element.(type) {
		case int64:
			resultTable.RawSetString(key, lua.LNumber(element.(int64)))
		case string:
			resultTable.RawSetString(key, lua.LString(element.(string)))
		case bool:
			resultTable.RawSetString(key, lua.LBool(element.(bool)))
		case map[string]interface{}:
			tble := MapToTable(element.(map[string]interface{}))
			resultTable.RawSetString(key, tble)
		}
	}

	return resultTable
}

// TableToMap converts a LUA table to a Go map[string]interface{}
func TableToMap(table *lua.LTable) map[string]interface{} {
	if table == nil {
		return nil
	}
	m := make(map[string]interface{})
	table.ForEach(func(i lua.LValue, v lua.LValue) {
		switch v.Type() {
		case lua.LTTable:
			n := TableToMap(v.(*lua.LTable))
			m[i.String()] = n
		case lua.LTNumber:
			num, err := strconv.ParseInt(v.String(), 10, 64)
			if err != nil {
				m[i.String()] = err.Error()
			} else {
				m[i.String()] = num
			}
		case lua.LTBool:
			b, err := strconv.ParseBool(v.String())
			if err != nil {
				m[i.String()] = err.Error()
			} else {
				m[i.String()] = b
			}
		default:
			m[i.String()] = v.String()
		}
	})
	return m
}

// QueryToTable converts a slice of interfaces to a lua table
func QueryToTable(r [][]interface{}, names []string) *lua.LTable {
	// Main table pointer
	resultTable := &lua.LTable{}

	// Loop query results
	for i := range r {

		// Table for current result set
		t := &lua.LTable{}

		// Loop result fields
		for x := range r[i] {

			// Set table fields
			v := r[i][x]
			switch v.(type) {
			case []uint8:
				t.RawSetString(names[x], lua.LString(string(r[i][x].([]uint8))))
			case time.Time:
				t.RawSetString(names[x], lua.LNumber(v.(time.Time).Unix()))
			}
		}

		// Append current table to main table
		resultTable.Append(t)
	}
	return resultTable
}

// TableToStringSlice converts a LUA table to a Go slice of strings
func TableToStringSlice(table *lua.LTable) []string {
	result := []string{}

	// Loop the lua table
	table.ForEach(func(i lua.LValue, v lua.LValue) {

		switch v.Type() {
		case lua.LTTable:

			// Convert table to slice using recursive algorithm
			r := TableToStringSlice(v.(*lua.LTable))

			// Append to main result
			result = append(result, r...)
		default:

			// Append to main result
			result = append(result, v.String())
		}
	})
	return result
}

// URLValuesToTable converts a map[string][]string to a LUA table
func URLValuesToTable(m url.Values) *lua.LTable {
	t := lua.LTable{}

	// Loop the map
	for i, v := range m {

		// Set the table fields
		t.RawSetString(
			i,
			lua.LString(strings.Join(v, ", ")),
		)
	}
	return &t
}