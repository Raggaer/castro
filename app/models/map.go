package models

import "time"

// Map struct used for the castro encoded map
type Map struct {
	ID         uint64
	Name       string
	Data       []byte
	Created_at time.Time
	Updated_at time.Time
}
