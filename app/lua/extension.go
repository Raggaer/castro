package lua

import (
	"github.com/raggaer/castro/app/util"
	lua "github.com/yuin/gopher-lua"
)

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
	if err := CompiledPageList.CompileExtensions("pages"); err != nil {
		L.RaiseError("Cannot reload extension page list: %v", err)
	}

	// Reload extension widgets
	if err := WidgetList.LoadExtensions(); err != nil {
		L.RaiseError("Cannot reload extension widget list: %v", err)
	}

	// Reload extension widget list
	if err := util.Widgets.LoadExtensions(); err != nil {
		L.RaiseError("Cannot reload extension widget list: %v", err)
	}

	return 0
}
