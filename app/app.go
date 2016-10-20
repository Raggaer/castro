package app

import (
	"io/ioutil"

	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/util"
	"github.com/raggaer/castro/dialect"
	"github.com/raggaer/castro/dialect/tfs"
)

// Start the main execution point for Castro
func Start() {
	// Load the configration file
	file, err := ioutil.ReadFile("config.toml")
	if err != nil {
		util.Logger.Fatalf("Cannot read configuration file: %v", err)
	}
	if err = util.LoadConfig(string(file), util.Config); err != nil {
		util.Logger.Fatalf("Cannot read configuration file: %v", err)
	}

	// Load templates
	if err := util.LoadTemplates(&util.Template); err != nil {
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

	// Run dialect start function
	if err := dialect.Current.Start(); err != nil {
		util.Logger.Fatal(err)
	}

	// Load server stages
	if err := dialect.Current.LoadStages(); err != nil {
		util.Logger.Fatalf("Cannot load server stages: %v", err)
	}
}
