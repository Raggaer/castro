package app

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/util"
	"github.com/raggaer/castro/dialect"
	"github.com/raggaer/castro/dialect/tfs"
	"github.com/urfave/negroni"
)

// Start the main execution point for Castro
func Start() *httprouter.Router {
	// Load the configration file
	file, err := ioutil.ReadFile("config.toml")
	if err != nil {
		util.Logger.Fatalf("Cannot read configuration file: %v", err)
	}
	if err = util.LoadConfig(string(file), util.Config); err != nil {
		util.Logger.Fatalf("Cannot read configuration file: %v", err)
	}

	// Load templates
	if err := util.LoadTemplates(util.Template); err != nil {
		util.Logger.Fatalf("Cannot load templates: %v", err)
	}

	// Connect to the MySQL database
	if database.DB, err = database.Open(util.Config.Database.Username, util.Config.Database.Password, util.Config.Database.Name); err != nil {
		util.Logger.Fatalf("Cannot connect to MySQL database: %v", err)
	}
	defer database.DB.Close()

	// Load applicattion dialect
	dialect.SetDialect(&tfs.TFS{})
	util.Logger.Infof("Using dialect: %v - %v", dialect.Current.Name(), dialect.Current.Version())

	// Load server stages
	if err := dialect.Current.LoadStages(); err != nil {
		util.Logger.Fatalf("Cannot load server stages: %v", err)
	}

	// Create the http router instance
	mux := httprouter.New()

	// Create the middleware negroini instance with
	// some prefredined basics
	n := negroni.Classic()

	// Tell negroni to use our http router
	n.UseHandler(mux)

	util.Logger.Infof("Starting Castro http server on port :%v", util.Config.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", util.Config.Port), n); err != nil {
		// This should only happen when a port is
		// already in use
		util.Logger.Fatalf("Cannot start Castro http server: %v", err)
	}

	return mux
}
