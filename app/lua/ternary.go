package lua

import (
	lua "github.com/yuin/gopher-lua"
)

// Ternary creates a ternary operator short-hand
func Ternary(L *lua.LState) int {
	// Get ternary condition
	condition := L.ToBool(1)

	// Return first argument if valid condition
	if condition {
		L.Push(L.Get(2))
		return 1
	}

	// Return wrong condition argument
	L.Push(L.Get(3))
	return 1
}
