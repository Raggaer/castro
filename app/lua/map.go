package lua

import (
	"fmt"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
)

// SetMapMetaTable sets a map metatable for the given state
func SetMapMetaTable(luaState *lua.LState) {
	// Create and set map metatable
	mapMetaTable := luaState.NewTypeMetatable(MapMetaTableName)
	luaState.SetGlobal(MapMetaTableName, mapMetaTable)

	// Set all map metatable functions
	luaState.SetFuncs(mapMetaTable, mapMethods)
}

// HouseList returns the server house list as a lua table
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
		util.Config.Configuration.Cache.Default.Duration,
	)

	// Push table to stack
	L.Push(tbl)

	return 1
}

// TownList returns the server town list as a lua table
func TownList(L *lua.LState) int {

	// Check if the list is on the cache
	list, found := util.Cache.Get("town_list")

	if found {

		// Push town list
		L.Push(list.(*lua.LTable))

		return 1
	}

	// Data holder
	result := &lua.LTable{}

	// Loop town list
	for _, town := range util.OTBMap.Map.Towns {

		// Convert town to table
		tbl := StructToTable(&town)

		// Append to main table
		result.Append(tbl)
	}

	// Save list to cache
	util.Cache.Add("town_list", result, util.Config.Configuration.Cache.Default.Duration)

	// Push result
	L.Push(result)

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

	// Check if town is on cache
	t, found := util.Cache.Get(
		fmt.Sprintf("town_%v", name.String()),
	)

	if found {

		// Push town table
		L.Push(t.(*lua.LTable))

		return 1
	}

	// Get town
	for _, town := range util.OTBMap.Map.Towns {

		// If its the town we are looking for
		if town.Name == name.String() {

			twn := StructToTable(&town)

			// Save town to cache
			util.Cache.Add(
				fmt.Sprintf("town_%v", name.String()),
				twn,
				util.Config.Configuration.Cache.Default.Duration,
			)

			// Convert town to lua table and push
			L.Push(twn)

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
	for _, town := range util.OTBMap.Map.Towns {

		// If its the town we are looking for
		if town.ID == townid {

			// Check if town is on cache
			t, found := util.Cache.Get(
				fmt.Sprintf("town_%v", town.Name),
			)

			if found {

				// Push town table
				L.Push(t.(*lua.LTable))

				return 1
			}

			// Convert town to table
			twn := StructToTable(&town)

			// Save town to cache
			util.Cache.Add(
				fmt.Sprintf("town_%v", town.Name),
				twn,
				util.Config.Configuration.Cache.Default.Duration,
			)

			// Convert town to lua table and push
			L.Push(twn)

			return 1
		}
	}

	return 0
}
