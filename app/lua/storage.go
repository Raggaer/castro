package lua

import (
	"github.com/raggaer/castro/app/database"
	"github.com/yuin/gopher-lua"
)

// StorageValue struct for storage values
type StorageValue struct {
	PlayerID int
	Key      int
	Value    int
}

// SetStorageMetaTable sets the storage metatable of the given state
func SetStorageMetaTable(luaState *lua.LState) {
	// Create and set the time metatable
	storageMetaTable := luaState.NewTypeMetatable(StorageMetaTableName)
	luaState.SetGlobal(StorageMetaTableName, storageMetaTable)

	// Set all time metatable functions
	luaState.SetFuncs(storageMetaTable, storageMethods)
}

// GetStorageValue gets a storage value from the given player
func GetStorageValue(L *lua.LState) int {
	// Get player id
	playerid := L.ToInt(2)

	// Get storage key
	key := L.ToInt(3)

	// Data holder
	s := StorageValue{}

	// Execute query
	if err := database.DB.Get(&s, "SELECT player_id, key, value FROM player_storage WHERE key = ? AND player_id = ?", key, playerid); err != nil {
		L.RaiseError("Cannot get storage value: %v", err)
		return 0
	}

	// Push storage value as table
	L.Push(StructToTable(&s))

	return 1
}

// SetStorageValue sets a storage value for the given player
func SetStorageValue(L *lua.LState) int {
	// Get player id
	playerid := L.ToInt(2)

	// Get storage key
	key := L.ToInt(3)

	// Get storage value
	v := L.ToInt(4)

	// Execute insert query
	if _, err := database.DB.Exec("INSERT INTO player_storage (player_id, key, value) VALUES (?, ?, ?)", playerid, key, v); err != nil {
		L.RaiseError("Cannot insert storage value: %v", err)
		return 0
	}

	return 0
}
