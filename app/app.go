package app

import (
	"fmt"
	"html/template"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"
	"github.com/raggaer/otmap"
	"log"
	"net/url"
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
}

func loadWidgets(wg *sync.WaitGroup) {
	// Load subtopic list
	if err := lua.Widgets.Load("widgets"); err != nil {
		util.Logger.Fatalf("Cannot load application widget list: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func loadSubtopics(wg *sync.WaitGroup) {
	// Load subtopic list
	if err := lua.Subtopics.Load("pages"); err != nil {
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

func loadMap() {
	// Parse OTBM file
	m, err := otmap.Parse(util.Config.Datapack + "/data/world/" + lua.Config.MapName + ".otbm")

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
				strings.Replace(text, "\n", "<br>", -1),
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
	}
}
