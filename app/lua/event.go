package lua

import (
	"github.com/yuin/gopher-lua"
)

// SetEventsMetaTable sets the event metatable of the given state
func SetEventsMetaTable(luaState *lua.LState) {
	// Create and set the events metatable
	eventMetaTable := luaState.NewTypeMetatable(EventsMetaTableName)
	luaState.SetGlobal(EventsMetaTableName, eventMetaTable)

	// Set all events metatable functions
	luaState.SetFuncs(eventMetaTable, eventsMethods)
}

// LuaEvent executes a background event
func LuaEvent(L *lua.LState) int {
	// Get function
	f := L.ToFunction(2)

	// Create new thread
	thread, _ := L.NewThread()

	// Close thread
	defer thread.Close()

	// Infinite loop
	go func() {

		for {

			// Resume function using  a new state thread
			status, err, _ := L.Resume(thread, f)

			// Check if event finished execution
			if status == lua.ResumeOK {
				break
			}

			if err != nil {

				L.RaiseError("Running event returned an error: %v", err)
				break
			}
		}
	}()

	return 0
}
