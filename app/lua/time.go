package lua

import (
	"fmt"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"time"
)

// LuaDate lua object date converted from go
type LuaDate struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Second int
	Result string
}

// SetTimeMetaTable sets the time metatable of the given state
func SetTimeMetaTable(luaState *lua.LState) {
	// Create and set the time metatable
	timeMetaTable := luaState.NewTypeMetatable(TimeMetaTableName)
	luaState.SetGlobal(TimeMetaTableName, timeMetaTable)

	// Set all time metatable functions
	luaState.SetFuncs(timeMetaTable, timeMethods)
}

// ParseUnixTimestamp returns a date object for the given timestamp
func ParseUnixTimestamp(L *lua.LState) int {
	// Get timestamp
	unix := L.Get(2)

	// Check for valid timestamp type
	if unix.Type() != lua.LTNumber {
		L.ArgError(1, "Invalid timestamp type. Expected number")
		return 0
	}

	// Get time as int64
	unix64 := L.ToInt64(2)

	// Check if time table is saved on cache
	t, found := util.Cache.Get(
		fmt.Sprintf("time_result_%v", unix64),
	)

	if found {

		// Push date table
		L.Push(t.(*lua.LTable))

		return 1
	}

	// If timestamp is empty return empty lua date
	if unix64 == 0 {

		// Push empty struct
		ldate := LuaDate{
			Result: "Never",
		}

		// Push result as table
		L.Push(StructToTable(&ldate))

		return 1
	}

	// Parse unix timestamp
	d := time.Unix(unix64, 0)

	// Create lua date struct
	ldate := LuaDate{
		Year:   d.Year(),
		Month:  int(d.Month()),
		Day:    d.Day(),
		Hour:   d.Hour(),
		Minute: d.Minute(),
		Second: d.Second(),
		Result: d.Format("Mon Jan 2 15:04:05"),
	}

	// Convert date struct to table
	ldateTable := StructToTable(&ldate)

	// Save date struct to cache
	util.Cache.Add(
		fmt.Sprintf("time_result_%v", unix64),
		ldateTable,
		time.Minute*3,
	)

	// Push result as table
	L.Push(ldateTable)

	return 1
}

// ParseDurationString parses the given duration string to return the time in seconds
func ParseDurationString(L *lua.LState) int {
	// Get duration string
	dur := L.Get(2)

	// Check valid duration
	if dur.Type() != lua.LTString {
		L.ArgError(1, "Invalid time duration string. Expected string")
		return 0
	}

	// Parse duration string
	duration, err := time.ParseDuration(dur.String())

	if err != nil {
		L.RaiseError("Cannot parse duration: %v", err)
		return 0
	}

	// Push duration seconds
	L.Push(lua.LNumber(duration.Seconds()))

	return 1
}
