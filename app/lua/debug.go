package lua

import (
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
)

// SetDebugMetaTable sets the debug metatable of the given state
func SetDebugMetaTable(luaState *lua.LState) {
	// Create and set the debug metatable
	debugMetaTable := luaState.NewTypeMetatable(DebugMetaTableName)
	luaState.SetGlobal(DebugMetaTableName, debugMetaTable)

	// Set all debug metatable functions
	luaState.SetFuncs(debugMetaTable, debugMethods)
}

// DebugValue prints the given values
func DebugValue(L *lua.LState) int {
	// Get value
	val := L.Get(2)

	// Set increment value
	i := 2

	// Loop all values
	for val.Type() != lua.LTNil {

		// Debug statement
		util.Logger.Logger.Debugf("type: %v, value: %v", val.Type().String(), val.String())

		// Increment value index
		i++

		// Get new value
		val = L.Get(i)
	}

	return 0
}
