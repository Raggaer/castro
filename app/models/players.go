package models

import (
	"time"

	"github.com/raggaer/castro/app/database"
)

type PlayerColumn struct {
	Name string
}

// Player struct used for server players
type Player struct {
	ID         int64
	Name       string
	Level      int
	Sex        int
	Vocation   int
	Town_id    uint32
	Account_id int64
	Experience int
}

// GetPlayerByID returns a player by the identifier
func GetPlayerByID(id int64) (*Player, error) {
	// Data holder
	p := &Player{}

	if err := database.DB.Get(p, "SELECT id, sex, account_id, name, level, vocation, town_id FROM players WHERE id = ?", id); err != nil {
		return nil, err
	}

	return p, nil
}

// GetPlayerByName returns a player by the name
func GetPlayerByName(name string) (*Player, error) {
	// Data holder
	p := &Player{}

	if err := database.DB.Get(p, "SELECT id, sex, account_id, name, level, vocation, town_id FROM players WHERE name = ?", name); err != nil {
		return nil, err
	}

	return p, nil
}

// GetBalance returns the player balance
func (p *Player) GetBalance() (int, error) {
	// Data holder
	balance := 0

	// Get balance value
	if err := database.DB.Get(&balance, "SELECT balance FROM players WHERE id = ?", p.ID); err != nil {
		return 0, err
	}

	return balance, nil
}

// SetBalance updates a player balance
func (p *Player) SetBalance(balance int) error {
	_, err := database.DB.Exec("UPDATE players SET balance = ? WHERE id = ?", balance, p.ID)
	return err
}

// IsOnline checks if the player is online
func (p *Player) IsOnline() (bool, error) {
	// Data holder
	online := false

	// Get online value
	if err := database.DB.Get(&online, "SELECT 1 FROM players_online WHERE player_id = ?", p.ID); err != nil {
		return false, err
	}

	return online, nil
}

// GetStorageValue returns a storage value by its key
func (p *Player) GetStorageValue(key int) (*Storage, error) {
	// Data holder
	storage := &Storage{}

	// Get storage value
	if err := database.DB.Get(storage, "SELECT key, value FROM players_storage WHERE player_id = ? AND key= ?", p.ID, key); err != nil {
		return nil, err
	}

	return storage, nil
}

// SetStorageValue sets a player storage value
func (p *Player) SetStorageValue(key, value int) error {
	_, err := database.DB.Exec("INSERT INTO player_storage (player_id, key, value) VALUES (?, ?, ?)", p.ID, key, value)
	return err
}

// GetPremiumTime returns the player premium time end
func (p *Player) GetPremiumEndsAt() (int, error) {
	// Premium time holder
	endsAt := 0

	// Get premium time
	if err := database.DB.Get(&endsAt, "SELECT premium_ends_at FROM accounts WHERE id = ?", p.Account_id); err != nil {
		return 0, err
	}

	return endsAt, nil
}

// GetPremiumTime returns the player premium time end
func (p *Player) GetPremiumTime() (int, error) {
	// End time holder
	endsAt, err := p.GetPremiumEndsAt()
	if err != nil {
		return 0, err
	}

	// Calculate remaining time (seconds)
	timeLeft := time.Unix(int64(endsAt), 0).Sub(time.Now()).Seconds()
	if timeLeft <= 0 {
		return 0, nil
	}

	return int(timeLeft), nil
}

// GetPremiumDays returns the player premium days
func (p *Player) GetPremiumDays() (int, error) {
	// Premium time holder
	timeLeft, err := p.GetPremiumTime()
	if err != nil {
		return 0, err
	}

	return int(timeLeft / 86400), nil
}

// GetExperience returns the player
func (p *Player) GetExperience() (int, error) {
	// Experience placeholder
	experience := 0

	// Retrieve experience from database
	if err := database.DB.Get(&experience, "SELECT experience FROM players WHERE id = ?", p.ID); err != nil {
		return 0, err
	}

	return experience, nil
}

// GetCapacity returns the player capacity
func (p *Player) GetCapacity() (int, error) {
	// Capacity placeholder
	capacity := 0

	// Retrieve experience from database
	if err := database.DB.Get(&capacity, "SELECT capacity FROM players WHERE id = ?", p.ID); err != nil {
		return 0, err
	}

	return capacity, nil
}
