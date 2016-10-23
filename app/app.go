package app

import (
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/patrickmn/go-cache"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
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

	// Create a new cache instance with the given options
	// first parametter is the default item duration on the cache
	// second parametter is the tick time to purge all dead cache items
	util.Cache = cache.New(util.Config.Cache.Default.Duration, util.Config.Cache.Purge.Duration)

	// Create applicattion template
	util.Template = util.NewTemplate("castro")

	// Set template functions
	util.Template.FuncMap(templateFuncs())
	util.FuncMap = templateFuncs()

	// Load templates
	if err := util.LoadTemplates(&util.Template); err != nil {
		util.Logger.Fatalf("Cannot load templates: %v", err)
	}

	// Connect to the MySQL database
	if database.DB, err = database.Open(util.Config.Database.Username, util.Config.Database.Password, util.Config.Database.Name); err != nil {
		util.Logger.Fatalf("Cannot connect to MySQL database: %v", err)
	}

	// Migrate database models
	if err := database.DB.AutoMigrate(&models.Article{}).Error; err != nil {
		util.Logger.Fatalf("Cannot migrate database models: %v", err)
	}
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"isDev": func() bool {
			return util.Config.IsDev()
		},
		"url": func(args ...interface{}) string {
			url := ""
			for _, arg := range args {
				url = url + fmt.Sprintf("%v", arg)
			}
			return url
		},
	}
}
