package lua

import (
	"github.com/yuin/gopher-lua"
	"github.com/raggaer/castro/app/util"
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
