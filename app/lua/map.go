package lua

import (
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
)

// HouseList returns the server house list
// as a lua table
func HouseList(L *lua.LState) int {
	// Check if user wants specific town
	town := L.Get(2)

	// Result table
	var tbl *lua.LTable

	if town.Type() != lua.LTNumber {

		// Convert house list to table
		tbl = HouseListToTable(util.ServerHouseList.List.Houses, 0)
	} else {

		// Convert house list to table
		tbl = HouseListToTable(util.ServerHouseList.List.Houses, uint32(L.ToInt(2)))
	}

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
			L.Push(TownToTable(town))

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
			L.Push(TownToTable(town))

			return 1
		}
	}

	return 0
}
