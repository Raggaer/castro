package lua

import (
	"github.com/yuin/gopher-lua"
	"os"
)

// SetEnvMetaTable sets the env metatable on the given lua state
func SetEnvMetaTable(luaState *lua.LState) {
	// Create and set env metatable
	envMetaTable := luaState.NewTypeMetatable(EnvMetaTableName)
	luaState.SetGlobal(EnvMetaTableName, envMetaTable)

	// Set all HTTP metatable functions
	luaState.SetFuncs(envMetaTable, envMethods)
}

// SetEnvVariable sets the given environment variable
func SetEnvVariable(L *lua.LState) int {
	// Get variable key
	key := L.Get(2)

	// Check valid key
	if key.Type() != lua.LTString {
		L.ArgError(1, "Invalid key type. Expected string")
		return 0
	}

	// Get value
	val := L.Get(3)

	// Check valid value
	if val.Type() != lua.LTString {
		L.ArgError(2, "Invalid content value. Expected string")
		return 0
	}

	// Set variable
	if err := os.Setenv(key.String(), val.String()); err != nil {
		L.RaiseError("Cannot set env variable: %v", err)
	}

	return 0
}

// GetEnvVariable gets the given environment variable
func GetEnvVariable(L *lua.LState) int {
	// Get variable key
	key := L.Get(2)

	// Check valid key
	if key.Type() != lua.LTString {
		L.ArgError(1, "Invalid key type. Expected string")
		return 0
	}

	// Get variable
	v := os.Getenv(key.String())

	// Push nil if variable is not set
	if v == "" {
		L.Push(lua.LNil)

		return 1
	}

	// Push variable
	L.Push(lua.LString(v))

	return 1
}
