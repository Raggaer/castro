package lua

import (
	"github.com/yuin/gopher-lua"
	"github.com/raggaer/castro/app/util"
	"strconv"
)

// getSessionData gets the user data struct from the
// session metatable and returns the session pointer
func getSessionData(L *lua.LState) *util.Session {
	// Get metatable
	meta := L.GetTypeMetatable(JWTMetaTable)

	// Get user data field
	data := L.GetField(meta, JWTTokenName).(*lua.LUserData)

	// Return session struct
	return data.Value.(*util.Session)
}

// DestroySession removes the session data from the database
func DestroySession(L *lua.LState) int {

	// Get session data from the user data field
	session := getSessionData(L)

	// Destroy user session
	if err := session.Destroy(); err != nil {

		L.RaiseError("Cannot destroy user session: %v", err)
	}

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
		session.Data[key.String()] = val.String()
	case lua.LTNumber:

		// Convert element to int64
		num, err := strconv.ParseInt(val.String(), 10, 64)

		if err != nil {

			L.ArgError(2, "Invalid number format. Cannot convert to Go type int64")
			return 0
		}

		// Assign element as int64
		session.Data[key.String()] = num
	case lua.LTBool:

		// Convert element to boolean
		b, err := strconv.ParseBool(val.String())

		if err != nil {

			L.ArgError(2, "Invalid boolean format. Cannot convert to Go type bool")
			return 0
		}

		// Assign element as bool
		session.Data[key.String()] = b
	case lua.LTTable:

		// Convert table to map
		m := TableToMap(val.(*lua.LTable))

		// Assign element as map
		session.Data[key.String()] = m
	}

	// Save user session
	if err := session.Save(); err != nil {

		L.RaiseError("Cannot save user session: %v", err)
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
	val, ok := session.Data[key.String()]

	// Check if element exists
	if !ok {

		L.Push(lua.LNil)
		return 1
	}

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
	b, ok := session.Data["logged"].(bool)

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
	v, ok := session.Flash[key.String()]

	// Check if flash value exists
	if !ok {

		L.Push(lua.LString(""))
		return 0
	}

	// Delete value from flash map
	delete(session.Flash, key.String())

	// Update session
	if err := session.Save(); err != nil {

		L.RaiseError("Cannot save user session: %v", err)
	}

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