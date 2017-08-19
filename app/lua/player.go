package lua

import (
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"html"
)

// PlayerConstructor returns a new player metatable for the given ID or name
func PlayerConstructor(L *lua.LState) int {
	// Get name or ID
	i := L.Get(1)

	// Get player by ID
	if i.Type() == lua.LTNumber {

		// Data holder
		player := models.Player{}

		// Get player by ID
		if err := database.DB.Get(&player, "SELECT id, sex, accound_id, name, level, vocation, town_id FROM players WHERE id = ?", L.ToInt64(1)); err != nil {
			L.Push(lua.LNil)
			return 1
		}

		// Create player metatable
		L.Push(createPlayerMetaTable(&player, L))

		return 1
	}

	if i.Type() != lua.LTString {
		L.ArgError(1, "Invalid player ID or name")
		return 0
	}

	// Data holder
	player := models.Player{}

	// Get player by name
	if err := database.DB.Get(&player, "SELECT id, sex, account_id, name, level, vocation, town_id FROM players WHERE name = ?", L.ToString(1)); err != nil {
		L.Push(lua.LNil)
		return 1
	}

	// Create player metatable
	L.Push(createPlayerMetaTable(&player, L))

	return 1
}

// createPlayerMetaTable returns a new player metatable
func createPlayerMetaTable(player *models.Player, luaState *lua.LState) *lua.LTable {
	// Create a player metatable
	playerMetaTable := luaState.NewTable()

	// Set user data
	u := luaState.NewUserData()

	// Set user data value
	u.Value = player

	// Set user data field
	luaState.SetField(playerMetaTable, "__player", u)

	// Set all player metatable functions
	luaState.SetFuncs(playerMetaTable, playerMethods)

	return playerMetaTable
}

// getPlayerObject returns a player struct from a user data value
func getPlayerObject(luaState *lua.LState) *models.Player {
	// Get metatable
	tbl := luaState.ToTable(1)

	// Get user data field
	data := luaState.GetField(tbl, "__player").(*lua.LUserData)

	// Return user data as pointer to struct
	return data.Value.(*models.Player)
}

// GetPlayerAccountID gets a player account ID
func GetPlayerAccountID(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Push account ID
	L.Push(lua.LNumber(player.Account_id))

	return 1
}

// GetPlayerBankBalance gets a player bank balance
func GetPlayerBankBalance(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Data holder
	balance := 0

	// Get balance value
	if err := database.DB.Get(&balance, "SELECT balance FROM players WHERE id = ?", player.ID); err != nil {
		L.RaiseError("Cannot get player bank balance: %v", err)
		return 0
	}

	// Push value
	L.Push(lua.LNumber(balance))

	return 1
}

// SetPlayerBankBalance sets a player bank balance
func SetPlayerBankBalance(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Retrieve bank balance number
	newBalance := L.ToInt(2)

	// Update bank balance
	if _, err := executeQueryHelper(L, "UPDATE players SET balance = ? WHERE id = ?", newBalance, player.ID); err != nil {
		L.RaiseError("Cannot update player balance: %v")
		return 0
	}

	return 0
}

// IsPlayerOnline checks if the given player is online
func IsPlayerOnline(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Data holder
	online := false

	// Get online value
	if err := database.DB.Get(&online, "SELECT 1 FROM players_online WHERE player_id = ?", player.ID); err != nil {
		L.RaiseError("Cannot check if player is online: %v", err)
		return 0
	}

	// Push online value
	L.Push(lua.LBool(online))

	return 1
}

// GetPlayerStorageValue gets a player storage value by the given key
func GetPlayerStorageValue(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Get key
	key := L.Get(2)

	// Check for valid key type
	if key.Type() != lua.LTNumber {
		L.ArgError(1, "Invalid key type. Expected number")
		return 0
	}

	// Data holder
	storage := models.Storage{}

	// Get storage value
	if err := database.DB.Get(&storage, "SELECT key, value FROM players_storage WHERE player_id = ?", player.ID); err != nil {
		L.RaiseError("Cannot get player storage: %v", err)
		return 0
	}

	// Push storage as table
	L.Push(StructToTable(&storage))

	return 1
}

