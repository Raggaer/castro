package main

import (
	"net/http"
	"time"

	"github.com/dchest/uniuri"
	"github.com/goincremental/negroni-sessions"
	"github.com/raggaer/castro/app/models"
	"golang.org/x/net/context"
	"strings"
)

// microtimeHandler used to record all requests
// time spent
type microtimeHandler struct{}

// csrfHandler used to add a token to all
// requests
type csrfHandler struct{}

// newCsrfHandler creates and returns a new csrfHandler
// instance
func newCsrfHandler() *csrfHandler {
	return &csrfHandler{}
}

func (c *csrfHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if !strings.HasPrefix(req.RequestURI, "/subtopic/") {
		next(w, req)
		return
	}

	// Get request session
	session := sessions.GetSession(req)

	// Check if token is set
	t := session.Get("csrf-token")

	// Token does not exists
	if t == nil {

		// If request is POST stop execution
		if req.Method == http.MethodPost {
			return
		}

		// Set token
		tkn := models.CsrfToken{
			Token: uniuri.New(),
			At:    time.Now(),
		}

		// Set session value
		session.Set("csrf-token", tkn)

		// Set token on the request context
		ctx := context.WithValue(req.Context(), "csrf-token", tkn)

		// Execute next handler
		next(w, req.WithContext(ctx))
		return
	}

	// Get token
	tkn, ok := t.(*models.CsrfToken)

	if !ok {
		return
	}

	// Check if token is valid
	if req.Method == http.MethodPost && tkn.Token != req.PostFormValue("_csrf") {
		return
	}

	// Check if token needs to be refreshed
	if time.Now().Before(tkn.At) {
		tkn.Token = uniuri.New()
		tkn.At = time.Now()
	}

	// Set session value
	session.Set("csrf-token", tkn)

	// Set token on the request context
	ctx := context.WithValue(req.Context(), "csrf-token", tkn)

	// Execute next handler
	next(w, req.WithContext(ctx))
}

// newMicrotimeHandler creates and returns a new microtimeHandler
// instance
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
