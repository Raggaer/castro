package app

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	"strconv"
)

// Start the main execution point for Castro
func Start() {
	// Load the TOML configuration file
	file, err := ioutil.ReadFile("config.toml")
	if err != nil {
		util.Logger.Fatalf("Cannot read configuration file: %v", err)
	}
	if err = util.LoadConfig(string(file), util.Config); err != nil {
		util.Logger.Fatalf("Cannot read configuration file: %v", err)
	}

	// Load the LUA configuration file
	if err := lua.LoadConfig(util.Config.Datapack, lua.Config); err != nil {
		util.Logger.Fatalf("Cannot read lua configuration file: %v", err)
	}

	// Create a new cache instance with the given options
	// first parameter is the default item duration on the cache
	// second parameter is the tick time to purge all dead cache items
	util.Cache = cache.New(time.Duration(util.Config.Cache.Default), time.Duration(util.Config.Cache.Purge))

	// Create application template
	util.Template = util.NewTemplate("castro")

	// Set template functions
	util.Template.FuncMap(templateFuncs())
	util.FuncMap = templateFuncs()

	// Load templates
	if err := util.LoadTemplates(&util.Template); err != nil {
		util.Logger.Fatalf("Cannot load templates: %v", err)
	}

	// Connect to the MySQL database
	if database.DB, err = database.Open(lua.Config.MySQLUser, lua.Config.MySQLPass, lua.Config.MySQLDatabase); err != nil {
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
		"url": func(args ...interface{}) template.URL {
			url := fmt.Sprintf("%v:%v", util.Config.URL, util.Config.Port)
			for _, arg := range args {
				url = url + fmt.Sprintf("/%v", arg)
			}
			return template.URL(url)
		},
		"queryResults": func(m map[interface{}]interface{}) []interface{} {
			n := len(m)
			r := []interface{}{}
			for i := 0; i < n; i++ {
				r = append(r, m[strconv.Itoa(i+1)])
			}
			return r
		},
		"unixToDate": func(m int64) template.HTML {
			date := time.Unix(0, m)
			return template.HTML(
				date.Format("2006 - Mon Jan 2 15:04:05"),
			)
		},
	}
}
