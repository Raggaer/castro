package lua

import (
	"reflect"
	"strconv"

	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	"github.com/raggaer/otmap"
	"github.com/yuin/gopher-lua"
	"net/url"
	"strings"
	"time"
)

// GetStructVariables loads all the global variables
// from a lua file into a struct using reflect
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
		}
	}

	return resultTable
}

// TableToMap converts a LUA table to a Go map[string]interface{}
func TableToMap(table *lua.LTable) map[string]interface{} {
	// Check for valid table
	if table == nil {
		return nil
	}

	// Data holder
	m := make(map[string]interface{})

	// Loop lua table
	table.ForEach(func(i lua.LValue, v lua.LValue) {

		// Switch value type
		switch v.Type() {
		case lua.LTTable:

			// Convert table to map
			n := TableToMap(v.(*lua.LTable))
			m[i.String()] = n
		case lua.LTNumber:

			// Convert to number to float64
			num, err := strconv.ParseFloat(v.String(), 64)
			if err != nil {
				m[i.String()] = err.Error()
			} else {
				m[i.String()] = num
			}
		case lua.LTBool:

			// Convert value to boolean
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
				t.RawSetString(names[x], lua.LString(string(v.([]uint8))))
			case time.Time:
				t.RawSetString(names[x], lua.LNumber(v.(time.Time).Unix()))
			case int64:
				t.RawSetString(names[x], lua.LNumber(v.(int64)))
			case string:
				t.RawSetString(names[x], lua.LString(v.(string)))
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

// HouseListToTable converts a slice of houses to a lua table
func HouseListToTable(list []*util.House, townid uint32) *lua.LTable {
	t := &lua.LTable{}

	// Loop house list
	for _, house := range list {

		if townid != 0 && townid != house.TownID {

			continue
		}

		// Create a table for each house
		houseTable := &lua.LTable{}

		// Set table fields
		houseTable.RawSetString("name", lua.LString(house.Name))
		houseTable.RawSetString("size", lua.LNumber(house.Size))
		houseTable.RawSetString("townid", lua.LNumber(house.TownID))
		houseTable.RawSetString("id", lua.LNumber(house.ID))

		// Append house table to main table
		t.Append(houseTable)
	}

	return t
}

// TownListToTable converts a slice of towns to a lua table
func TownListToTable(list []otmap.Town) *lua.LTable {
	t := &lua.LTable{}

	// Loop town list
	for _, town := range list {

		// Append town table to main table
		t.Append(TownToTable(town))
	}

	return t
}

// TownToTable converts the given town to a lua table
func TownToTable(town otmap.Town) *lua.LTable {
	// Create a table for the town
	townTable := &lua.LTable{}

	// Set table fields
	townTable.RawSetString("name", lua.LString(town.Name))
	townTable.RawSetString("id", lua.LNumber(town.ID))

	return townTable
}

// VocationListToTable converts the server vocation list to a lua table
// executing the given condition
func VocationListToTable(list []*util.Vocation, cond func(*util.Vocation) bool) *lua.LTable {
	t := &lua.LTable{}

	// Loop vocation list
	for _, vocation := range list {

		// Check condition
		if !cond(vocation) {
			continue
		}

		// Append new table
		t.Append(VocationToTable(vocation))
	}

	return t
}

// VocationToTable converts a vocation to a lua table
func VocationToTable(voc *util.Vocation) *lua.LTable {
	t := &lua.LTable{}

	// Set table fields
	t.RawSetString("id", lua.LNumber(voc.ID))
	t.RawSetString("name", lua.LString(voc.Name))
	t.RawSetString("description", lua.LString(voc.Description))

	return t
}

// AccountToTable converts a tfs account to a lua table
func AccountToTable(account models.Account) *lua.LTable {
	t := &lua.LTable{}

	// Set account fields
	t.RawSetString("id", lua.LNumber(account.ID))
	t.RawSetString("premdays", lua.LNumber(account.Premdays))
	t.RawSetString("lastday", lua.LNumber(account.Lastday))
	t.RawSetString("name", lua.LString(account.Name))
	t.RawSetString("email", lua.LString(account.Email))

	return t
}

// CastroAccountToTable converts a castro account to a lua table
func CastroAccountToTable(account models.CastroAccount) *lua.LTable {
	t := &lua.LTable{}

	// Set account fields
	t.RawSetString("id", lua.LNumber(account.ID))
	t.RawSetString("points", lua.LNumber(account.Points))
	t.RawSetString("name", lua.LString(account.Name))

	return t
}
