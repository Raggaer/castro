package lua

import (
	"reflect"
	"strconv"

	"github.com/yuin/gopher-lua"
)

// GetStructVariables loads all the global variables
// from a lua file into a struct using reflect
func GetStructVariables(src interface{}, L *lua.LState) error {
	v := reflect.ValueOf(src).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldTag := v.Type().Field(i)
		if t, ok := fieldTag.Tag.Lookup("lua"); ok {
			if t == "" {
				continue
			}
			variable := L.GetGlobal(t)
			if variable.Type() == lua.LTNil {
				continue
			}
			switch variable.Type() {
			case lua.LTNumber:
				n, err := strconv.ParseInt(variable.String(), 10, 64)
				if err != nil {
					return err
				}
				field.SetInt(n)
			case lua.LTBool:
				field.SetBool(lua.LVAsBool(variable))
			case lua.LTString:
				field.SetString(variable.String())
			}

		}
	}
	return nil
}
