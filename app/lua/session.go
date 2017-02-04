package lua

import (
	"github.com/goincremental/negroni-sessions"
	"github.com/raggaer/castro/app/models"
	"github.com/yuin/gopher-lua"
	"strconv"
)

// SetSessionMetaTable sets the session metatable on the given
// lua state
func SetSessionMetaTable(luaState *lua.LState) {
	// Create and set session metatable
	jwtMetaTable := luaState.NewTypeMetatable(SessionMetaTable)
	luaState.SetGlobal(SessionMetaTable, jwtMetaTable)

	// Set all map metatable functions
	luaState.SetFuncs(jwtMetaTable, sessionMethods)
}

// SetSessionMetaTableUserData sets the session metatable user data
func SetSessionMetaTableUserData(luaState *lua.LState, sessionData sessions.Session) {
	// Get session metatable
	jwtMetaTable := luaState.GetTypeMetatable(SessionMetaTable)

	// Set session field
	sess := luaState.NewUserData()
	sess.Value = sessionData
	luaState.SetField(jwtMetaTable, SessionInstanceName, sess)
}

// getSessionData gets the user data struct from the
// session metatable and returns the session pointer
func getSessionData(L *lua.LState) sessions.Session {
	// Get metatable
	meta := L.GetTypeMetatable(SessionMetaTable)

	// Get user data field
	data := L.GetField(meta, SessionInstanceName).(*lua.LUserData)

	// Return session struct
	return data.Value.(sessions.Session)
}

// GetLoggedAccount gets the user account if any
func GetLoggedAccount(L *lua.LState) int {
	// Get session data from the user data field
	session := getSessionData(L)

	// Check if user is logged
	logged, ok := session.Get("logged").(bool)

	if !ok {

		// Return nil if user is not logged in
		L.Push(lua.LNil)
		return 1
	}

	if !logged {

		// Return nil if user is not logged in
		L.Push(lua.LNil)
		return 1
	}

	// Get logged account name
	accountName, ok := session.Get("logged-account").(string)

	if !ok {

		// Return nil if invalid account name
		L.Push(lua.LNil)
		return 1
	}

	// Get accounts from database
	account, castroAccount, err := models.GetAccountByName(accountName)

	if err != nil {

		L.RaiseError("Cannot get account by name: %v", err)
		return 0
	}

	// Convert tfs account to lua table
	t := StructToTable(&account)

	// Set castro account inside the table
	t.RawSetString("castro", StructToTable(&castroAccount))

	// Send table to stack
	L.Push(t)

	return 1
}

// DestroySession removes the session data from the database
func DestroySession(L *lua.LState) int {
	// Get session data from the user data field
	session := getSessionData(L)

	// Destroy user session
	session.Clear()

	return 0
}

// SetSessionData saves an item to the session map
func SetSessionData(L *lua.LState) int {
	// Get key
	key := L.Get(2)

	// Check for valid key type
	if key.Type() != lua.LTString {

		L.ArgError(1, "Invalid key format. Expected string")
		return 0
	}

	// Get session data from the user data field
	session := getSessionData(L)

	// Get value
	val := L.Get(3)

	// Transform value to Go type
	switch val.Type() {
	case lua.LTString:

		// Assign element as string
		session.Set(key.String(), val.String())

	case lua.LTNumber:

		// Convert element to flaot64
		num, err := strconv.ParseFloat(val.String(), 64)

		if err != nil {

			L.ArgError(2, "Invalid number format. Cannot convert to Go type int64")
			return 0
		}

		// Assign element as flaot64
		session.Set(key.String(), num)

	case lua.LTBool:

		// Convert element to boolean
		b, err := strconv.ParseBool(val.String())

		if err != nil {

			L.ArgError(2, "Invalid boolean format. Cannot convert to Go type bool")
			return 0
		}

		// Assign element as bool
		session.Set(key.String(), b)

	case lua.LTTable:

		// Convert table to map
		m := TableToMap(val.(*lua.LTable))

		// Assign element as map
		session.Set(key.String(), m)
	}

	return 0
}

// GetSessionData retrieves an element from the session map
func GetSessionData(L *lua.LState) int {
	// Get key
	key := L.Get(2)

	// Check for valid key type
	if key.Type() != lua.LTString {

		L.ArgError(1, "Invalid key format. Expected string")
		return 0
	}

	// Get session data from the user data field
	session := getSessionData(L)

	// Get element from session
	val := session.Get(key.String())

	// Push element depending on the Go type
	switch val.(type) {
	case int64:

		// Push element as number
		L.Push(lua.LNumber(val.(int64)))
	case string:

		// Push element as string
		L.Push(lua.LString(val.(string)))
	case bool:

		// Push element as boolean
		L.Push(lua.LBool(val.(bool)))
	case map[string]interface{}:

		// Convert map to lua table
		tble := MapToTable(val.(map[string]interface{}))

		// Push element as table
		L.Push(tble)
	default:
		L.RaiseError("Unexpected data format")
	}

	return 1
}

// IsLogged checks if the current user is logged in
func IsLogged(L *lua.LState) int {
	// Get session data from the user data field
	session := getSessionData(L)

	// Try to get logged field from data
	b, ok := session.Get("logged").(bool)

	// If element does not exist push false
	if !ok {
		L.Push(lua.LBool(false))

		return 1
	}

	// Check the logged field
	L.Push(
		lua.LBool(b),
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
	v, ok := session.Get(key.String()).(string)

	if !ok {
		L.Push(lua.LString(""))
		return 1
	}

	// Delete value from flash map
	session.Delete(key.String())

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
	session.Set(key.String(), content.String())

	return 0
}
