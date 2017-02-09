package main

import (
	"net/http"
	"time"

	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	"golang.org/x/net/context"
)

// microtimeHandler used to record all requests time spent
type microtimeHandler struct{}

// csrfHandler used to add a token to all requests
type csrfHandler struct{}

// sessionHandler used for application session
type sessionHandler struct{}

// newSessionHandler creates and returns a new sessionHandler instance
func newSessionHandler() *sessionHandler {
	return &sessionHandler{}
}

func (s *sessionHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Get application cookie
	cookie, err := req.Cookie(util.Config.Cookies.Name)

	if err != nil {

		// Cookie data
		v := make(map[string]interface{})

		// Set issuer field
		v["issuer"] = "Castro"

		// Encode cookie value
		encoded, err := util.SessionStore.Encode(util.Config.Cookies.Name, v)

		if err != nil {
			util.Logger.Fatalf("Cannot encode cookie value: %v", err)
			return
		}

		// Create cookie
		c := &http.Cookie{
			Name:  util.Config.Cookies.Name,
			Value: encoded,
			Path:  "/",
		}

		// Set cookie
		http.SetCookie(w, c)

		// Create new context with cookie value
		ctx := context.WithValue(req.Context(), "session", v)

		// Run next handler
		next(w, req.WithContext(ctx))

		return
	}

	// Cookie data holder
	v := make(map[string]interface{})

	// Decode cookie
	if err := util.SessionStore.Decode(
		util.Config.Cookies.Name,
		cookie.Value,
		&v,
	); err != nil {

		util.Logger.Fatalf("Cannot decode cookie value: %v", err)
		return
	}

	// Check issuer
	if v["issuer"].(string) != "Castro" {
		return
	}

	// Create new context with cookie value
	ctx := context.WithValue(req.Context(), "session", v)

	// Run next handler
	next(w, req.WithContext(ctx))
}

// newCsrfHandler creates and returns a new csrfHandler instance
func newCsrfHandler() *csrfHandler {
	return &csrfHandler{}
}

func (c *csrfHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Get session
	_ = req.Context().Value("session").(map[string]interface{})

	// Create new context
	ctx := context.WithValue(req.Context(), "csrf-token", &models.CsrfToken{
		Token: "12",
	})

	// Run next handler
	next(w, req.WithContext(ctx))
}

// newMicrotimeHandler creates and returns a new microtimeHandler instance
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
