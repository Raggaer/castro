package lua

import (
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"net/url"
)

// SetURLMetaTable sets the url metatable of the given state
func SetURLMetaTable(luaState *lua.LState) {
	// Create and set the url metatable
	urlMetaTable := luaState.NewTypeMetatable(URLMetaTableName)
	luaState.SetGlobal(URLMetaTableName, urlMetaTable)

	// Set all url metatable functions
	luaState.SetFuncs(urlMetaTable, urlMethods)
}

// DecodeURL decodes the given string uri
func DecodeURL(L *lua.LState) int {
	// Get uri
	uri := L.Get(2)

	// Check for valid uri type
	if uri.Type() != lua.LTString {
		L.ArgError(1, "Invalid uri type. Expected string")
		return 0
	}

	// Decode uri
	decoded, err := url.QueryUnescape(uri.String())

	if err != nil {
		util.Logger.Errorf("Cannot decode uri: %v", err)

		// Push nil
		L.Push(lua.LNil)

		return 1
	}

	// Push decoded string
	L.Push(lua.LString(decoded))

	return 1
}
