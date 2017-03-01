package lua

import (
	"database/sql"
	"fmt"
	"github.com/yuin/gopher-lua"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type castroInterfaceField struct {
	Name string
	Type reflect.Kind
}

// MapToTable converts a Go map to a lua table
func MapToTable(m map[string]interface{}) *lua.LTable {
	// Main table pointer
	resultTable := &lua.LTable{}

	// Loop map
	for key, element := range m {

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

			// Get current table
			tbl := v.(*lua.LTable)

			// Check for lua table
			if tbl.MaxN() == 0 {

				// Convert table to map
				n := TableToMap(tbl)
				m[index] = n

			} else {

				// Data holder
				ret := make([]interface{}, 0, tbl.MaxN())

				// Loop table
				tbl.ForEach(func(i lua.LValue, v lua.LValue) {

					// Append to array
					ret = append(ret, ValueToGo(v))
				})

				// Set array
				m[index] = ret
			}

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

// ValueToGo converts a lua value to a go type
func ValueToGo(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[string]interface{})
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(ValueToGo(key))
				ret[keystr] = ValueToGo(value)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, ValueToGo(v.RawGetInt(i)))
			}
			return ret
		}
	default:
		return v
	}
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

		// Interface layer
		inter := field.Interface()

		// Switch field type
		switch inter.(type) {
		case string:

			// Set value as string
			t.RawSetString(fieldName, lua.LString(inter.(string)))

		case int64:

			// Set value as number
			t.RawSetString(fieldName, lua.LNumber(inter.(int64)))

		case uint32:

			// Set value as number
			t.RawSetString(fieldName, lua.LNumber(inter.(uint32)))

		case int:

			// Set value as number
			t.RawSetString(fieldName, lua.LNumber(inter.(int)))

		case float64:

			// Set value as number
			t.RawSetString(fieldName, lua.LNumber(inter.(float64)))

		case bool:

			// Set value as bool
			t.RawSetString(fieldName, lua.LBool(inter.(bool)))

		case time.Time:

			// Get time value
			timeStamp := inter.(time.Time).Unix()

			// Set value as number
			t.RawSetString(fieldName, lua.LNumber(timeStamp))

		case sql.NullString:

			// Get null string
			str := inter.(sql.NullString)

			// Check if string is not NULL
			if str.Valid {

				// Set value as string
				t.RawSetString(fieldName, lua.LString(str.String))
			} else {

				// Set value as nil
				t.RawSetString(fieldName, lua.LNil)
			}
		}
	}

	return t
}
