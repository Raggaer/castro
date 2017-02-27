package lua

import "github.com/yuin/gopher-lua"

// SetPayPalMetaTable sets the paypal metatable of the given state
func SetPayPalMetaTable(luaState *lua.LState) {
	// Create and set the paypal metatable
	paypalMetaTable := luaState.NewTypeMetatable(PayPalMetaTableName)
	luaState.SetGlobal(PayPalMetaTableName, paypalMetaTable)

	// Set all mail metatable functions
	luaState.SetFuncs(paypalMetaTable, paypalMethods)
}
