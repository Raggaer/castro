package lua

import (
	"encoding/base64"
	"fmt"

	"github.com/yuin/gopher-lua"
)

// SetBase64MetaTable sets the base64 metatable for the given state
func SetBase64MetaTable(luaState *lua.LState) {
	// Create and set the log metatable
	base64MetaTable := luaState.NewTypeMetatable(Base64MetaTableName)
	luaState.SetGlobal(Base64MetaTableName, base64MetaTable)

	// Set all base64 metatable functions
	luaState.SetFuncs(base64MetaTable, base64Methods)
}

// Base64Encode encodes the given string to base64
func Base64Encode(L *lua.LState) int {
	// Get string to be encoded
	str := L.Get(2)

	// Check for valid string type
	if str.Type() != lua.LTString {

		L.ArgError(1, "Invalid string format. Expected string")
		return 0
	}

	// Encode the string using base64
	encodedString := base64.URLEncoding.EncodeToString([]byte(str.String()))

	// Convert byte array to string and push to stack
	L.Push(lua.LString(encodedString))

	return 1
}

// Base64Decode decodes the given string from base64
func Base64Decode(L *lua.LState) int {
	// Get string to be decoded
	str := L.Get(2)

	// Check for valid string type
	if str.Type() != lua.LTString {

		L.ArgError(1, "Invalid string format. Expected string")
		return 0
	}

	// Decode base64 to string
	decodedString, err := base64.URLEncoding.DecodeString((str.String()))
	if err != nil {
		L.RaiseError(fmt.Sprintf("Failed to encode string: %v", err))
		return 0
	}

	// Convert byte array to string and push to stack
	L.Push(lua.LString(decodedString))

	return 1
}
