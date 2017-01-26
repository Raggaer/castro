package models

import "time"

// CsrfToken struct used for the application XSRF tokens
type CsrfToken struct {
	Token string
	At    time.Time
}
