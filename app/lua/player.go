package lua

import (
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
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
		if err := database.DB.Get(&player, "SELECT id, accound_id, name, level, vocation, town_id FROM players WHERE id = ?", L.ToInt64(1)); err != nil {
			L.RaiseError("Cannot get player for ID: %v. Error: %v", L.ToInt64(1), err)
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

	// Get player by ID
	if err := database.DB.Get(&player, "SELECT id, account_id, name, level, vocation, town_id FROM players WHERE name = ?", L.ToString(1)); err != nil {
		L.RaiseError("Cannot get player for name: %v. Error: %v", L.ToInt64(1), err)
	}

	// Create player metatable
	L.Push(createPlayerMetaTable(&player, L))

	return 1
}

// createPlayerMetaTable returns a new player metatable
func createPlayerMetaTable(player *models.Player, luaState *lua.LState) *lua.LTable {
	// Create a player metatable
	playerMetaTable := luaState.NewTypeMetatable(PlayerMetaTableName)

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
	tbl := luaState.GetTypeMetatable(PlayerMetaTableName)

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

// IsPlayerOnline checks if the given player is online
func IsPlayerOnline(L *lua.LState) int {
	// Get player struct
	player := getPlayerObject(L)

	// Data holder
	online := false

	// Get online value
	database.DB.Get(&online, "SELECT 1 FROM players_online WHERE player_id = ?", player.ID)

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
	database.DB.Get(&storage, "SELECT key, value FROM players_storage WHERE player_id = ?", player.ID)

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
	if _, err := database.DB.Exec("INSERT INTO player_storage (player_id, key, value) VALUES (?, ?, ?)", player.ID, L.ToInt(2), L.ToInt(3)); err != nil {
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
