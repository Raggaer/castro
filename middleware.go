package main

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
	"github.com/raggaer/castro/app/util"
)

type notFoundHandler struct {
}

// cookieHandler used to make sure all requests
// contain a castro specific cookie
type cookieHandler struct {}

// microtimeHandler used to record all requests
// time spent
type microtimeHandler struct {}

// csrfHandler used to grant CSRF tokens
// also checks POST requests
type csrfHandler struct {}

// newCsrfHandler creates and returns a new csrfHandler instance
func newCsrfHandler() *csrfHandler {
	return &csrfHandler{}
}

// newNotFoundHandler creates and returns a new notFoundHandler
// instance
func newNotFoundHandler() *notFoundHandler {
	return &notFoundHandler{}
}

// newCookieHandler creates and returns a new cookieHandler
// instance with the given options
func newCookieHandler() *cookieHandler {
	return &cookieHandler{}
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

// ServeHTTP makes notFoundHandler compatible with negroni
func (c *notFoundHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("Page was not found"))
}

// ServeHTTP makes cookieHandler compatible with negroni
func (c *cookieHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Check if the castro cookie is present
	_, err := req.Cookie(util.Config.Cookies.Name)

	if err != nil {

		// Create new unique token
		newToken, err := util.CreateUniqueToken(35)

		if err != nil {

			// Throw error with header 500
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))

			return
		}

		// Create JWT
		token, err := util.CreateJWToken(util.CastroClaims{
			CreatedAt: time.Now().Unix(),
			Token: newToken,
		})

		if err != nil {

			// Throw error with header 500
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))

			return
		}

		// Cookie is not found so we create one
		newCookie := http.Cookie{
			Name:   util.Config.Cookies.Name,
			MaxAge: util.Config.Cookies.MaxAge,
			Value: token,
			Secure: util.Config.SSL.Enabled,
			HttpOnly: true,
			Path: "/",
		}

		// Set the new cookie into the user
		http.SetCookie(w, &newCookie)
	}

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
