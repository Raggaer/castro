package lua

import (
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/models"
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
	database.DB.Get(&balance, "SELECT balance FROM players WHERE id = ?", player.ID)

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
