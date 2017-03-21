package controllers

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// SSLRedirect redirects all request to HTTPS
func SSLRedirect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Create target URL
	target := "https://" + r.Host + r.URL.Path

	// Add values
	if len(r.URL.RawQuery) > 0 {
		target += "?" + r.URL.RawQuery
	}

	// Redirect user
	http.Redirect(w, r, target, 302)
}
