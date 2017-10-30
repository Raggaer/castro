package lua

import lua "github.com/yuin/gopher-lua"

// SetExtensionMetaTable sets the extension metatable for the given state
func SetExtensionMetaTable(luaState *lua.LState) {
	// Create and set the log metatable
	extMetaTable := luaState.NewTypeMetatable(ExtensionMetaTableName)
	luaState.SetGlobal(ExtensionMetaTableName, extMetaTable)

	// Set all mail metatable functions
	luaState.SetFuncs(extMetaTable, extensionMethods)
}

// ReloadExtensions reloads all extensions
func ReloadExtensions(L *lua.LState) int {
	// Reload extension pages
	if err := PageList.LoadExtensions(); err != nil {
		L.RaiseError("Cannot reload extension list: %v", err)
	}

	return 0
}
