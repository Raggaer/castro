package lua

import (
	"github.com/yuin/gopher-lua"
	"github.com/raggaer/castro/app/util"
)

// VocationList returns the server vocations xml file
func VocationList(L *lua.LState) int {
	// Check if user wants base vocations
	base := L.Get(2)

	// Check for valid base type
	if base.Type() != lua.LTBool {

		// Get vocation list table
		t := VocationListToTable(util.ServerVocationList.List.Vocations, func(v *util.Vocation) bool {
			return true
		})

		// Push table to stack
		L.Push(t)

		return 1
	}

	// Convert case to boolean
	b := L.ToBool(2)

	// Stupid case
	if !b {

		// Get vocation list table
		t := VocationListToTable(util.ServerVocationList.List.Vocations, func(v *util.Vocation) bool {
			return true
		})

		// Push table to stack
		L.Push(t)

		return 1
	}

	// Get base vocation list
	t := VocationListToTable(util.ServerVocationList.List.Vocations, func(v *util.Vocation) bool {
		return v.FromVoc == v.ID
	})

	// Push table to stack
	L.Push(t)

	return 1
}