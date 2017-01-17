package lua

import (
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"strings"
	"time"
)

// SetDatabaseMetaTable sets the database metatable of the given state
func SetDatabaseMetaTable(luaState *lua.LState) {
	// Create and set the MySQL metatable
	mysqlMetaTable := luaState.NewTypeMetatable(MySQLMetaTableName)
	luaState.SetGlobal(MySQLMetaTableName, mysqlMetaTable)

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

	// Execute the query
	if err := database.DB.Exec(query.String(), args...).Error; err != nil {

		L.RaiseError("Cannot execute query: %v", err)
		return 0
	}

	return 0
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
	cache := L.Get(3 + n)

	// Save cache variable
	saveToCache := false
	cacheKey := query.String()

	if cache.Type() != lua.LTNil {

		// Build cache key as the full query with arguments inside
		for _, arg := range args {

			// Replace "?" with argument
			cacheKey = strings.Replace(cacheKey, "?", arg.(string), 1)
		}

		// Try to load from cache
		q, found := util.Cache.Get(cacheKey)

		if found {

			// Push the cached lua table to stack
			L.Push(q.(*lua.LTable))

			return 1
		}

		saveToCache = true
	}

	// Execute query and get rows
	rows, err := database.DB.Raw(query.String(), args...).Rows()
	if err != nil {

		// Raise error
		L.RaiseError("Cannot get result: %v", err)
		return 0
	}

	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {

		// Raise error
		L.RaiseError("Cannot get article: missing QUERY")
		return 0
	}
	var results [][]interface{}

	// Loop result rows
	for rows.Next() {

		// Hold all result columns
		columns := make([]interface{}, len(columnNames))

		// Hold all result pointers
		columnPointers := make([]interface{}, len(columnNames))

		// Loop result columns and assign to pointer
		for i := range columnNames {
			columnPointers[i] = &columns[i]
		}

		// Scan the column pointers
		rows.Scan(columnPointers...)

		// Append to results
		results = append(results, columns)
	}

	// If there are no query results push nil
	if len(results) <= 0 {

		L.Push(lua.LNil)
		return 1
	}

	finalTable := QueryToTable(results, columnNames)

	// If there is only one result do not
	// return as table of results
	if len(results) == 1 {

		// Get first element of the result table
		finalTable = finalTable.RawGetInt(1).(*lua.LTable)
	}

	// If user wants to use cache save table
	if saveToCache {
		util.Cache.Add(cacheKey, finalTable, time.Minute*3)
	}

	// Push the converted query to the stack
	L.Push(finalTable)

	return 1
}
