package main

import (
	"net/http"
	"time"

	"github.com/raggaer/castro/app/util"
	"golang.org/x/net/context"
)

// microtimeHandler used to record all requests
// time spent
type microtimeHandler struct{}

// csrfHandler used to grant CSRF tokens
// also checks POST requests
type csrfHandler struct{}

// newCsrfHandler creates and returns a new csrfHandler instance
func newCsrfHandler() *csrfHandler {
	return &csrfHandler{}
}

// newMicrotimeHandler creates and returns a new microtimeHandler
// instance with the given format
func newMicrotimeHandler() *microtimeHandler {
	return &microtimeHandler{}
}

// ServeHTTP makes csrfHandler compatible with negroni
func (c *csrfHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Get main cooke
	_, err := req.Cookie(util.Config.Cookies.Name)

	// If cant retrieve main cookie refresh
	// main cookie should always be available since
	// cookie handler is the first handler executed
	if err != nil {
		http.Redirect(w, req, req.RequestURI, 302)
		return
	}

	//c.Value

	// Execute next handler
	next(w, req)
}

// ServeHTTP makes microtimeHandler compatible with negroni
func (m *microtimeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Set timestamp on the request context
	ctx := context.WithValue(req.Context(), "microtime", time.Now())

	// Execute next handler
	next(w, req.WithContext(ctx))
}
