package lua

import glua "github.com/yuin/gopher-lua"

// NewState creates and returns a new lua.LState
// pointer with some options
func NewState() *glua.LState {
	luaState := glua.NewState(
		glua.Options{
			IncludeGoStackTrace: true,
		},
	)
	return luaState
}
