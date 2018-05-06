package main

import (
	"fmt"
	"net/http"

	"crypto/tls"
	"encoding/gob"
	"log"
	"net/http/pprof"
	_ "net/http/pprof"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app"
	"github.com/raggaer/castro/app/controllers"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/store/memory"
	"github.com/urfave/negroni"
	"github.com/yuin/gopher-lua"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// Register gob data
	gob.Register(&models.CsrfToken{})
	gob.Register(&lua.LTable{})
	gob.Register(&util.CastroMap{})
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})

	// Show credits and application name
	fmt.Printf(`
Castro - High performance content management system for Open Tibia servers

Running version: %v
Compiled at: %v

`, util.VERSION, util.BUILD_DATE)

	// Check if application is installed
	if !isInstalled() {

		// Run the installation process
		if err := startInstallerApplication(); err != nil {
			log.Fatal(err)
		}

		return
	}

	// Run main app entry point
	app.Start()

	// Create rate-limiter instance
	rate := limiter.Rate{
		Period: util.Config.Configuration.RateLimit.Time.Duration,
		Limit:  util.Config.Configuration.RateLimit.Number,
	}

	// Create rate-limiter storage
	store := memory.NewStore()

	// Create rate-limiter
	limiter := limiter.New(store, rate)

	// Declare our new http router
	router := httprouter.New()

	// Declare application endpoints
	router.GET("/", controllers.LuaPage)
	router.POST("/", controllers.LuaPage)
	router.POST("/subtopic/*filepath", controllers.LuaPage)
	router.GET("/subtopic/*filepath", controllers.LuaPage)
	router.GET("/extensions/:id/static/*filepath", controllers.ExtensionStatic)
	router.POST("/nocsrf/*filepath", controllers.LuaPage)
	router.NotFound = PageNotFound

	// Register pprof router only on development mode
	if util.Config.Configuration.IsDev() {
		router.GET("/pprof/heap", wrapHandler(pprof.Handler("heap")))
	}

	// Create the session storage
	util.SessionStore = securecookie.New(
		[]byte(util.Config.Configuration.Cookies.HashKey),
		[]byte(util.Config.Configuration.Cookies.BlockKey),
	)

	// Create the middleware negroni instance with some application middleware
	n := negroni.New(
		newRateLimitHandler(limiter),
		newSecurityHandler(),
		newSessionHandler(),
		newMicrotimeHandler(),
		newCsrfHandler(),
		newI18nHandler(),
		negroni.NewStatic(http.Dir("public/")),
	)

	// Use negroni logger only in development mode
	if util.Config.Configuration.IsDev() || util.Config.Configuration.IsLog() {
		n.Use(negroni.NewLogger())
	}

	// Disable httprouter not found handler
	router.HandleMethodNotAllowed = false

	// Tell negroni to use our http router
	n.UseHandler(router)

	// Close database handle when the main function ends
	defer database.DB.Close()

	// Create castro server
	server := http.Server{
		Addr:         fmt.Sprintf(":%v", util.Config.Configuration.Port),
		Handler:      n,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Check if Castro should run on SSL mode
	if util.Config.Configuration.SSL.Enabled {

		// Check if user is using auto-certificate
		if util.Config.Configuration.SSL.Auto {

			// Create auto-certificate manager
			m := autocert.Manager{
				Prompt: autocert.AcceptTOS,
				Cache:  autocert.DirCache("tls"),
			}

			// Set auto-certificate hosts
			if strings.HasPrefix(util.Config.Configuration.URL, "www") {
				m.HostPolicy = autocert.HostWhitelist(util.Config.Configuration.URL, strings.Replace(util.Config.Configuration.URL, "www.", "", 1))
			} else {
				m.HostPolicy = autocert.HostWhitelist(util.Config.Configuration.URL, "www."+util.Config.Configuration.URL)
			}

			// Set server TLS option
			server.TLSConfig = &tls.Config{
				GetCertificate: m.GetCertificate,
			}

			// Listen to non-https ACME challenges connections
			go http.ListenAndServe(":http", m.HTTPHandler(nil))

			// Listen to https connections using autocert
			if err := server.ListenAndServeTLS("", ""); err != nil {
				util.Logger.Logger.Fatalf("Cannot start Castro autocert HTTPS server: %v", err)
			}
		}

		// Redirect all non https connections
		go httpsRedirect()

		// If SSL is enabled listen with cert and key
		if err := server.ListenAndServeTLS(
			util.Config.Configuration.SSL.Cert,
			util.Config.Configuration.SSL.Key,
		); err != nil {
			util.Logger.Logger.Fatalf("Cannot start Castro HTTPS server: %v", err)
		}

	} else {

		// Listen without using ssl
		if err := server.ListenAndServe(); err != nil {
			util.Logger.Logger.Fatalf("Cannot start Castro HTTP server: %v", err)
		}
	}
}

// wrapHandler converts a normal http handler to a httprouter handler
func wrapHandler(h http.Handler) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		h.ServeHTTP(rw, req)
	}
}

// httpsRedirect gets all non-https traffic and redirects to https
func httpsRedirect() {
	// Create router
	mux := httprouter.New()
	mux.GET("/*filepath", controllers.SSLRedirect)

	// Create server
	server := http.Server{
		Addr:         fmt.Sprintf(":%v", 80),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		util.Logger.Logger.Fatalf("Cannot start HTTP redirect server: %v", err)
	}
}
