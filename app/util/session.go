package util

import (
	"github.com/gorilla/securecookie"
	"net/http"
)

// SessionStore main application session storage
var SessionStore *securecookie.SecureCookie

// SessionCookie returns a session cookie pointer
func SessionCookie(v string) *http.Cookie {
	return &http.Cookie{
		Name:     Config.Cookies.Name,
		Value:    v,
		Path:     "/",
		Secure:   Config.IsSSL(),
		MaxAge:   Config.Cookies.MaxAge,
		HttpOnly: true,
	}
}
