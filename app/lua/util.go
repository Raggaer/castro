package lua

import (
	"reflect"
	"strconv"

	"github.com/yuin/gopher-lua"
	"net/url"
	"strings"
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

// TableToMap converts a LUA table to a Go map[interface{}]interface{}
func TableToMap(table *lua.LTable) map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	table.ForEach(func(i lua.LValue, v lua.LValue) {
		switch v.Type() {
		case lua.LTTable:
			n := TableToMap(v.(*lua.LTable))
			m[i.String()] = n
		default:
			m[i.String()] = v.String()
		}
	})
	return m
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