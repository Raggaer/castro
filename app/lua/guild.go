package lua

import (
	"reflect"

	"github.com/kataras/go-errors"
	"github.com/raggaer/castro/app/models"
	"github.com/yuin/gopher-lua"
)

// GuildConstructor returns a new guild metatable for the given ID or name
func GuildConstructor(L *lua.LState) int {
	// Retrieve guild
	guild, err := guildTableConstructor(L.Get(1))

	if err != nil {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(createGuildMetaTable(guild, L))
	return 1
}

func guildTableConstructor(i interface{}) (*models.Guild, error) {
	// Get guild by ID
	if reflect.TypeOf(i).Kind() == reflect.Int64 {
		return models.GetGuildByID(i.(int64))
	}

	if reflect.TypeOf(i).Kind() != reflect.String {
		return nil, errors.New("Invalid guild name or id")
	}

	// Get guild by name
	return models.GetGuildByName(i.(string))
}

func createGuildMetaTable(guild *models.Guild, luaState *lua.LState) *lua.LTable {
	// Create a guild metatable
	guildMetaTable := luaState.NewTable()

	// Set user data
	u := luaState.NewUserData()

	// Set user data value
	u.Value = guild

	// Set user data field
	luaState.SetField(guildMetaTable, "__guild", u)

	// Set all player metatable functions
	luaState.SetFuncs(guildMetaTable, guildMethods)

	// Set guild public fields
	MergeTableFields(StructToTable(guild), guildMetaTable)

	return guildMetaTable
}

func getGuildObject(luaState *lua.LState) *models.Guild {
	// Get metatable
	tbl := luaState.ToTable(1)

	// Get user data field
	data := luaState.GetField(tbl, "__guild").(*lua.LUserData)

	// Return user data as pointer to struct
	return data.Value.(*models.Guild)
}

// GetGuildOwner retrieves a guild owner
func GetGuildOwner(L *lua.LState) int {
	// Retrieve guild object
	guild := getGuildObject(L)

	// Push owner as player metatable
	L.Push(lua.LNumber(guild.Ownerid))
	return 1
}
