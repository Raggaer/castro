package lua

import (
	"fmt"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
)

// SetReflectMetaTable sets the reflect metatable of the given state
func SetReflectMetaTable(luaState *lua.LState) {
	// Create and set the reflect metatable
	reflectMetaTable := luaState.NewTypeMetatable(ReflectMetaTableName)
	luaState.SetGlobal(ReflectMetaTableName, reflectMetaTable)

	// Set all reflect metatable functions
	luaState.SetFuncs(reflectMetaTable, reflectMethods)
}

// GetGlobal retrieves a global lua value from other script
func GetGlobal(L *lua.LState) int {
	// Get script location
	path := L.ToString(2)

	// Get value name
	val := L.ToString(3)

	// Check if value is on cache
	source, found := util.Cache.Get(
		fmt.Sprintf("reflect_global_%v", val),
	)

	if found {

		// Push source
		L.Push(source.(lua.LValue))

		return 1
	}

	// Get a new lua state
	state := Pool.Get()

	// Return state
	defer Pool.Put(state)

	// Execute the script
	if err := state.DoFile(path); err != nil {
		L.RaiseError("Cannot execute the given script: %v", err)
		return 0
	}

	// Get value from state
	v := state.GetGlobal(val)

	// Save global to cache
	util.Cache.Add(
		fmt.Sprintf("reflect_global_%v", val),
		v,
		util.Config.Configuration.Cache.Default,
	)

	// Push value
	L.Push(v)

	return 1
}
