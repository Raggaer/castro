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