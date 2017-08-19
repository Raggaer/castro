package lua

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"github.com/raggaer/castro/app/database"
	"github.com/yuin/gopher-lua"
)

// SetGlobalMetaTable sets the file metatable of the given state
func SetGlobalMetaTable(luaState *lua.LState) {
	// Create and set the file metatable
	globalMetaTable := luaState.NewTypeMetatable(GlobalMetaTableName)
	luaState.SetGlobal(GlobalMetaTableName, globalMetaTable)

	// Set all global metatable functions
	luaState.SetFuncs(globalMetaTable, globalMethods)
}

// SetGlobalLuaValue saves the given value into the global database table
func SetGlobalLuaValue(L *lua.LState) int {
	// Get value key
	key := L.ToString(2)

	// Get content table
	content := TableToMap(L.ToTable(3))

	// Create buffer
	buff := bytes.NewBuffer([]byte{})

	// Create gob encoder
	encoder := gob.NewEncoder(buff)

	// Encode table to byte array
	if err := encoder.Encode(content); err != nil {
		L.RaiseError("Cannot encode given lua table: %v", err)
		return 0
	}

	// Global value exists placeholder
	var exists = true

	// Check if key already exists
	if err := database.DB.Get(&exists, "SELECT id FROM castor_global WHERE `key` = ?", key); err != nil {

		// Global value does not exist
		if err == sql.ErrNoRows {

			// Insert value into database
			if _, err := executeQueryHelper(L, "INSERT INTO castro_global (`key`, value) VALUES (?, ?)", key, buff.Bytes()); err != nil {
				L.RaiseError("Cannot save encoded lua table to database: %v", err)
			}

			return 0
		}

		L.RaiseError("Cannot check for global value: %v", err)
		return 0
	}

	// Update global value
	if _, err := executeQueryHelper(L, "UPDATE castro_global SET value = ? WHERE `key` = ?", buff.Bytes(), key); err != nil {
		L.RaiseError("Cannot update global value: %v", err)
	}

	return 0
}

// GetGlobalLuaValue gets the given value from the global database table
func GetGlobalLuaValue(L *lua.LState) int {
	// Get value key
	key := L.ToString(2)

	// Result placeholder
	b := []byte{}

	// Retrieve value from database
	if err := database.DB.Get(&b, "SELECT value FROM castro_global WHERE `key` = ?", key); err != nil {

		// Check if there is no match
		if err == sql.ErrNoRows {

			// Push nil value
			L.Push(lua.LNil)
			return 1
		}

		L.RaiseError("Cannot get encoded value from database: %v", err)
		return 0
	}

	// Create buffer
	buff := bytes.NewBuffer(b)

	// Create decoder
	decoder := gob.NewDecoder(buff)

	// Decoded value placeholder
	result := map[string]interface{}{}

	// Decode value
	if err := decoder.Decode(&result); err != nil {
		L.RaiseError("Cannot decode lua table: %v", err)
		return 0
	}

	// Push result as lua table
	L.Push(MapToTable(result))

	return 1
}
