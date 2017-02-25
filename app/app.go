package app

import (
	"fmt"
	glua "github.com/yuin/gopher-lua"
	"html/template"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"
	"github.com/raggaer/otmap"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Start the main execution point for Castro
func Start() {
	// Wait for all goroutines to make their work
	wait := &sync.WaitGroup{}

	// Wait for all tasks
	wait.Add(7)

	// Load application logger
	loadAppLogger()

	// Load application config
	loadAppConfig()

	// Run logger renew service
	go util.RenewLogger()

	// Execute our tasks
	go func(wait *sync.WaitGroup) {

		loadLUAConfig()
		connectDatabase()
		loadMap()
		go loadHouses(wait)
		go loadVocations(wait)
	}(wait)

	// Create application cache
	createCache()

	go loadWidgetList(wait)
	go appTemplates(wait)
	go widgetTemplates(wait)
	go loadSubtopics(wait)
	go loadWidgets(wait)

	// Wait for the tasks
	wait.Wait()

	// Execute the init lua file
	executeInitFile()
}

func executeInitFile() {
	// Get lua state
	luaState := glua.NewState()

	// Close state
	defer luaState.Close()

	// Create events metatable
	lua.SetEventsMetaTable(luaState)

	// Create storage metatable
	lua.SetStorageMetaTable(luaState)

	// Create time metatable
	lua.SetTimeMetaTable(luaState)

	// Create url metatable
	lua.SetURLMetaTable(luaState)

	// Create debug metatable
	lua.SetDebugMetaTable(luaState)

	// Create XML metatable
	lua.SetXMLMetaTable(luaState)

	// Create captcha metatable
	lua.SetCaptchaMetaTable(luaState)

	// Create crypto metatable
	lua.SetCryptoMetaTable(luaState)

	// Create validator metatable
	lua.SetValidatorMetaTable(luaState)

	// Create database metatable
	lua.SetDatabaseMetaTable(luaState)

	// Create config metatable
	lua.SetConfigMetaTable(luaState)

	// Create map metatable
	lua.SetMapMetaTable(luaState)

	// Create mail metatable
	lua.SetMailMetaTable(luaState)

	// Create cache metatable
	lua.SetCacheMetaTable(luaState)

	// Create reflect metatable
	lua.SetReflectMetaTable(luaState)

	// Create json metatable
	lua.SetJSONMetaTable(luaState)

	// Execute init file
	if err := luaState.DoFile(filepath.Join("engine", "init.lua")); err != nil {
		util.Logger.Fatalf("Cannot execute init lua file: %v", err)
	}
}

func loadWidgets(wg *sync.WaitGroup) {
	// Load subtopic list
	if err := lua.WidgetList.Load("widgets"); err != nil {
		util.Logger.Fatalf("Cannot load application widget list: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func loadSubtopics(wg *sync.WaitGroup) {
	// Load subtopic list
	if err := lua.PageList.Load("pages"); err != nil {
		util.Logger.Fatalf("Cannot load application subtopic list: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func loadAppLogger() {
	// Create logger file
	f, day, err := util.CreateLogFile()

	if err != nil {
		log.Fatal(err)
	}

	// Set logger output variable
	util.LoggerOutput = f

	// Set last logger day
	util.LastLoggerDay = day

	// Create main application logger instance
	util.Logger = util.CreateLogger(f)
}

func loadVocations(wg *sync.WaitGroup) {
	// Load server vocations
	if err := util.LoadVocations(
		filepath.Join(util.Config.Datapack, "data", "xml", "vocations.xml"),
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
		filepath.Join(util.Config.Datapack, "data", "world", util.OTBMap.HouseFile),
		util.ServerHouseList,
	); err != nil {
		util.Logger.Fatalf("Cannot load map house list: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func loadMap() {
	// Parse OTBM file
	m, err := otmap.Parse(filepath.Join(util.Config.Datapack, "data", "world", lua.Config.MapName+".otbm"))

	if err != nil {
		util.Logger.Fatalf("Cannot parse OTBM file: %v", err)
	}

	util.OTBMap = m
}

func loadAppConfig() {
	// Load the TOML configuration file
	if err := util.LoadConfig("config.toml", util.Config); err != nil {
		util.Logger.Fatalf("Cannot read configuration file: %v", err)
	}
}

func loadLUAConfig() {
	// Load the LUA configuration file
	if err := lua.LoadConfig(util.Config.Datapack, lua.Config); err != nil {
		util.Logger.Fatalf("Cannot read lua configuration file: %v", err)
	}
}

func createCache() {
	// Create a new cache instance with the given options
	// first parameter is the default item duration on the cache
	// second parameter is the tick time to purge all dead cache items
	util.Cache = cache.New(util.Config.Cache.Default, util.Config.Cache.Purge)
}

func loadWidgetList(wg *sync.WaitGroup) {
	// Load widget list
	if err := util.Widgets.Load("widgets/"); err != nil {
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
	if err := util.Template.LoadTemplates("views/"); err != nil {
		util.Logger.Fatalf("Cannot load templates: %v", err)
	}

	// Load subtopic templates
	if err := util.Template.LoadTemplates("pages/"); err != nil {
		util.Logger.Error(err.Error())
		return
	}

	// Tell the wait group we are done
	wg.Done()
}

func widgetTemplates(wg *sync.WaitGroup) {
	// Create widget template
	util.WidgetTemplate = util.NewTemplate("widget")

	util.WidgetTemplate.FuncMap(templateFuncs())

	// Load widget templates
	if err := util.WidgetTemplate.LoadTemplates("widgets/"); err != nil {
		util.Logger.Fatalf("Cannot load widget templates: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func connectDatabase() {
	var err error

	// Connect to the MySQL database
	if database.DB, err = database.Open(lua.Config.MySQLUser, lua.Config.MySQLPass, lua.Config.MySQLDatabase); err != nil {
		util.Logger.Fatalf("Cannot connect to MySQL database: %v", err)
	}
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"vocation": func(voc float64) string {
			for _, v := range util.ServerVocationList.List.Vocations {
				if v.ID == int(voc) {
					return v.Name
				}
			}
			return ""
		},
		"isDev": func() bool {
			return util.Config.IsDev()
		},
		"url": func(args ...interface{}) template.URL {
			u := fmt.Sprintf("%v:%v", util.Config.URL, util.Config.Port)
			for _, arg := range args {
				u = u + fmt.Sprintf("/%v", arg)
			}
			if util.Config.SSL.Enabled {
				return template.URL("https://" + u)
			}
			return template.URL("http://" + u)
		},
		"queryResults": func(m map[string]interface{}) []interface{} {
			n := len(m)
			r := []interface{}{}
			for i := 0; i < n; i++ {
				r = append(r, m[strconv.Itoa(i+1)])
			}
			return r
		},
		"unixToDate": func(m float64) template.HTML {
			date := time.Unix(int64(m), 0)
			return template.HTML(
				date.Format("2006 - Mon Jan 2 15:04:05"),
			)
		},
		"nl2br": func(text string) template.HTML {
			return template.HTML(
				strings.Replace(text, "\r\n", "<br>", -1),
			)
		},
		"urlEncode": func(t string) template.URL {
			return template.URL(url.QueryEscape(t))
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
		"eqNumber": func(a, b float64) bool {
			return a == b
		},
	}
}