// SetPlayerStorageValue sets a player storage value with the given key
func SetPlayerStorageValue(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Get key
	key := L.Get(2)

	// Check for valid key type
	if key.Type() != lua.LTNumber {
		L.ArgError(1, "Invalid key type. Expected number")
		return 0
	}

	// Get value
	val := L.Get(3)

	// Check for valid value type
	if val.Type() != lua.LTNumber {
		L.ArgError(1, "Invalid value type. Expected number")
		return 0
	}

	// Insert storage value
	if _, err := executeQueryHelper(L, "INSERT INTO player_storage (player_id, key, value) VALUES (?, ?, ?)", player.ID, L.ToInt(2), L.ToInt(3)); err != nil {
		L.RaiseError("Cannot set player storage value: %v", err)
		return 0
	}

	return 0
}

// GetPlayerVocation gets the player vocation
func GetPlayerVocation(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Loop server vocations
	for _, voc := range util.ServerVocationList.List.Vocations {

		// Check vocation
		if voc.ID == player.Vocation {

			// Convert vocation to lua table
			L.Push(StructToTable(voc))

			return 1
		}
	}

	// Vocation is not found
	L.RaiseError("Cannot find player vocation")

	return 0
}

// GetPlayerGender gets the player gender
func GetPlayerGender(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Push gender as number
	L.Push(lua.LNumber(player.Sex))

	return 1
}

// GetPlayerPremiumDays gets the player number of premium days
func GetPlayerPremiumDays(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Premium days holder
	days := 0

	// Get premium days
	if err := database.DB.Get(&days, "SELECT premdays FROM accounts WHERE id = ?", player.Account_id); err != nil {
		L.RaiseError("Cannot get player premium days: %v", err)
		return 0
	}

	// Push days as number
	L.Push(lua.LNumber(days))

	return 1
}

// GetPlayerTown gets the player town
func GetPlayerTown(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Loop towns
	for _, town := range util.OTBMap.Map.Towns {

		// Check for player town
		if town.ID == player.Town_id {

			// Push town as table
			L.Push(StructToTable(&town))

			return 1
		}
	}

	L.RaiseError("Cannot find player town")

	return 0
}

// GetPlayerLevel gets the player level
func GetPlayerLevel(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Push player level as number
	L.Push(lua.LNumber(player.Level))

	return 1
}

// GetPlayerName gets the player name
func GetPlayerName(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Push player name as string
	L.Push(lua.LString(player.Name))

	return 1
}

// GetPlayerExperience gets the player experience
func GetPlayerExperience(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Experience placeholder
	experience := 0

	// Retrieve experience from database
	if err := database.DB.Get(&experience, "SELECT experience FROM players WHERE id = ?", player.ID); err != nil {
		L.RaiseError("Cannot get player experience from database: %v", err)
		return 0
	}

	// Push player experience as number
	L.Push(lua.LNumber(experience))

	return 1
}

// GetPlayerCapacity gets the player capacity
func GetPlayerCapacity(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Capacity placeholder
	capacity := 0

	// Retrieve experience from database
	if err := database.DB.Get(&capacity, "SELECT capacity FROM players WHERE id = ?", player.ID); err != nil {
		L.RaiseError("Cannot get player capacity from database: %v", err)
		return 0
	}

	// Push player capacity as number
	L.Push(lua.LNumber(capacity))

	return 1
}

// GetPlayerCustomField retrieves a field from the player table as string
func GetPlayerCustomField(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Get field name
	fieldName := L.ToString(2)

	// Field placeholder
	fieldValue := ""

	// Retrieve current schema
	schema := Config.GetGlobal("mysqlDatabase").String()

	// Column name placeholder
	nameList := []models.PlayerColumn{}

	// Get all player column names
	if err := database.DB.Select(&nameList, "SELECT COLUMN_NAME AS name FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = ? AND TABLE_SCHEMA = ?", "players", schema); err != nil {
		L.RaiseError("Cannot get list of column names from information_schema: %v", err)
		return 0
	}

	// Loop column list
	for _, column := range nameList {

		// Check for valid column name
		if column.Name == fieldName {

			// Retrieve custom field
			if err := database.DB.Get(&fieldValue, "SELECT "+html.EscapeString(fieldName)+" FROM players WHERE id = ?", player.ID); err != nil {
				L.RaiseError("Cannot get custom field %s: %v", fieldName, err)
				return 0
			}

			// Push value as string
			L.Push(lua.LString(fieldValue))

			return 1
		}
	}

	// Push nil if the field is not valid
	L.Push(lua.LNil)

	return 1
}
