package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/util"
	"github.com/urfave/negroni"
)

// Start the main execution point for Castro
func Start() {
	// Create the http router instance
	mux := httprouter.New()

	// Create the middleware negroini instance with
	// some prefredined basics
	n := negroni.Classic()

	// Tell negroni to use our http router
	n.UseHandler(mux)

	// Everything set lets start the http server
	util.Logger.Infof("Starting Castro http server on port: %v", 8080)

	if err := http.ListenAndServe(":8080", n); err != nil {
		// This should only happen when a port is
		// already in use
		util.Logger.Errorf("Cannot start Castro http server: %v", err)
	}
}
