package main

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

// microtimeHandler used to record all requests
// time spent
type microtimeHandler struct{}

// newMicrotimeHandler creates and returns a new microtimeHandler
// instance with the given format
func newMicrotimeHandler() *microtimeHandler {
	return &microtimeHandler{}
}

// ServeHTTP makes microtimeHandler compatible with negroni
func (m *microtimeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Set timestamp on the request context
	ctx := context.WithValue(req.Context(), "microtime", time.Now())

	// Execute next handler
	next(w, req.WithContext(ctx))
}
