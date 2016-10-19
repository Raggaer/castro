package app

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/util"
	"github.com/urfave/negroni"
)

var (
	// Config variable to hold the main
	// configuration file
	Config = &util.Config{}

	// DB is the main handle for the database
	DB *gorm.DB
)

// Start the main execution point for Castro
func Start() {

	///////

	tfs := database.TFS{}
	database.SetDialect(tfs)
	util.Logger.Info(tfs.Name())

	///////

	// Load the configration file
	file, err := ioutil.ReadFile("config.toml")
	if err != nil {
		util.Logger.Fatalf("Cannot read configuration file: %v", err)
	}
	if err = util.LoadConfig(string(file), Config); err != nil {
		util.Logger.Fatalf("Cannot read configuration file: %v", err)
	}

	// Connect to the MySQL database
	if DB, err = database.Open(Config.Database.Username, Config.Database.Password, Config.Database.Name); err != nil {
		util.Logger.Fatalf("Cannot connect to MySQL database: %v", err)
	}
	defer DB.Close()

	// Create the http router instance
	mux := httprouter.New()

	// Create the middleware negroini instance with
	// some prefredined basics
	n := negroni.Classic()

	// Tell negroni to use our http router
	n.UseHandler(mux)

	util.Logger.Infof("Starting Castro http server on port :%v", Config.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", Config.Port), n); err != nil {
		// This should only happen when a port is
		// already in use
		util.Logger.Fatalf("Cannot start Castro http server: %v", err)
	}
}
