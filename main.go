package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app"
	"github.com/raggaer/castro/app/controllers"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/util"
	"github.com/urfave/negroni"
)

func main() {
	// Show credits and application name
	fmt.Println(`
			Castro - Open Tibia automatic account creator

			Developed by Raggaer
	`)

	// Declare our new http router
	router := httprouter.New()
	router.GET("/", controllers.LuaPage)
	router.POST("/", controllers.LuaPage)
	router.GET("/signature/:name", controllers.Signature)
	router.POST("/subtopic/:page", controllers.LuaPage)
	router.GET("/subtopic/:page", controllers.LuaPage)
	router.GET("/public/*filepath", serveStatic)

	// Check if Castro is installed if not we create the
	// config file on runtime
	if !isInstalled() {

		// Create the config file with the given name
		if err := createConfigFile("config.toml"); err != nil {
			util.Logger.Fatalf("Cannot create %v file: %v", "config.toml", err)
		}

		util.Logger.Info("Config.toml file is now installed. Edit the configuration file and start Castro")

		// Exit app
		return
	}

	// Create the middleware negroni instance with
	// some middlewares
	n := negroni.New(
		newCookieHandler(),
		newMicrotimeHandler(),
		negroni.NewRecovery(),
		newCsrfHandler(),
	)

	// Run main app entry point
	app.Start()

	// Use negroni logger only in development mode
	if util.Config.IsDev() {
		n.Use(negroni.NewLogger())
	}

	// Disable httprouter not found handler
	router.HandleMethodNotAllowed = false

	// Tell negroni to use our http router
	n.UseHandler(router)

	// Close database handle when the main function ends
	defer database.DB.Close()

	// Show running port
	util.Logger.Infof("Starting Castro http server on port :%v", util.Config.Port)

	// Check if Castro should run on SSL mode
	if util.Config.SSL.Enabled {

		// If SSL is enabled listen with cert and key
		if err := http.ListenAndServeTLS(
			fmt.Sprintf("%v:%v", util.Config.URL, util.Config.Port),
			util.Config.SSL.Cert,
			util.Config.SSL.Key,
			n,
		); err != nil {
			util.Logger.Fatalf("Cannot start Castro HTTP server: %v", err)
		}
	} else {
		if err := http.ListenAndServe(fmt.Sprintf("%v:%v", util.Config.URL, util.Config.Port), n); err != nil {
			// This should only happen when a port is
			// already in use
			util.Logger.Fatalf("Cannot start Castro HTTP server: %v", err)
		}
	}
}
