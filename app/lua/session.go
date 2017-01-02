package lua

import (
	"github.com/yuin/gopher-lua"
	"github.com/raggaer/castro/app/util"
)

func getSessionData(L *lua.LState) *util.Session {
	// Get metatable
	meta := L.GetTypeMetatable(JWTMetaTable)

	// Get user data field
	data := L.GetField(meta, JWTTokenName).(*lua.LUserData)

	// Return session struct
	return data.Value.(*util.Session)
}

// IsLogged checks if the current user is logged in
func IsLogged(L *lua.LState) int {
	// Get session data from the user data field
	session := getSessionData(L)

	// Check the logged field
	L.Push(
		lua.LBool(session.Logged),
	)

	return 1
}

// GetFlash gets a flash value from the user session
func GetFlash(L *lua.LState) int {
	// Get session data from the user data field
	session := getSessionData(L)

	// Get flash key
	key := L.Get(2)

	// Check for valid key
	if key.Type() != lua.LTString {

		L.ArgError(1, "Invalid flash key. Expected string")
		return 0
	}

	// Get value from the flash map
	v, ok := session.Flash[key.String()]

	// Check if flash value exists
	if !ok {

		L.Push(lua.LString(""))
		return 0
	}

	// Delete value from flash map
	delete(session.Flash, key.String())

	// Push value to stack
	L.Push(lua.LString(v))

	return 1
}

// SetFlash sets a flash value to the user session
func SetFlash(L *lua.LState) int {
	// Get session data from the user data field
	session := getSessionData(L)

	// Get flash key
	key := L.Get(2)

	// Check for valid key
	if key.Type() != lua.LTString {

		L.ArgError(1, "Invalid flash key. Expected string")
		return 0
	}

	// Get flash data
	content := L.Get(3)

	// Check for valid content
	if content.Type() != lua.LTString {

		L.ArgError(1, "Invalid flash content. Expected string")
		return 0
	}

	// Set flash value
	session.Flash[key.String()] = content.String()

	// Update session
	if err := session.Save(); err != nil {

		L.RaiseError("Cannot save user session: %v", err)
	}

	return 0
}