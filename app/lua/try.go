package lua

import (
	lua "github.com/yuin/gopher-lua"
)

// TryCatch implements a try-catch block for lua
func TryCatch(L *lua.LState) int {
	// Get values
	tryfunc := L.ToFunction(1)
	catchfunc := L.ToFunction(2)

	// Create state thread
	th, _ := L.NewThread()

	// Push try function to thread
	th.Push(tryfunc)

	// Protected call try catch function with the catch function as error
	if err := th.PCall(0, 0, catchfunc); err != nil && err.Error() != "nil" {
		L.RaiseError("Unable to protected call try function: %v", err)
	}

	// Close state thread
	th.Close()

	return 0
}
