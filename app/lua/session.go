package lua

import (
	"net/http"

	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
)

// SetSessionMetaTable sets the session metatable on the given lua state
func SetSessionMetaTable(luaState *lua.LState) {
	// Create and set session metatable
	jwtMetaTable := luaState.NewTypeMetatable(SessionMetaTable)
	luaState.SetGlobal(SessionMetaTable, jwtMetaTable)

	// Set all map metatable functions
	luaState.SetFuncs(jwtMetaTable, sessionMethods)
}

// SetSessionMetaTableUserData sets the session metatable user data
func SetSessionMetaTableUserData(luaState *lua.LState, sessionData map[string]interface{}) {
	// Get session metatable
	jwtMetaTable := luaState.GetTypeMetatable(SessionMetaTable)

	// Set session field
	sess := luaState.NewUserData()
	sess.Value = sessionData
	luaState.SetField(jwtMetaTable, SessionInstanceName, sess)
}

// getSessionData gets the user data struct from the session metatable and returns the session pointer
func getSessionData(L *lua.LState) map[string]interface{} {
	// Get metatable
	meta := L.GetTypeMetatable(SessionMetaTable)

	// Get user data field
	data := L.GetField(meta, SessionInstanceName).(*lua.LUserData)

	// Return session struct
	return data.Value.(map[string]interface{})
}

// updateSessionData saves a new cookie with the encoded map
func updateSessionData(L *lua.LState) {
	// Get response writer from state
	_, w := getRequestAndResponseWriter(L)

	// Get session
	session := getSessionData(L)

	// Encode session map
	encoded, err := util.SessionStore.Encode(util.Config.Configuration.Cookies.Name, session)

	if err != nil {
		util.Logger.Logger.Fatalf("Cannot encode cookie value: %v", err)
	}

	// Create cookie
	c := util.SessionCookie(encoded)

	// Set cookie
	http.SetCookie(w, c)
}

// GetLoggedAccount gets the user account if any
func GetLoggedAccount(L *lua.LState) int {
	// Get session data from the user data field
	session := getSessionData(L)

	// Check if user is logged
	logged, ok := session["logged"].(bool)

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
	accountName, ok := session["loggedAccount"].(string)

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

// IsAdmin checks if the logged account is admin
func IsAdmin(L *lua.LState) int {
	// Get session data from the user data field
	session := getSessionData(L)

	// Try to get logged field from data
	b, ok := session["logged"].(bool)

	// If element does not exist push false
	if !ok {
		L.Push(lua.LBool(false))
		return 1
	}

	// Check the session value
	if !b {
		L.Push(lua.LBool(false))
		return 1
	}

	// Get logged account name
	accountName, ok := session["loggedAccount"].(string)

	if !ok {

		// Return nil if invalid account name
		L.Push(lua.LNil)
		return 1
	}

	// Get accounts from database
	_, castroAccount, err := models.GetAccountByName(accountName)

	if err != nil {
		L.RaiseError("Cannot get account by name: %v", err)
		return 0
	}

	// Push admin status
	L.Push(lua.LBool(castroAccount.Admin))

	return 1
}

// DestroySession removes the session data from the database
func DestroySession(L *lua.LState) int {
	// Get session data from the user data field
	session := getSessionData(L)

	// Loop map
	for key := range session {

		// Omit issuer element
		if key == "issuer" || key == "csrf-token" {
			continue
		}

		// Delete each element
		delete(session, key)
	}

	// Update session data
	updateSessionData(L)

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
	switch lv := val.(type) {
	case lua.LString:

		// Assign element as string
		session[key.String()] = string(lv)

	case lua.LNumber:

		// Assign element as float64
		session[key.String()] = float64(lv)

	case lua.LBool:

		// Assign element as bool
		session[key.String()] = bool(lv)

	case *lua.LTable:

		// Convert table to map
		m := TableToMap(val.(*lua.LTable))

		// Assign element as map
		session[key.String()] = m
	}

	// Update session data
	updateSessionData(L)

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
	val := session[key.String()]

	// Push element depending on the Go type
	switch val.(type) {
	case float64:

		// Push element as number
		L.Push(lua.LNumber(val.(float64)))
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
		L.Push(lua.LNil)
	}

	return 1
}

// IsLogged checks if the current user is logged in
func IsLogged(L *lua.LState) int {
	// Get session data from the user data field
	session := getSessionData(L)

	// Try to get logged field from data
	b, ok := session["logged"].(bool)

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

	// Get flash key
	key := L.Get(2)

	// Check for valid key
	if key.Type() != lua.LTString {

		L.ArgError(1, "Invalid flash key. Expected string")
		return 0
	}

	// Get session data from the user data field
	session := getSessionData(L)

	// Get element from session
	val := session[key.String()]

	// Push element depending on the Go type
	switch val.(type) {
	case float64:

		// Push element as number
		L.Push(lua.LNumber(val.(float64)))
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
		L.Push(lua.LNil)
		return 1
	}

	// Delete element from map
	delete(session, key.String())

	// Update session data
	updateSessionData(L)

	return 1
}

// SetFlash sets a flash value to the user session
func SetFlash(L *lua.LState) int {
	
	// Get flash key
	key := L.Get(2)

	// Check for valid key
	if key.Type() != lua.LTString {

		L.ArgError(1, "Invalid flash key. Expected string")
		return 0
	}

	// Get session data from the user data field
	session := getSessionData(L)

	// Get value
	val := L.Get(3)

	// Transform value to Go type
	switch lv := val.(type) {
	case lua.LString:

		// Assign element as string
		session[key.String()] = string(lv)

	case lua.LNumber:

		// Assign element as float64
		session[key.String()] = float64(lv)

	case lua.LBool:

		// Assign element as bool
		session[key.String()] = bool(lv)

	case *lua.LTable:

		// Convert table to map
		m := TableToMap(val.(*lua.LTable))

		// Assign element as map
		session[key.String()] = m
	}

	// Update session data
	updateSessionData(L)

	return 0
}
