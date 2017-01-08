package lua

import (
	"github.com/yuin/gopher-lua"
	"github.com/raggaer/castro/app/util"
)

// GetVocationByName gets a vocation by the given name
func GetVocationByName(L *lua.LState) int {
	// Get name
	name := L.Get(2)

	// Check for valid name type
	if name.Type() != lua.LTString {

		L.ArgError(1, "Invalid name format. Expected string")
		return 0
	}

	// Get vocation
	for _, voc := range util.ServerVocationList.List.Vocations {

		// If it is the vocation we are looking for
		if voc.Name == name.String() {

			// Push vocation as lua table
			L.Push(VocationToTable(voc))

			return 1
		}
	}

	return 0
}

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