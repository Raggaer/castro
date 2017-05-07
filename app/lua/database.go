package lua

import (
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"strings"
)

// SetDatabaseMetaTable sets the database metatable of the given state
func SetDatabaseMetaTable(luaState *lua.LState) {
	// Create and set the MySQL metatable
	mysqlMetaTable := luaState.NewTypeMetatable(DatabaseMetaTableName)
	luaState.SetGlobal(DatabaseMetaTableName, mysqlMetaTable)

	// Set all MySQL metatable functions
	luaState.SetFuncs(mysqlMetaTable, mysqlMethods)
}

// Execute executes a query without returning the result
func Execute(L *lua.LState) int {
	// Get query
	query := L.Get(2)

	// Check if query is valid
	if query.Type() != lua.LTString {

		// Raise error
		L.ArgError(1, "Invalid query type. Expected string")
		return 0
	}

	// Count number of params
	n := strings.Count(query.String(), "?")

	args := []interface{}{}

	// Get all arguments matching the number of params
	for i := 0; i < n; i++ {

		// Append argument to slice
		args = append(args, L.Get(3+i).String())
	}

	// Log query on development mode
	if util.Config.Configuration.IsDev() || util.Config.Configuration.IsLog() {
		util.Logger.Infof("execute: "+strings.Replace(query.String(), "?", "%v", -1), args...)
	}

	// Execute query
	result, err := database.DB.Exec(query.String(), args...)

	if err != nil {
		L.RaiseError("Cannot execute query: %v", err)
		return 0
	}

	// Check if query is INSERT
	if !strings.HasPrefix(query.String(), "INSERT") {
		return 0
	}

	// Get last inserted id
	id, err := result.LastInsertId()

	if err != nil {
		L.RaiseError("Cannot get last inserted id: %v", err)
		return 0
	}

	// Push id
	L.Push(lua.LNumber(id))

	return 1
}

// SingleQuery performs and ad-hoc query returning the first result as a table
func SingleQuery(L *lua.LState) int {
	// Get query
	query := L.Get(2)

	// Check if query is valid
	if query.Type() != lua.LTString {

		// Raise error
		L.ArgError(1, "Invalid query type. Expected string")
		return 0
	}

	// Count number of params
	n := strings.Count(query.String(), "?")

	args := []interface{}{}

	// Get all arguments matching the number of params
	for i := 0; i < n; i++ {

		// Append argument to slice
		args = append(args, L.Get(3+i).String())
	}

	// Check if user wants to use cache
	cache := L.ToBool(3 + n)

	// Save cache variable
	saveToCache := false
	cacheKey := query.String()

	if cache {

		// Build cache key as the full query with arguments inside
		for _, arg := range args {

			// Replace "?" with argument
			cacheKey = strings.Replace(cacheKey, "?", arg.(string), 1)
		}

		// Try to load from cache
		q, found := util.Cache.Get(cacheKey)

		if found {

			results := q.(*lua.LTable)

			// Set cache status
			L.Push(lua.LBool(true))

			// If there are no results return nil
			if results.Len() == 0 {
				L.Push(lua.LNil)

				// Set cache status
				L.Push(lua.LBool(true))

				return 2
			}

			// Push the cached lua table to stack
			L.Push(results.RawGetInt(1))

			// Set cache status
			L.Push(lua.LBool(true))

			return 2
		}

		saveToCache = true
	}

	// Log query on development mode
	if util.Config.Configuration.IsDev() || util.Config.Configuration.IsLog() {
		util.Logger.Infof("query: "+strings.Replace(query.String(), "?", "%v", -1), args...)
	}

	// Execute query
	rows, err := database.DB.Queryx(query.String(), args...)

	if err != nil {
		L.RaiseError("Cannot execute query: %v", err)
		return 0
	}

	// Close rows
	defer rows.Close()

	// Result holder
	results := L.NewTable()

	// Loop rows
	for rows.Next() {

		// Hold current row
		result := make(map[string]interface{})

		// Scan row to map
		if err := rows.MapScan(result); err != nil {
			L.RaiseError("Cannot map row to map: %v", err)
			return 0
		}

		// Append to lua table
		results.Append(MapToTable(result))
	}

	// If user wants to use cache save table
	if saveToCache {
		util.Cache.Add(cacheKey, results, util.Config.Configuration.Cache.Default)
	}

	// If there are no results return nil
	if results.Len() == 0 {
		L.Push(lua.LNil)

		// Set cache status
		L.Push(lua.LBool(false))

		return 2
	}

	// Push result
	L.Push(results.RawGetInt(1))

	// Set cache status
	L.Push(lua.LBool(false))

	return 2
}

// Query executes an ad-hoc query
func Query(L *lua.LState) int {
	// Get query
	query := L.Get(2)

	// Check if query is valid
	if query.Type() != lua.LTString {

		// Raise error
		L.ArgError(1, "Invalid query type. Expected string")
		return 0
	}

	// Count number of params
	n := strings.Count(query.String(), "?")

	args := []interface{}{}

	// Get all arguments matching the number of params
	for i := 0; i < n; i++ {

		// Append argument to slice
		args = append(args, L.Get(3+i).String())
	}

	// Check if user wants to use cache
	cache := L.ToBool(3 + n)

	// Save cache variable
	saveToCache := false
	cacheKey := query.String()

	if cache {

		// Build cache key as the full query with arguments inside
		for _, arg := range args {

			// Replace "?" with argument
			cacheKey = strings.Replace(cacheKey, "?", arg.(string), 1)
		}

		// Try to load from cache
		q, found := util.Cache.Get(cacheKey)

		if found {

			results := q.(*lua.LTable)

			// Set cache status
			L.Push(lua.LBool(true))

			// If there are no results return nil
			if results.Len() == 0 {
				L.Push(lua.LNil)

				// Set cache status
				L.Push(lua.LBool(true))

				return 2
			}

			// Push the cached lua table to stack
			L.Push(results)

			// Set cache status
			L.Push(lua.LBool(true))

			return 2
		}

		saveToCache = true
	}

	// Log query on development mode
	if util.Config.Configuration.IsDev() || util.Config.Configuration.IsLog() {
		util.Logger.Infof("query: "+strings.Replace(query.String(), "?", "%v", -1), args...)
	}

	// Execute query
	rows, err := database.DB.Queryx(query.String(), args...)

	if err != nil {
		L.RaiseError("Cannot execute query: %v", err)
		return 0
	}

	// Close rows
	defer rows.Close()

	// Result holder
	results := L.NewTable()

	// Loop rows
	for rows.Next() {

		// Hold current row
		result := make(map[string]interface{})

		// Scan row to map
		if err := rows.MapScan(result); err != nil {
			L.RaiseError("Cannot map row to map: %v", err)
			return 0
		}

		// Append to lua table
		results.Append(MapToTable(result))
	}

	// If user wants to use cache save table
	if saveToCache {
		util.Cache.Add(cacheKey, results, util.Config.Configuration.Cache.Default)
	}

	// If there are no results return nil
	if results.Len() == 0 {
		L.Push(lua.LNil)

		// Set cache status
		L.Push(lua.LBool(false))

		return 2
	}

	// Push result
	L.Push(results)

	// Set cache status
	L.Push(lua.LBool(false))

	return 2
}
