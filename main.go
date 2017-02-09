package main

import (
	"fmt"
	"net/http"

	"encoding/gob"
	"github.com/gorilla/securecookie"
	"github.com/julienschmidt/httprouter"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/raggaer/castro/app"
	"github.com/raggaer/castro/app/controllers"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	"github.com/urfave/negroni"
	"github.com/yuin/gopher-lua"
	"log"
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

	// Check if application is installed
	if !isInstalled() {

		fmt.Println("Castro is not installed. Running installation process")

		// Run the installation process
		if err := installApplication(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Configuration file created (%v) . Installation process is now done", configFileName)

		return
	}

	// Declare application endpoints
	router.GET("/", controllers.LuaPage)
	router.POST("/", controllers.LuaPage)
	router.GET("/signature/:name", controllers.Signature)
	router.POST("/subtopic/*filepath", controllers.LuaPage)
	router.GET("/subtopic/*filepath", controllers.LuaPage)
	router.GET("/pprof/heap", wrapHandler(pprof.Handler("heap")))

	// Run main app entry point
	app.Start()

	// Create the session storage
	util.SessionStore = securecookie.New(
		[]byte(util.Config.Cookies.HashKey),
		[]byte(util.Config.Cookies.BlockKey),
	)

	// Create the middleware negroni instance with
	// some application middleware
	n := negroni.New(
		newSessionHandler(),
		newMicrotimeHandler(),
		newCsrfHandler(),
		gzip.Gzip(gzip.DefaultCompression),
		negroni.NewStatic(http.Dir("public/")),
	)

	// Use negroni logger only in development mode
	if util.Config.IsDev() || util.Config.IsLog() {
		n.Use(negroni.NewLogger())

	} else {

		// Use rate-limiter on production mode
		/*n.Use(tollbooth_negroni.LimitHandler(
			tollbooth.NewLimiter(
				util.Config.RateLimit.Number,
				util.Config.RateLimit.Time,
			),
		))*/
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

		// Listen without using ssl
		if err := http.ListenAndServe(fmt.Sprintf("%v:%v", util.Config.URL, util.Config.Port), n); err != nil {
			// This should only happen when a port is
			// already in use
			util.Logger.Fatalf("Cannot start Castro HTTP server: %v", err)
		}
	}
}

// wrapHandler converts a normal http handler to a httprouter handler
func wrapHandler(h http.Handler) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		h.ServeHTTP(rw, req)
	}
}
