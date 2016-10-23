package lua

import "github.com/yuin/gopher-lua"

// Config holds the current lua configuration file
var Config = &Configuration{}

// Configuration struct used for LUA configuration
// file
type Configuration struct {
	WorldType           string
	HotkeyAimbotEnabled bool
}

// LoadConfig loads the lua configuration file using
// lua vm to get the global variables
func LoadConfig(path string, dest *Configuration) error {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile(path + "/config.lua"); err != nil {
		return err
	}
	dest.WorldType = lua.LVAsString(L.GetGlobal("worldType"))
	dest.HotkeyAimbotEnabled = lua.LVAsBool(L.GetGlobal("hotkeyAimbotEnabled"))
	return nil
}
