package models

// Player struct used for server players
type Player struct {
	ID         int64
	Name       string
	Level      int
	Sex        int
	Vocation   int
	Town_id    uint32
	Account_id int64
}
