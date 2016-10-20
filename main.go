package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app"
	"github.com/raggaer/castro/app/controllers"
	"github.com/raggaer/castro/app/util"
	"github.com/urfave/negroni"
)

func main() {
	// Declare our new http router
	router := httprouter.New()
	router.GET("/", controllers.Home)
	router.GET("/public/*filepath", serveStatic)

	// Create the middleware negroini instance with
	// some prefredined basics
	n := negroni.New(
		negroni.NewLogger(),
		negroni.NewRecovery(),
	)

	// Tell negroni to use our http router
	n.UseHandler(router)

	// Run main app entry point
	app.Start()

	util.Logger.Infof("Starting Castro http server on port :%v", util.Config.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", util.Config.Port), n); err != nil {
		// This should only happen when a port is
		// already in use
		util.Logger.Fatalf("Cannot start Castro http server: %v", err)
	}
}
