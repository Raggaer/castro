package lua

import (
	"github.com/yuin/gopher-lua"
	"time"
)

// LuaSleep sleeps the current running lua state
func LuaSleep(L *lua.LState) int {
	// Get sleep duration
	duration, err := time.ParseDuration(L.ToString(1))

	if err != nil {
		L.RaiseError("Cannot parse sleep duration: %v", err)
		return 0
	}

	// Sleep goroutine
	time.Sleep(duration)

	return 0
}
