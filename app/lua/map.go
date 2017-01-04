package lua

import (
	"github.com/yuin/gopher-lua"
	"github.com/raggaer/castro/app/util"
)

// HouseList returns the server house list
// as a lua table
func HouseList(L *lua.LState) int {
	// Convert house list to table
	tbl := HouseListToTable(util.ServerHouseList.List.Houses)

	// Push table to stack
	L.Push(tbl)

	return 1
}
