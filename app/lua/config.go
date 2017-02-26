package lua

import (
	"github.com/yuin/gopher-lua"
)

// Config holds the current lua configuration file state
var Config = lua.NewState()

// SetConfigMetaTable sets the config metatable of the given state
func SetConfigMetaTable(luaState *lua.LState) {
	// Create and set Config metatable
	configMetaTable := luaState.NewTypeMetatable(ConfigMetaTableName)
	luaState.SetGlobal(ConfigMetaTableName, configMetaTable)

	// Set all Config metatable functions
	luaState.SetFuncs(configMetaTable, configMethods)
}

// LoadConfig loads the lua configuration file using lua vm to get the global variables
func LoadConfig(path string) error {
	// Execute config lua file
	return Config.DoFile(path)
}

// GetConfigLuaValue gets a value from the config struct using reflect
func GetConfigLuaValue(L *lua.LState) int {
	// Get value name
	name := L.ToString(2)

	// Get global from main state
	L.Push(Config.GetGlobal(name))

	return 1
}
