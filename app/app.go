package app

import (
	"fmt"
	"html/template"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	"github.com/raggaer/otmap"
	"strconv"
	"strings"
	"sync"
)

// Start the main execution point for Castro
func Start() {
	// Wait for all goroutines to make their work
	wait := &sync.WaitGroup{}

	// Wait for all tasks
	wait.Add(11)

	// Load application config
	loadAppConfig(wait)

	// Execute our tasks
	go func(wait *sync.WaitGroup) {

		loadLUAConfig(wait)
		connectDatabase(wait)
		migrateDatabase(wait)
		loadMap(wait)
		go loadHouses(wait)
		go loadVocations(wait)
	}(wait)

	// Create application cache
	createCache(wait)

	go loadWidgetList(wait)
	go appTemplates(wait)
	go widgetTemplates(wait)

	// Wait for the tasks
	wait.Wait()
}

func loadVocations(wg *sync.WaitGroup) {
	// Load server vocations
	if err := util.LoadVocations(
		util.Config.Datapack+"/data/xml/vocations.xml",
		util.ServerVocationList,
	); err != nil {
		util.Logger.Fatalf("Cannot load map house list: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func loadHouses(wg *sync.WaitGroup) {
	// Load server houses
	if err := util.LoadHouses(
		util.Config.Datapack+"/data/world/"+util.OTBMap.HouseFile,
		util.ServerHouseList,
	); err != nil {
		util.Logger.Fatalf("Cannot load map house list: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func loadMap(wg *sync.WaitGroup) {
	// Parse OTBM file
	m, err := otmap.Parse(util.Config.Datapack + "/data/world/" + lua.Config.MapName + ".otbm")

	if err != nil {
		util.Logger.Fatalf("Cannot parse OTBM file: %v", err)
	}

	util.OTBMap = m

	// Tell the wait group we are done
	wg.Done()
}

func loadAppConfig(wg *sync.WaitGroup) {
	// Load the TOML configuration file
	if err := util.LoadConfig("config.toml", util.Config); err != nil {
		util.Logger.Fatalf("Cannot read configuration file: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func loadLUAConfig(wg *sync.WaitGroup) {
	// Load the LUA configuration file
	if err := lua.LoadConfig(util.Config.Datapack, lua.Config); err != nil {
		util.Logger.Fatalf("Cannot read lua configuration file: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func createCache(wg *sync.WaitGroup) {
	// Create a new cache instance with the given options
	// first parameter is the default item duration on the cache
	// second parameter is the tick time to purge all dead cache items
	util.Cache = cache.New(time.Duration(util.Config.Cache.Default), time.Duration(util.Config.Cache.Purge))

	// Tell the wait group we are done
	wg.Done()
}

func loadWidgetList(wg *sync.WaitGroup) {
	// Load widget list
	if err := util.Widgets.LoadWidgetList("widgets/"); err != nil {
		util.Logger.Fatalf("Cannot load widget list: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func appTemplates(wg *sync.WaitGroup) {
	// Create application template
	util.Template = util.NewTemplate("castro")

	// Set template functions
	util.Template.FuncMap(templateFuncs())
	util.FuncMap = templateFuncs()

	// Load templates
	if err := util.LoadTemplates("views/", &util.Template); err != nil {
		util.Logger.Fatalf("Cannot load templates: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func widgetTemplates(wg *sync.WaitGroup) {
	// Create widget template
	util.WidgetTemplate = util.NewTemplate("widget")

	util.WidgetTemplate.FuncMap(templateFuncs())

	// Load widget templates
	if err := util.LoadTemplates("widgets/", &util.WidgetTemplate); err != nil {
		util.Logger.Fatalf("Cannot load widget templates: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func connectDatabase(wg *sync.WaitGroup) {
	var err error

	// Connect to the MySQL database
	if database.DB, err = database.Open(lua.Config.MySQLUser, lua.Config.MySQLPass, lua.Config.MySQLDatabase); err != nil {
		util.Logger.Fatalf("Cannot connect to MySQL database: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func migrateDatabase(wg *sync.WaitGroup) {
	// Migrate database models
	if err := database.DB.AutoMigrate(&models.Article{}, &models.CastroAccount{}).Error; err != nil {
		util.Logger.Fatalf("Cannot migrate database models: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"isDev": func() bool {
			return util.Config.IsDev()
		},
		"url": func(args ...interface{}) template.URL {
			url := fmt.Sprintf("%v:%v", util.Config.URL, util.Config.Port)
			for _, arg := range args {
				url = url + fmt.Sprintf("/%v", arg)
			}
			if util.Config.SSL.Enabled {
				return template.URL("https://" + url)
			}
			return template.URL("http://" + url)
		},
		"queryResults": func(m map[string]interface{}) []interface{} {
			n := len(m)
			r := []interface{}{}
			for i := 0; i < n; i++ {
				r = append(r, m[strconv.Itoa(i+1)])
			}
			return r
		},
		"unixToDate": func(m int64) template.HTML {
			date := time.Unix(m, 0)
			return template.HTML(
				date.Format("2006 - Mon Jan 2 15:04:05"),
			)
		},
		"nl2br": func(text string) template.HTML {
			return template.HTML(
				strings.Replace(text, "\n", "<br>", -1),
			)
		},
		"serverName": func() string {
			return lua.Config.ServerName
		},
		"serverMotd": func() string {
			return lua.Config.Motd
		},
		"widgetList": func() []*util.Widget {
			return util.Widgets.List
		},
		"captchaKey": func() string {
			return util.Config.Captcha.Public
		},
		"captchaEnabled": func() bool {
			return util.Config.Captcha.Enabled
		},
	}
}
