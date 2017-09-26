package lua

import (
	"os"

	"github.com/raggaer/castro/app/util"
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

// SetConfigCustomValue sets the a config custom value
func SetConfigCustomValue(L *lua.LState) int {
	// Get config key
	key := L.ToString(2)

	// Get config value
	value := L.Get(3)

	switch lv := value.(type) {
	case *lua.LTable:

		// Convert table to go map
		v := TableToMap(value.(*lua.LTable))

		// Insert map into custom field
		util.Config.SetCustomValue(key, v)

	case lua.LBool:

		// Insert value as bool
		util.Config.SetCustomValue(key, bool(lv))

	case lua.LNumber:

		// Insert value as float64
		util.Config.SetCustomValue(key, float64(lv))

	case lua.LString:

		// Insert value as string
		util.Config.SetCustomValue(key, string(lv))
	}

	// Open configuration file
	configFile, err := os.OpenFile("config.toml", os.O_RDWR|os.O_TRUNC, 0660)

	if err != nil {
		L.RaiseError("Cannot open configuration file: %v", err)
		return 0
	}

	// Close configuration file
	defer configFile.Close()

	if err := util.EncodeConfig(configFile, util.Config.Configuration); err != nil {
		L.RaiseError("Cannot encode configuration file: %v", err)
	}

	return 0
}
