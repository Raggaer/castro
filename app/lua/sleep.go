package lua

import (
	"time"

	"github.com/yuin/gopher-lua"
)

// ThreadSleep sleeps the current running lua state
func ThreadSleep(L *lua.LState) int {
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
