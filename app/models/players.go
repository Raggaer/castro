package models

// Player struct used for server players
type Player struct {
	ID         int64
	Name       string
	Level      int
	Vocation   int
	Town_id    int
	Account_id int64
}
