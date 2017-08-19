package lua

import (
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"time"
)

// getEventChannel retrieves a channel from the user data
func getEventChannel(luaState *lua.LState) chan int {
	// Get the metatable
	table := luaState.GetTypeMetatable("event")

	// Get user data field
	data := luaState.GetField(table, "__channel")

	return data.(*lua.LUserData).Value.(chan int)
}

// SetEventMetaTable sets the event metatable for the given state
func SetEventMetaTable(luaState *lua.LState) {
	// Create and set the event metatable
	eventMetaTable := luaState.NewTypeMetatable(EventMetaTableName)
	luaState.SetGlobal(EventMetaTableName, eventMetaTable)

	// Set all event metatable functions
	luaState.SetFuncs(eventMetaTable, eventMethods)
}

// SetEventsMetaTable sets the event metatable of the given state
func SetEventsMetaTable(luaState *lua.LState) {
	// Create and set the events metatable
	eventMetaTable := luaState.NewTypeMetatable(EventsMetaTableName)
	luaState.SetGlobal(EventsMetaTableName, eventMetaTable)

	// Set all events metatable functions
	luaState.SetFuncs(eventMetaTable, eventsMethods)
}

// AddEvent adds a background job
func AddEvent(L *lua.LState) int {
	// Get function
	f := L.Get(2)

	// Check valid function type
	if f.Type() != lua.LTString {
		L.ArgError(1, "Invalid function type. Expected string")
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

	go executeEvent(f.String(), duration)

	return 0
}

// StopEvent stops the given background event
func StopEvent(L *lua.LState) int {
	// Get channel
	ch := getEventChannel(L)

	// Signal stop
	ch <- 0

	return 0
}

// executeEvents runs the lua background event
func executeEvent(file string, duration time.Duration) {
	// Create a state
	state := lua.NewState()

	// Close state
	defer state.Close()

	GetApplicationState(state)

	// Create event channel
	eventChannel := make(chan int, 100)

	// Close event channel
	defer close(eventChannel)

	// Create ticker
	tick := time.NewTicker(duration)

	// Close ticker
	defer tick.Stop()

	// Get event metatable
	meta := state.GetTypeMetatable(EventMetaTableName)

	// Set event channel user data
	channelUserData := state.NewUserData()
	channelUserData.Value = eventChannel
	state.SetField(meta, "__channel", channelUserData)

	// Push function
	if err := state.DoFile(file); err != nil {
		util.Logger.Logger.Errorf("Cannot run event: %v", err)
		return
	}

	// Run background task
	for {
		select {
		case <-tick.C:

			// Execute function
			state.CallByParam(lua.P{
				Fn:      state.GetGlobal("run"),
				Protect: !util.Config.Configuration.IsDev(),
			})

		case i := <-eventChannel:

			// Stop signal
			if i == 0 {
				return
			}
		}
	}
}
