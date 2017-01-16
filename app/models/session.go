package models

import "time"

// Session struct used to store user sessions
// on the database
type Session struct {
	ID        int64
	Token     string
	Data      []byte `gorm:"type:blob"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
