package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/raggaer/castro/app/database"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	"github.com/raggaer/otmap"
	glua "github.com/yuin/gopher-lua"
)

// Start the main execution point for Castro
func Start() {
	// Wait for all goroutines to make their work
	wait := &sync.WaitGroup{}

	// Wait for all tasks
	wait.Add(10)

	// Load application logger
	loadAppLogger()

	// Load application config
	loadAppConfig()

	// Run logger renew service
	go util.RenewLogger()

	loadLUAConfig()
	connectDatabase()

	// Execute our tasks
	go func(wait *sync.WaitGroup) {

		if util.Config.Configuration.LoadMap {
			loadMap()
			go mapWatcher()
		} else {
			loadConfigMap()
		}

		go loadHouses(wait)
		go loadVocations(wait)
		go loadServerMonsters(wait)
	}(wait)

	// Create application cache
	createCache()

	// Load local config files
	overwriteConfigFile()

	go loadLanguageFiles(wait)
	go loadWidgetList(wait)
	go appTemplates(wait)
	go widgetTemplates(wait)
	go loadSubtopics(wait)
	go loadWidgets(wait)

	loadExtensionStaticResources(wait)

	// Wait for the tasks
	wait.Wait()

	// Execute migrations
	executeMigrations()

	// Execute the init lua file
	executeInitFile()
}

func loadServerMonsters(wg *sync.WaitGroup) {
	// Load server monsters
	if err := util.LoadServerMonsters(util.Config.Configuration.Datapack); err != nil {
		util.Logger.Logger.Fatalf("Cannot load server monsters: %v", err)
	}

	// Sort server monsters list
	sort.Slice(util.MonstersList, func(i, j int) bool {
		return util.MonstersList[i].Name < util.MonstersList[j].Name
	})
	wg.Done()
}

func loadLanguageFiles(wg *sync.WaitGroup) {
	// Load language files
	if err := util.Loadi18n("i18n"); err != nil {
		util.Logger.Logger.Fatalf("Cannot load language files: %v", err)
	}

	// Finish task
	wg.Done()
}

func overwriteConfigFile() {
	// Overwrite config file
	if err := lua.OverwriteConfigFile(); err != nil {
		util.Logger.Logger.Fatalf("Cannot overwrite config file: %v", err)
	}
}

func loadExtensionStaticResources(wg *sync.WaitGroup) {
	// Load extension static resources
	if err := util.ExtensionStatic.Load("extensions"); err != nil {
		util.Logger.Logger.Fatalf("Cannot load extensions static resources: %v", err)
	}

	// Finish task
	wg.Done()
}

func mapWatcher() {
	// Check if watcher is enabled
	if !util.Config.Configuration.MapWatch.Enabled {
		return
	}

	// Create watcher ticker
	ticker := time.NewTicker(util.Config.Configuration.MapWatch.Check.Duration)
	defer ticker.Stop()

	// Start watcher loop
	for {
		select {
		case <-ticker.C:
			// Reload map
			loadMap()
		}
	}
}

func loadConfigMap() {
	mapTowns := []otmap.Town{}

	// Convert config towns to map towns
	for _, t := range util.Config.Configuration.Towns {
		mapTowns = append(mapTowns, otmap.Town{
			Name: t.Name,
			ID:   t.ID,
		})
	}

	// Set map global
	util.OTBMap.Load(&util.CastroMap{
		HouseFile: util.Config.Configuration.MapHouseFile,
		Towns:     mapTowns,
	})
}

