package lua

import "github.com/yuin/gopher-lua"

// Config holds the current lua configuration file
var Config = &Configuration{}

// Configuration struct used for LUA configuration
// file
type Configuration struct {
	WorldType           string `lua:"worldType"`
	HotkeyAimbotEnabled bool   `lua:"hotkeyAimbotEnabled"`
	ProtectionLevel     int    `lua:"protectionLevel"`
	RedSkullKills       int    `lua:"killsToRedSkull"`
	BlackSkullKills     int    `lua:"killsToBlackSkull"`
	PzLocked            int    `lua:"pzLocked"`
}

// LoadConfig loads the lua configuration file using
// lua vm to get the global variables
func LoadConfig(path string, dest *Configuration) error {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile(path + "/config.lua"); err != nil {
		return err
	}
	return GetStructVariables(dest, L)
}
