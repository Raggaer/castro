package lua

import (
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"strconv"
	"time"
)

// SetCacheMetaTable sets the cache metatable of the given state
func SetCacheMetaTable(luaState *lua.LState) {
	// Create and set the cache metatable
	cacheMetaTable := luaState.NewTypeMetatable(CacheMetaTableName)
	luaState.SetGlobal(CacheMetaTableName, cacheMetaTable)

	// Set all captcha metatable functions
	luaState.SetFuncs(cacheMetaTable, cacheMethods)
}

// GetCacheValue retrieves a value from the application cache
func GetCacheValue(L *lua.LState) int {
	// Get key
	key := L.Get(2)

	// Check valid key
	if key.Type() != lua.LTString {
		L.ArgError(1, "Invalid cache key type. Expected string")
		return 0
	}

	// Get value from cache
	v, found := util.Cache.Get(key.String())

	// If there is no value return nil
	if !found {
		L.Push(lua.LNil)
		return 1
	}

	// Switch cache value type
	switch v.(type) {

	case string:
		L.Push(lua.LString(v.(string)))
	case float64:
		L.Push(lua.LNumber(v.(float64)))
	case bool:
		L.Push(lua.LBool(v.(bool)))
	case map[string]interface{}:
		L.Push(MapToTable(v.(map[string]interface{})))
	}

	return 1
}

// SetCacheValue sets a cache value with the given key and the given duration string
func SetCacheValue(L *lua.LState) int {
	// Get key
	key := L.Get(2)

	// Check valid key
	if key.Type() != lua.LTString {
		L.ArgError(1, "Invalid cache key type. Expected string")
		return 0
	}

	// Get value
	val := L.Get(3)

	// Check for invalid lua value
	if val.Type() == lua.LTNil {
		L.ArgError(2, "Invalid cache value type")
		return 0
	}

	// Get optional time value
	t := L.Get(4)

	// Duration time placeholder. Cache default time
	dur := util.Config.Configuration.Cache.Default.Duration

	if t.Type() == lua.LTString {

		// Parse time
		d, err := time.ParseDuration(t.String())

		if err != nil {
			L.ArgError(3, "Invalid time format. Unexpected format")
			return 0
		}

		// Set cache placeholder
		dur = d
	}

	// Switch cache value type
	switch val.Type() {

	case lua.LTString:

		util.Cache.Set(key.String(), val.String(), dur)
	case lua.LTNumber:

		// Convert number to float64
		f, err := strconv.ParseFloat(val.String(), 64)

		if err != nil {
			L.ArgError(2, "Invalid cache value type. Expected number")
			return 0
		}

		// Set cache value
		util.Cache.Set(key.String(), f, dur)

	case lua.LTBool:

		// Convert bool to go bool
		b, err := strconv.ParseBool(val.String())

		if err != nil {
			L.ArgError(2, "Invalid cache value type. Expected boolean")
			return 0
		}

		// Set cache value
		util.Cache.Set(key.String(), b, dur)

	case lua.LTTable:

		// Convert table to map
		m := TableToMap(val.(*lua.LTable))

		// Set cache value
		util.Cache.Set(key.String(), m, dur)
	}

	return 0
}

// DeleteCacheValue removes a key from the cache storage
func DeleteCacheValue(L *lua.LState) int {
	// Get cache key
	key := L.Get(2)

	if key.Type() != lua.LTString {
		L.ArgError(1, "Invalid cache key type. Expected string")
		return 0
	}

	// Delete element from the cache
	util.Cache.Delete(key.String())

	return 0
}