func loadMap() {
	// Map holder
	m := models.Map{}

	// Get map mod time
	fileInformation, err := os.Stat(filepath.Join(util.Config.Configuration.Datapack, "data", "world", lua.Config.GetGlobal("mapName").String()+".otbm"))
	if err != nil {
		util.Logger.Logger.Fatalf("Cannot get map file information: %v", err)
	}

	// Check if map is encoded
	err = database.DB.Get(&m, "SELECT id, name, data, created_at, updated_at, last_modtime FROM castro_map WHERE name = ? ORDER BY id DESC", lua.Config.GetGlobal("mapName").String())

	// Unexpected error
	if err != nil && err != sql.ErrNoRows {
		util.Logger.Logger.Fatalf("Cannot retrieve map from database: %v", err)
	}

	// Check if map is not encoded
	if err == sql.ErrNoRows {

		fmt.Println(">> Encoding map. This process can take several minutes")
		util.Logger.Logger.Info("Encoding map. This process can take several minutes")

		// Encode map
		mapData, err := util.EncodeMap(
			filepath.Join(util.Config.Configuration.Datapack, "data", "world", lua.Config.GetGlobal("mapName").String()+".otbm"),
		)

		if err != nil {
			util.Logger.Logger.Fatalf("Cannot encode map file: %v", err)
		}

		// Update map struct
		m.Name = lua.Config.GetGlobal("mapName").String()
		m.Data = mapData
		m.Created_at = time.Now()
		m.Updated_at = time.Now()
		m.Last_modtime = fileInformation.ModTime()

		// Save map
		if _, err := database.DB.Exec(
			"INSERT INTO castro_map (name, data, created_at, updated_at, last_modtime) VALUES (?, ?, ?, ?, ?)",
			m.Name,
			m.Data,
			m.Created_at,
			m.Updated_at,
			m.Last_modtime,
		); err != nil {
			util.Logger.Logger.Fatalf("Cannot save encoded map file: %v", err)
		}
	}

	// Check if map is old
	if !fileInformation.ModTime().Round(time.Second).Equal(m.Last_modtime) {

		fmt.Println(">> Encoded map is outdated. Generating new map data")
		util.Logger.Logger.Info("Encoded map is outdated. Generating new map data")

		// Encode map
		mapData, err := util.EncodeMap(
			filepath.Join(util.Config.Configuration.Datapack, "data", "world", lua.Config.GetGlobal("mapName").String()+".otbm"),
		)

		if err != nil {
			util.Logger.Logger.Fatalf("Cannot encode map file: %v", err)
		}

		// Update map struct
		m.Name = lua.Config.GetGlobal("mapName").String()
		m.Data = mapData
		m.Created_at = time.Now()
		m.Updated_at = time.Now()
		m.Last_modtime = fileInformation.ModTime()

		// Save map
		if _, err := database.DB.Exec("UPDATE castro_map SET data = ?, created_at = ?, updated_at = ? WHERE name = ?", m.Data, m.Created_at, m.Updated_at, m.Name); err != nil {
			util.Logger.Logger.Fatalf("Cannot save encoded map file: %v", err)
		}

		// Log messages
		util.Logger.Logger.Info("New map data saved to database")
	}

	// Decode map
	castroMap, err := util.DecodeMap(m.Data)

	if err != nil {
		util.Logger.Logger.Fatalf("Cannot decode map file: %v", err)
	}

	// Set map global
	util.OTBMap.Load(castroMap)
}

func executeMigrations() {
	// Create migration state
	state := glua.NewState()

	// Set database metatable
	lua.SetDatabaseMetaTable(state)

	// Close state
	defer state.Close()

	// Walk migrations directory
	if err := filepath.Walk("migrations", func(path string, info os.FileInfo, err error) error {

		// Check if lua file
		if !strings.HasSuffix(path, ".lua") {
			return nil
		}

		// Do lua file
		if err := state.DoFile(path); err != nil {
			return err
		}

		// Call migration function
		if err := state.CallByParam(
			glua.P{
				Fn:      state.GetGlobal("migration"),
				NRet:    0,
				Protect: !util.Config.Configuration.IsDev(),
			},
		); err != nil {
			return err
		}

		// Pop state
		state.Pop(-1)

		return nil

	}); err != nil {
		util.Logger.Logger.Fatalf("Cannot run migration files: %v", err)
	}
}

func executeInitFile() {
	// Get lua state
	luaState := glua.NewState()

	// Close state
	defer luaState.Close()

	// Get application ready state
	lua.GetApplicationState(luaState)

	// Execute init file
	if err := lua.ExecuteFile(luaState, filepath.Join("engine", "init.lua")); err != nil {
		util.Logger.Logger.Fatalf("Cannot execute init lua file: %v", err)
	}
}

