package util

import (
	c "github.com/patrickmn/go-cache"
	"net/http"
)

// Cache variable that holds the main cache instance
// of the applicattion
var Cache *c.Cache

// GetSession returns an user session by its
// cookie ID
func GetSession(req *http.Request) {

}