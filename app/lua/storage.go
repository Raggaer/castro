package lua

import "github.com/yuin/gopher-lua"

// SetStorageMetaTable sets the storage metatable of the given state
func SetStorageMetaTable(luaState *lua.LState) {
	// Create and set the time metatable
	storageMetaTable := luaState.NewTypeMetatable(StorageMetaTableName)
	luaState.SetGlobal(StorageMetaTableName, storageMetaTable)

	// Set all time metatable functions
	luaState.SetFuncs(storageMetaTable, storageMethods)
}
