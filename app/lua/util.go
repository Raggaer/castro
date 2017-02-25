package lua

import (
	"github.com/yuin/gopher-lua"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"log"
)

type castroInterfaceField struct {
	Name string
	Type reflect.Kind
}

// GetStructVariables loads all the global variables from a lua file into a struct using reflect
func GetStructVariables(src interface{}, L *lua.LState) error {
	// Get base element
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

		log.Println(key)

		switch element.(type) {
		case float64:
			resultTable.RawSetString(key, lua.LNumber(element.(float64)))
		case int64:
			resultTable.RawSetString(key, lua.LNumber(element.(int64)))
		case string:
			resultTable.RawSetString(key, lua.LString(element.(string)))
		case bool:
			resultTable.RawSetString(key, lua.LBool(element.(bool)))
		case []byte:
			resultTable.RawSetString(key, lua.LString(string(element.([]byte))))
		case map[string]interface{}:

			// Get table from map
			tble := MapToTable(element.(map[string]interface{}))

			log.Println(key)

			resultTable.RawSetString(key, tble)

		case time.Time:
			resultTable.RawSetString(key, lua.LNumber(element.(time.Time).Unix()))

		case []map[string]interface{}:

			// Create slice table
			sliceTable := &lua.LTable{}

			// Loop element
			for _, s := range element.([]map[string]interface{}) {

				// Get table from map
				tble := MapToTable(s)

				sliceTable.Append(tble)
			}

			// Set slice table
			resultTable.RawSetString(key, sliceTable)

		case []interface{}:

			// Create slice table
			sliceTable := &lua.LTable{}

			// Loop interface slice
			for _, s := range element.([]interface{}) {

				// Switch interface type
				switch s.(type) {
				case map[string]interface{}:

					// Convert map to table
					t := MapToTable(s.(map[string]interface{}))

					// Append result
					sliceTable.Append(t)
				}
			}

			// Append to main table
			resultTable.RawSetString(key, sliceTable)
		}
	}

	return resultTable
}

// TableToMap converts a LUA table to a Go map[string]interface{}
func TableToMap(table *lua.LTable) map[string]interface{} {
	// Check for valid table
	if table == nil {
		return map[string]interface{}{}
	}

	// Data holder
	m := make(map[string]interface{})

	// Loop lua table
	table.ForEach(func(i lua.LValue, v lua.LValue) {

		// Get string index
		index := i.String()

		// Switch value type
		switch lv := v.(type) {
		case *lua.LTable:

			// Convert table to map
			n := TableToMap(v.(*lua.LTable))
			m[index] = n

		case lua.LNumber:

			// Convert to number to float64
			m[index] = float64(lv)

		case lua.LBool:

			// Convert value to boolean
			m[index] = bool(lv)

		case lua.LString:

			// Convert value to string
			m[index] = string(lv)
		}
	})
	return m
}

// URLValuesToTable converts a map[string][]string to a LUA table
func URLValuesToTable(m url.Values) *lua.LTable {
	// Data holder
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

// StructToTable converts a go struct pointer to a lua table
func StructToTable(s interface{}) *lua.LTable {
	// Data holder
	t := &lua.LTable{}

	// Get interface element
	elem := reflect.ValueOf(s).Elem()

	// Loop struct fields
	for i := 0; i < elem.NumField(); i++ {

		// Get current field
		field := elem.Field(i)

		// Get field name
		fieldName := elem.Type().Field(i).Name

		// Switch field type
		switch field.Interface().(type) {
		case string:

			// Set value as string
			t.RawSetString(fieldName, lua.LString(field.Interface().(string)))

		case int64:

			// Set value as number
			t.RawSetString(fieldName, lua.LNumber(field.Interface().(int64)))

		case uint32:

			// Set value as number
			t.RawSetString(fieldName, lua.LNumber(field.Interface().(uint32)))

		case int:

			// Set value as number
			t.RawSetString(fieldName, lua.LNumber(field.Interface().(int)))

		case float64:

			// Set value as number
			t.RawSetString(fieldName, lua.LNumber(field.Interface().(float64)))

		case bool:

			// Set value as bool
			t.RawSetString(fieldName, lua.LBool(field.Interface().(bool)))

		case time.Time:

			// Get time value
			timeStamp := field.Interface().(time.Time).Unix()

			// Set value as number
			t.RawSetString(fieldName, lua.LNumber(timeStamp))

		}
	}

	return t
}
