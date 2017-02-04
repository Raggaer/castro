package lua

import (
	"fmt"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"time"
)

// SetMapMetaTable sets a map metatable for the given state
func SetMapMetaTable(luaState *lua.LState) {
	// Create and set map metatable
	mapMetaTable := luaState.NewTypeMetatable(MapMetaTableName)
	luaState.SetGlobal(MapMetaTableName, mapMetaTable)

	// Set all map metatable functions
	luaState.SetFuncs(mapMetaTable, mapMethods)
}

// HouseList returns the server house list
// as a lua table
func HouseList(L *lua.LState) int {
	// Check if user wants specific town
	town := uint32(L.ToInt(2))

	// Check if list is on the cache
	list, found := util.Cache.Get(
		fmt.Sprintf("house_list_%v", town),
	)

	// If list is on the cache return cache result
	if found {

		L.Push(list.(*lua.LTable))

		return 1
	}

	// Result table
	tbl := &lua.LTable{}

	// Loop house list
	for _, house := range util.ServerHouseList.List.Houses {

		// Check if user wants specific town
		if town == 0 {

			// Convert house to table
			h := StructToTable(house)

			// Append to final table
			tbl.Append(h)

		} else if town == house.TownID {

			// Convert house to table
			h := StructToTable(house)

			// Append to final table
			tbl.Append(h)
		}
	}

	// Save list to cache
	util.Cache.Add(
		fmt.Sprintf("house_list_%v", town),
		tbl,
		time.Minute*3,
	)

	// Push table to stack
	L.Push(tbl)

	return 1
}

// TownList returns the server town list
// as a lua table
func TownList(L *lua.LState) int {
	// Convert town list to table and push
	// to stack
	L.Push(TownListToTable(util.OTBMap.Towns))

	return 1
}

// GetTownByName grabs a town by the given name
func GetTownByName(L *lua.LState) int {
	// Get town name
	name := L.Get(2)

	// Check for valid ID type
	if name.Type() != lua.LTString {

		L.ArgError(1, "Invalid name format. Expected string")
		return 0
	}

	// Get town
	for _, town := range util.OTBMap.Towns {

		// If its the town we are looking for
		if town.Name == name.String() {

			// Convert town to lua table and push
			L.Push(StructToTable(&town))

			return 1
		}
	}

	return 0
}

// GetTownByID grabs a town by the given ID
func GetTownByID(L *lua.LState) int {
	// Get town ID
	id := L.Get(2)

	// Check for valid ID type
	if id.Type() != lua.LTNumber {

		L.ArgError(1, "Invalid ID format. Expected number")
		return 0
	}

	// Convert town id to uint32
	townid := uint32(L.ToInt(2))

	// Get town
	for _, town := range util.OTBMap.Towns {

		// If its the town we are looking for
		if town.ID == townid {

			// Convert town to lua table and push
			L.Push(StructToTable(&town))

			return 1
		}
	}

	return 0
}
