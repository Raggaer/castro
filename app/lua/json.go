package lua

import (
	"encoding/json"
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
	r := TableToMap(L.ToTable(2))

	// Marshal converted table
	buff, err := json.Marshal(r)

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

	// Data holder
	result := make(map[string]interface{})

	if err := json.Unmarshal(
		[]byte(src.String()),
		&result,
	); err != nil {
		L.RaiseError("Cannot unmarshal the given string: %v", err)
		return 0
	}

	// Push result as table
	L.Push(MapToTable(result))

	return 1
}