func loadWidgets(wg *sync.WaitGroup) {
	// Load subtopic list
	if err := lua.WidgetList.Load("widgets"); err != nil {
		util.Logger.Logger.Fatalf("Cannot load application widget list: %v", err)
	}

	// Load extension widgets
	if err := lua.WidgetList.LoadExtensions(); err != nil {
		util.Logger.Logger.Errorf("Cannot load extension widget list: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func loadSubtopics(wg *sync.WaitGroup) {
	// Compile pages files
	if err := lua.CompiledPageList.CompileFiles("pages"); err != nil {
		util.Logger.Logger.Fatalf("Cannot compile application subtopic list: %v", err)
	}

	// Compile extension pages files
	if err := lua.CompiledPageList.CompileExtensions("pages"); err != nil {
		util.Logger.Logger.Errorf("Cannot compile extension subtopic list: %v", err)
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
	util.Logger.LoggerOutput = f

	// Set last logger day
	util.Logger.LastLoggerDay = day

	// Create main application logger instance
	util.Logger.Logger = util.CreateLogger(f)
}

func loadVocations(wg *sync.WaitGroup) {
	// Load server vocations
	if err := util.LoadVocations(
		filepath.Join(util.Config.Configuration.Datapack, "data", "XML", "vocations.xml"),
		util.ServerVocationList,
	); err != nil {
		util.Logger.Logger.Fatalf("Cannot load map house list: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func loadHouses(wg *sync.WaitGroup) {
	// Load server houses
	if err := util.ServerHouseList.LoadHouses(
		filepath.Join(util.Config.Configuration.Datapack, "data", "world", util.OTBMap.Map.HouseFile),
	); err != nil {
		util.Logger.Logger.Fatalf("Cannot load map house list: %v", err)
	}

	util.Logger.Logger.Info("House list loaded")

	// Tell the wait group we are done
	wg.Done()
}

func loadAppConfig() {
	// Load the TOML configuration file
	if err := util.LoadConfig("config.toml"); err != nil {
		util.Logger.Logger.Fatalf("Cannot read configuration file: %v", err)
	}
}

func loadLUAConfig() {
	// Load the LUA configuration file
	if err := lua.LoadConfig(filepath.Join(util.Config.Configuration.Datapack, "config.lua")); err != nil {
		util.Logger.Logger.Fatalf("Cannot read lua configuration file: %v", err)
	}
}

func createCache() {
	// Create a new cache instance with the given options
	// first parameter is the default item duration on the cache
	// second parameter is the tick time to purge all dead cache items
	util.Cache = cache.New(
		util.Config.Configuration.Cache.Default.Duration,
		util.Config.Configuration.Cache.Purge.Duration,
	)
}

func loadWidgetList(wg *sync.WaitGroup) {
	// Load widget list
	if err := util.Widgets.Load("widgets/"); err != nil {
		util.Logger.Logger.Fatalf("Cannot load widget list: %v", err)
	}

	// Load extension widget list
	if err := util.Widgets.LoadExtensions(); err != nil {
		util.Logger.Logger.Errorf("Cannot load extension widget list: %v", err)
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
	if err := util.Template.LoadTemplates(util.Config.Configuration.Template); err != nil {
		util.Logger.Logger.Fatalf("Cannot load templates: %v", err)
	}

	// Load subtopic templates
	if err := util.Template.LoadTemplates("pages/"); err != nil {
		util.Logger.Logger.Error(err.Error())
		return
	}

	// Load extension subtopic templates
	if err := util.Template.LoadExtensionTemplates("pages"); err != nil {
		util.Logger.Logger.Errorf("Cannot load extension subtopic templates: %v", err)
	}

	// Load template hooks
	util.Template.LoadTemplateHooks()

	// Tell the wait group we are done
	wg.Done()
}

func widgetTemplates(wg *sync.WaitGroup) {
	// Create widget template
	util.WidgetTemplate = util.NewTemplate("widget")

	util.WidgetTemplate.FuncMap(templateFuncs())

	// Load widget templates
	if err := util.WidgetTemplate.LoadTemplates("widgets/"); err != nil {
		util.Logger.Logger.Fatalf("Cannot load widget templates: %v", err)
	}

	// Load extension widget templates
	if err := util.WidgetTemplate.LoadExtensionTemplates("widgets"); err != nil {
		util.Logger.Logger.Errorf("Cannot load extension widget templates: %v", err)
	}

	// Tell the wait group we are done
	wg.Done()
}

func connectDatabase() {
	var err error

	// Connect to the MySQL database
	if database.DB, err = database.Open(lua.Config.GetGlobal("mysqlUser").String(), 
		lua.Config.GetGlobal("mysqlPass").String(), 
		lua.Config.GetGlobal("mysqlHost").String(), 
		lua.Config.GetGlobal("mysqlPort").String(),
		lua.Config.GetGlobal("mysqlDatabase").String(),
		""); err != nil {
		util.Logger.Logger.Fatalf("Cannot connect to MySQL database: %v", err)
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
		"isMap": func(i interface{}) bool {
			return reflect.TypeOf(i).Kind() == reflect.Map
		},
		"isDev": func() bool {
			return util.Config.Configuration.IsDev()
		},
		"safeURL": func(s string) template.URL {
			return template.URL(s)
		},
		"url": func(args ...interface{}) template.URL {
			u := fmt.Sprintf("%v", util.Config.Configuration.URL)
			for _, arg := range args {
				u = u + fmt.Sprintf("/%v", arg)
			}
			if util.Config.Configuration.SSL.Proxy {
				return template.URL("https://" + u)
			}
			if util.Config.Configuration.SSL.Enabled {
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
		"str2html": func(text string) template.HTML {
			return template.HTML(text)
		},
		"str2attr": func(text string) template.HTMLAttr {
			return template.HTMLAttr(text)
		},
		"str2url": func(text string) template.URL {
			return template.URL(text)
		},
		"nl2br": func(text string) template.HTML {
			text = template.HTMLEscapeString(text)
			return template.HTML(
				strings.Replace(text, "\n", "<br>", -1),
			)
		},
		"urlEncode": func(t string) template.URL {
			return template.URL(url.QueryEscape(t))
		},
		"urlDecode": func(t string) string {
			v, err := url.QueryUnescape(t)
			if err != nil {
				return v
			}
			return ""
		},
		"serverName": func() string {
			return lua.Config.GetGlobal("serverName").String()
		},
		"serverMotd": func() string {
			return lua.Config.GetGlobal("motd").String()
		},
		"widgetList": func() []*util.Widget {
			return util.Widgets.List
		},
		"captchaKey": func() string {
			return util.Config.Configuration.Captcha.Public
		},
		"captchaEnabled": func() bool {
			return util.Config.Configuration.Captcha.Enabled
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"eqNumber": func(a, b float64) bool {
			return a == b
		},
		"gtNumber": func(a, b float64) bool {
			return a > b
		},
		"lsNumber": func(a, b float64) bool {
			return a < b
		},
		"menuPages": func() interface{} {
			return util.Config.GetCustomValue("MenuPages")
		},
		"widget": func(list map[string]template.HTML, index string) template.HTML {
			v, ok := list[index]
			if !ok {
				return ""
			}
			return v
		},
		"i18n": func(ilang interface{}, index string, args ...interface{}) string {
			// Convert to string
			lang := fmt.Sprintf("%v", ilang)

			// Load language
			language, ok := util.LanguageFiles.Get(lang)
			if ok {
				str, ok := language.Data[index]
				if ok {
					return fmt.Sprintf(str, args...)
				}
			}

			// Load default language string
			language, ok = util.LanguageFiles.Get("default")
			if ok {
				str, ok := language.Data[index]
				if ok {
					return fmt.Sprintf(str, args...)
				}
			}

			return ""
		},
		"atoi": func(s string) int {
			v, err := strconv.Atoi(s)
			if err != nil {
				return 0
			}
			return v
		},
		"formatFloat": func(i float64) string {
			return strconv.FormatFloat(i, 'f', -1, 64)
		},
		"itoa": func(i int) string {
			return strconv.Itoa(i)
		},
		"addInt": func(a, b int) int {
			return a + b
		},
		"toLower": func(s string) string {
			return strings.ToLower(s)
		},
		"toTitle": func(s string) string {
			return strings.ToTitle(s)
		},
	}
}
