package main

import (
	"fmt"
	"net/http"

	"encoding/gob"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/thirdparty/tollbooth_negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/julienschmidt/httprouter"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/raggaer/castro/app"
	"github.com/raggaer/castro/app/controllers"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	"github.com/urfave/negroni"
	"github.com/yuin/gopher-lua"
	"net/http/pprof"
	_ "net/http/pprof"
)

var (
	// VERSION is the application current version. Changed at compile time
	VERSION string

	// BUILD_DATE date of application compilation
	BUILD_DATE string
)

func main() {
	// Register gob data
	gob.Register(&models.CsrfToken{})
	gob.Register(&lua.LTable{})

	// Show credits and application name
	fmt.Printf(`
Castro - Open Tibia automatic account creator

Running version: %v
Compiled at: %v

`, VERSION, BUILD_DATE)

	// Declare our new http router
	router := httprouter.New()
	router.GET("/", controllers.LuaPage)
	router.POST("/", controllers.LuaPage)
	router.GET("/signature/:name", controllers.Signature)
	router.POST("/subtopic/*filepath", controllers.LuaPage)
	router.GET("/subtopic/*filepath", controllers.LuaPage)
	router.GET("/pprof/heap", wrapHandler(pprof.Handler("heap")))

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

	// Run main app entry point
	app.Start()

	// Create the middleware negroni instance with
	// some application middleware
	n := negroni.New(
		newMicrotimeHandler(),
		negroni.NewRecovery(),
		sessions.Sessions(
			util.Config.Cookies.Name,
			cookiestore.New([]byte(util.Config.Secret)),
		),
		newCsrfHandler(),
		gzip.Gzip(gzip.DefaultCompression),
		negroni.NewStatic(http.Dir("public/")),
		tollbooth_negroni.LimitHandler(
			tollbooth.NewLimiter(
				util.Config.RateLimit.Number,
				util.Config.RateLimit.Time,
			),
		),
	)

	// Use negroni logger only in development mode
	if util.Config.IsDev() || util.Config.IsLog() {
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

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		h.ServeHTTP(rw, req)
	}
}
