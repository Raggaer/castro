package lua

import (
	"github.com/yuin/gopher-lua"
	"os"
)

// SetFileMetaTable sets the file metatable of the given state
func SetFileMetaTable(luaState *lua.LState) {
	// Create and set the file metatable
	fileMetaTable := luaState.NewTypeMetatable(FileMetaTableName)
	luaState.SetGlobal(FileMetaTableName, fileMetaTable)

	// Set all file metatable functions
	luaState.SetFuncs(fileMetaTable, fileMethods)
}

// CheckFileExists checks if the given file exists
func CheckFileExists(L *lua.LState) int {
	// Get file info
	_, err := os.Stat(L.ToString(2))

	// If file exists push true
	if err == nil {
		L.Push(lua.LBool(true))
		return 1
	}

	// Check if file does not exists then push false
	if os.IsNotExist(err) {
		L.Push(lua.LBool(false))
		return 1
	}

	L.RaiseError("Cannot get file information: %v", err)

	return 0
}

// GetFileModTime gets the file unix timestamp
func GetFileModTime(L *lua.LState) int {
	// Get file info
	info, err := os.Stat(L.ToString(2))

	if err != nil {
		L.RaiseError("Cannot get file information: %v", err)
	}

	// Push file mod time
	L.Push(lua.LNumber(info.ModTime().Unix()))

	return 1
}
