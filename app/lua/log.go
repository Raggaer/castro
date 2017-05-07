package lua

import (
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
)

// SetLogMetaTable sets the log metatable for the given state
func SetLogMetaTable(luaState *lua.LState) {
	// Create and set the log metatable
	logMetaTable := luaState.NewTypeMetatable(LogMetaTableName)
	luaState.SetGlobal(LogMetaTableName, logMetaTable)

	// Set all mail metatable functions
	luaState.SetFuncs(logMetaTable, logMethods)
}

// LogError logs a message with the error level
func LogError(L *lua.LState) int {
	// Get content to log
	content := L.Get(2)

	// Log content
	util.Logger.Logger.Errorf("%v", content.String())

	return 0
}

// LogFatal logs a message with the fatal level
func LogFatal(L *lua.LState) int {
	// Get content to log
	content := L.Get(2)

	// Log content
	util.Logger.Logger.Fatalf("%v", content.String())

	return 0
}

// LogInfo logs a message with the info level
func LogInfo(L *lua.LState) int {
	// Get content to log
	content := L.Get(2)

	// Log content
	util.Logger.Logger.Infof("%v", content.String())

	return 0
}
