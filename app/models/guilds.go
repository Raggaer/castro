package models

import (
	"github.com/raggaer/castro/app/database"
)

type Guild struct {
	ID           int64
	Name         string
	Ownerid      int64
	Creationdata int64
	Motd         string
}

// GetGuildByID retrieves a guild by its identifier
func GetGuildByID(id int64) (*Guild, error) {
	// Data holder
	g := &Guild{}

	//Retrieve guild basic information
	if err := database.DB.Get(g, "SELECT id, ownerid, motd, name creationdata FROM guilds WHERE id = ?", id); err != nil {
		return nil, err
	}

	return g, nil
}

// GetGuildByPlayerID retrieves a player guild
func GetGuildByPlayerID(id int64) (*Guild, error) {
	// Data holder
	g := &Guild{}

	//Retrieve guild basic information
	if err := database.DB.Get(g, "SELECT a.id, a.ownerid, a.motd, a.name, a.creationdata FROM guilds a, guild_membership b WHERE b.player_id = ? AND b.guild_id = a.id", id); err != nil {
		return nil, err
	}

	return g, nil
}

// GetGuildByName retrieves a guild by its name
func GetGuildByName(name string) (*Guild, error) {
	// Data holder
	g := &Guild{}

	//Retrieve guild basic information
	if err := database.DB.Get(g, "SELECT id, ownerid, motd, name creationdata FROM guilds WHERE name = ?", name); err != nil {
		return nil, err
	}

	return g, nil
}
