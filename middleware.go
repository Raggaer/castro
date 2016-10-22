package main

import "net/http"

// cookieHandler used to make sure all requests
// contain a castro specific cookie
type cookieHandler struct {
	cookieDuration int
	cookieName     string
}

// newCookieHandler creates and returns a new cookieHandler
// instance with the given options
func newCookieHandler(duration int, name string) *cookieHandler {
	return &cookieHandler{
		cookieDuration: duration,
		cookieName:     name,
	}
}

// ServeHTTP makes cookieHandler compatible with negroni
func (c *cookieHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Check if the castro cookie is present
	_, err := req.Cookie(c.cookieName)

	if err != nil {
		// Cookie is not found so we create one
		newCookie := http.Cookie{
			Name:   c.cookieName,
			MaxAge: c.cookieDuration,
		}

		// Set the new cookie into the user
		http.SetCookie(w, &newCookie)
	}

	// Execute next handler
	next(w, req)
}
