package lua

import (
	"github.com/raggaer/castro/app/util"
	lua "github.com/yuin/gopher-lua"
)

// SetOutfitMetaTable sets the outfit metatable of the given state
func SetOutfitMetaTable(luaState *lua.LState) {
	// Create and set the outfit metatable
	outfitMetaTable := luaState.NewTypeMetatable(OutfitMetaTableName)
	luaState.SetGlobal(OutfitMetaTableName, outfitMetaTable)

	// Set all outfit metatable functions
	luaState.SetFuncs(outfitMetaTable, outfitMethods)
}

// GenerateOutfit generates an outfit
func GenerateOutfit(L *lua.LState) int {
	// Get outfit bytes
	outfitBytes, err := util.GenerateOutfitImage(
		L.ToInt(2),
		L.ToInt(3),
		L.ToInt(4),
		L.ToInt(5),
		L.ToInt(6),
		L.ToInt(7),
	)

	if err != nil {
		L.RaiseError("Cannot generate outfit: %v", err)
		return 0
	}

	// Push outfitBytes to stack
	L.Push(lua.LString(string(outfitBytes)))

	return 1
}
