package models

import "time"

type CsrfToken struct {
	Token string
	At    time.Time
}
