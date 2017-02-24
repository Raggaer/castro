package lua

import (
	"github.com/yuin/gopher-lua"
	"time"
)

// SetEventMetaTable sets the event metatable of the given state
func SetEventMetaTable(luaState *lua.LState) {
	// Create and set the crypto metatable
	eventMetaTable := luaState.NewTypeMetatable(EventMetaTableName)
	luaState.SetGlobal(EventMetaTableName, eventMetaTable)

	// Set all event metatable functions
	luaState.SetFuncs(eventMetaTable, eventMethods)
}

// AddEvent adds a background job
func AddEvent(L *lua.LState) int {
	// Get function
	f := L.Get(2)

	// Check valid function type
	if f.Type() != lua.LTFunction {
		L.ArgError(1, "Invalid function type. Expected function")
		return 0
	}

	// Get duration string
	durationStr := L.Get(3)

	// Check valid duration type
	if durationStr.Type() != lua.LTString {
		L.ArgError(2, "Invalid duration type. Expected string")
		return 0
	}

	// Parse duration string
	duration, err := time.ParseDuration(durationStr.String())

	if err != nil {
		L.RaiseError("Cannot parse event duration: %v", err)
		return 0
	}

	go executeEvent(f.(*lua.LFunction), duration)

	return 0
}

// executeEvents runs the lua background event
func executeEvent(f *lua.LFunction, duration time.Duration) {
	// Create a state
	state := lua.NewState()

	// Close state
	defer state.Close()

	// Push function
	state.Push(f)

	// Create ticker
	tick := time.NewTicker(duration)

	// Close ticker
	defer tick.Stop()

	// Run background task
	for {
		select {
		case <-tick.C:

			// Execute function
			state.Call(0, 0)

			// Push function
			state.Push(f)
		}
	}
}
