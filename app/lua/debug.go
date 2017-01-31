package lua

import (
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
)

// SetDebugMetaTable sets the debug metatable of the given state
func SetDebugMetaTable(luaState *lua.LState) int {
	// Create and set the debug metatable
	debugMetaTable := luaState.NewTypeMetatable(DebugMetaTableName)
	luaState.SetGlobal(DebugMetaTableName, debugMetaTable)

	// Set all crypto metatable functions
	luaState.SetFuncs(debugMetaTable, debugMethods)
}

// DebugValue prints a the value type and all the contents of the value
func DebugValue(L *lua.LState) int {
	// Get value
	val := L.Get(2)

	// Log value type
	util.Logger.Infof(" >> DEBUG - Value type %v", val.Type().String())

	// Switch value type
	switch val.Type() {
	case lua.LTString:
		util.Logger.Infof(" >> DEBUG - Value %v", val.String())
	case lua.LTNumber:
		util.Logger.Infof(" >> DEBUG - Value %v", val.String())
	}

	return 0
}
