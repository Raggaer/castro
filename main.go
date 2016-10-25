package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app"
	"github.com/raggaer/castro/app/controllers"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/util"
	"github.com/urfave/negroni"
)

func main() {
	// Show credits and applicattion name
	fmt.Println(`
			Castro - Open Tibia automatic account creator

			Developed by Raggaer
	`)

	// Declare our new http router
	router := httprouter.New()
	router.GET("/", controllers.Home)
	router.GET("/page/:page", controllers.Home)
	router.GET("/public/*filepath", serveStatic)

	// Check if Castro is installed
	if !isInstalled() {
		configFile, err := os.Create("config.toml")
		if err != nil {
			util.Logger.Fatalf("Cannot create config.toml file: %v", err)
		}
		defer configFile.Close()
		if err := toml.NewEncoder(configFile).Encode(util.Configuration{
			Mode:     "dev",
			Port:     8080,
			URL:      "http://localhost",
			Datapack: "/",
			Cookies: util.CookieConfig{
				Name:   "castro",
				MaxAge: 1000000,
			},
			Cache: util.CacheConfig{
				Default: int(time.Minute) * 5,
				Purge:   int(time.Minute),
			},
			Custom: make(map[string]interface{}),
		}); err != nil {
			util.Logger.Fatalf("Cannot encode config.toml file: %v", err)
		}
	}

	// Create the middleware negroni instance with
	// some middlewares
	n := negroni.New(
		newMicrotimeHandler(),
		negroni.NewRecovery(),
		newCookieHandler(
			10000,
			"castro",
		),
	)

	// Run main app entry point
	app.Start()

	// Use negroni logger only in development mode
	if util.Config.IsDev() {
		n.Use(negroni.NewLogger())
	}

	// Disable httprouter not found handler
	router.HandleMethodNotAllowed = false

	// Use custom not found handler
	router.NotFound = newNotFoundHandler()

	// Tell negroni to use our http router
	n.UseHandler(router)

	// Close database handle when the main function ends
	defer database.DB.Close()

	// Show running port
	util.Logger.Infof("Starting Castro http server on port :%v", util.Config.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", util.Config.Port), n); err != nil {
		// This should only happen when a port is
		// already in use
		util.Logger.Fatalf("Cannot start Castro http server: %v", err)
	}
}
