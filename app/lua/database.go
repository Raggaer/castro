package lua

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
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

	// Set transaction field status
	luaState.SetField(mysqlMetaTable, DatabaseTransactionStatusFieldName, lua.LBool(false))
}

// GetDatabaseTransactionFieldStatus retrieves the transaction status field
func GetDatabaseTransactionFieldStatus(luaState *lua.LState) bool {
	// Get transaction field
	txFieldStatus, ok := luaState.GetField(luaState.GetTypeMetatable(DatabaseMetaTableName), DatabaseTransactionStatusFieldName).(lua.LBool)

	if !ok {
		luaState.RaiseError("Cannot retrieve transaction status field")
	}

	return bool(txFieldStatus)
}

func GetDatabaseTransactionField(luaState *lua.LState) *sqlx.Tx {
	// Get metatable
	m := luaState.GetTypeMetatable(DatabaseMetaTableName)

	// Get transaction user data
	data, ok := luaState.GetField(m, DatabaseTransactionFieldName).(*lua.LUserData)

	if !ok {
		luaState.RaiseError("Cannot retrieve database transaction user data")
	}

	// Get database transaction
	tx, ok := data.Value.(*sqlx.Tx)

	if !ok {
		luaState.RaiseError("Cannot retrieve database transaction from user data")
	}

	return tx
}

func setDatabaseTransactionField(luaState *lua.LState) *sqlx.Tx {
	// Create database transaction
	tx, err := database.DB.Beginx()

	if err != nil {
		luaState.RaiseError("Cannot begin database transaction: %v", err)
	}

	// Set transaction field status
	luaState.SetField(luaState.GetTypeMetatable(DatabaseMetaTableName), DatabaseTransactionStatusFieldName, lua.LBool(true))

	// Create transaction user data
	txUserData := luaState.NewUserData()
	txUserData.Value = tx

	// Set transaction field
	luaState.SetField(luaState.GetTypeMetatable(DatabaseMetaTableName), DatabaseTransactionFieldName, txUserData)

	return tx
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
		util.Logger.Logger.Infof("execute: "+strings.Replace(query.String(), "?", "%v", -1), args...)
	}

	// Result and error placeholders
	var result sql.Result
	var err error

	// Retrieve transaction status
	txStatus := GetDatabaseTransactionFieldStatus(L)

	if txStatus {

		// Retrieve transaction field
		tx := GetDatabaseTransactionField(L)

		// Execute query
		result, err = tx.Exec(query.String(), args...)

	} else {

		// Update transaction status
		tx := setDatabaseTransactionField(L)

		// Execute query
		result, err = tx.Exec(query.String(), args...)
	}

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
		util.Logger.Logger.Infof("query: "+strings.Replace(query.String(), "?", "%v", -1), args...)
	}

	// Run query
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
		util.Cache.Add(cacheKey, results, util.Config.Configuration.Cache.Default.Duration)
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
		util.Logger.Logger.Infof("query: "+strings.Replace(query.String(), "?", "%v", -1), args...)
	}

	// Run query
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
		util.Cache.Add(cacheKey, results, util.Config.Configuration.Cache.Default.Duration)
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
