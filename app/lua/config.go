package lua

import (
	"strconv"

	"github.com/yuin/gopher-lua"
)

// Config holds the current lua configuration file
var Config = &Configuration{}

// Configuration struct used for LUA configuration
// file
type Configuration struct {
	WorldType           string
	HotkeyAimbotEnabled bool
	ProtectionLevel     int
	RedSkullKills       int
	BlackSkullKills     int
	PzLocked            int
}

// LoadConfig loads the lua configuration file using
// lua vm to get the global variables
func LoadConfig(path string, dest *Configuration) error {
	L := lua.NewState()
	defer L.Close()
	err := L.DoFile(path + "/config.lua")
	if err != nil {
		return err
	}
	dest.WorldType = lua.LVAsString(L.GetGlobal("worldType"))
	dest.HotkeyAimbotEnabled = lua.LVAsBool(L.GetGlobal("hotkeyAimbotEnabled"))
	dest.ProtectionLevel, err = strconv.Atoi(lua.LVAsString(L.GetGlobal("protectionLevel")))
	if err != nil {
		return err
	}
	dest.RedSkullKills, err = strconv.Atoi(lua.LVAsString(L.GetGlobal("killsToRedSkull")))
	if err != nil {
		return err
	}
	dest.BlackSkullKills, err = strconv.Atoi(lua.LVAsString(L.GetGlobal("killsToBlackSkull")))
	if err != nil {
		return err
	}
	dest.PzLocked, err = strconv.Atoi(lua.LVAsString(L.GetGlobal("pzLocked")))
	if err != nil {
		return err
	}
	return nil
}
