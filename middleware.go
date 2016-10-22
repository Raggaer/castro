package main

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

// cookieHandler used to make sure all requests
// contain a castro specific cookie
type cookieHandler struct {
	cookieDuration int
	cookieName     string
}

// microtimeHandler used to record all requests
// time spent
type microtimeHandler struct {
}

// newCookieHandler creates and returns a new cookieHandler
// instance with the given options
func newCookieHandler(duration int, name string) *cookieHandler {
	return &cookieHandler{
		cookieDuration: duration,
		cookieName:     name,
	}
}

// newMicrotimeHandler creates and returns a new microtimeHandler
// instance with the given format
func newMicrotimeHandler() *microtimeHandler {
	return &microtimeHandler{}
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

// ServeHTTP makes microtimeHandler compatible with negroni
func (m *microtimeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	ctx := context.WithValue(req.Context(), "microtime", time.Now())
	next(w, req.WithContext(ctx))
}
