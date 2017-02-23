package lua

import (
	"github.com/clbanning/mxj"
	"github.com/yuin/gopher-lua"
)

// SetJSONMetaTable sets the json metatable of the given state
func SetJSONMetaTable(luaState *lua.LState) {
	// Create and set the json metatable
	jsonMetaTable := luaState.NewTypeMetatable(JSONMetaTableName)
	luaState.SetGlobal(JSONMetaTableName, jsonMetaTable)

	// Set all json metatable functions
	luaState.SetFuncs(jsonMetaTable, jsonMethods)
}

// MarshalJSON marshals the given lua table
func MarshalJSON(L *lua.LState) int {
	// Get table
	tbl := L.Get(2)

	// Check for valid table type
	if tbl.Type() != lua.LTTable {
		L.ArgError(1, "Invalid marshal object. Expected table")
		return 0
	}

	// Convert table to map
	r := mxj.Map(TableToMap(L.ToTable(2)))

	// Marshal converted table
	buff, err := r.Json()

	if err != nil {
		L.RaiseError("Cannot marshal the given table: %v", err)
		return 0
	}

	// Push result as string
	L.Push(lua.LString(string(buff)))

	return 1
}

// UnmarshalJSON unmarshals the given string to a lua table
func UnmarshalJSON(L *lua.LState) int {
	// Get string
	src := L.Get(2)

	// Check for valid string type
	if src.Type() != lua.LTString {
		L.ArgError(1, "Invalid unmarshal source. Expected string")
		return 0
	}

	// Unmarshal string
	result, err := mxj.NewMapJson([]byte(src.String()))

	if err != nil {
		L.RaiseError("Cannot unmarshal the given string: %v", err)
		return 0
	}

	// Push result as table
	L.Push(MapToTable(result))

	return 1
}
