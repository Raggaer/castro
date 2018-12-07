package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"
)

// LuaPage executes the given lua page
func LuaPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Create application paypal REST client
	lua.CreatePaypalClient(util.Config.Configuration.PayPal.SandBox)

	// Check if request is POST
	if r.Method == http.MethodPost {

		// Parse POST form
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(500)
			return
		}
	}

	// If development mode reload pages, widgets and config file
	if util.Config.Configuration.IsDev() {

		// Reload config file
		if err := util.LoadConfig("config.toml"); err != nil {

			// Set error header
			w.WriteHeader(500)
			util.Logger.Logger.Errorf("Cannot reload config file: %v", err)

			return
		}

		// Reload config overwrites
		if err := lua.OverwriteConfigFile(); err != nil {

			// Set error header
			w.WriteHeader(500)
			util.Logger.Logger.Errorf("Cannot reload external config files: %v", err)

			return
		}

		// Reload pages
		if err := lua.CompiledPageList.CompileFiles("pages"); err != nil {

			// Set error header
			w.WriteHeader(500)
			util.Logger.Logger.Errorf("Cannot reload subtopic %v: %v", ps.ByName("page"), err)

			return
		}

		// Reload extension pages
		if err := lua.PageList.LoadExtensions(); err != nil {

			// If AAC is running on development mode log error
			if util.Config.Configuration.IsDev() || util.Config.Configuration.IsLog() {
				util.Logger.Logger.Errorf("Cannot load extension subtopic %v: %v", ps.ByName("page"), err)
			}
		}

		// Reload extension static list
		if err := util.ExtensionStatic.Load("extensions"); err != nil {

			// If AAC is running on development mode log error
			if util.Config.Configuration.IsDev() || util.Config.Configuration.IsLog() {
				util.Logger.Logger.Errorf("Cannot load extension subtopic %v: %v", ps.ByName("page"), err)
			}
		}

		// Reload widgets
		if err := lua.WidgetList.Load("widgets"); err != nil {

			// Set error header
			w.WriteHeader(500)
			util.Logger.Logger.Errorf("Cannot reload widgets when executing %v subtopic: %v", ps.ByName("page"), err)

			return
		}

		// Reload extension widgets
		if err := lua.WidgetList.LoadExtensions(); err != nil {
			util.Logger.Logger.Errorf("Cannot load extension widgets when executing %v subtopic: %v", ps.ByName("page"), err)
		}

		// Reload widget list
		if err := util.Widgets.Load("widgets/"); err != nil {
			util.Logger.Logger.Fatalf("Cannot load widget list: %v", err)
		}

		// Reload extension widget list
		if err := util.Widgets.LoadExtensions(); err != nil {
			util.Logger.Logger.Errorf("Cannot load extension widget list: %v", err)
		}
	}

	// Get session
	session, ok := r.Context().Value("session").(map[string]interface{})

	if !ok {
		// Set error header
		w.WriteHeader(500)
		util.Logger.Logger.Error("Cannot get session as map")

		return
	}

	// Get request language
	language, ok := r.Context().Value("language").([]string)
	if !ok {
		// Set error header
		w.WriteHeader(500)
		util.Logger.Logger.Error("Cannot get language as string slice")

		return
	}

	// Set LUA file name
	pageName := ps.ByName("filepath")

	// If there is no subtopic request index
	if pageName == "" {
		pageName = "index"
	}

	// Get state from the pool
	s := lua.NewState()

	// Create HTTP metatable
	lua.SetHTTPMetaTable(s)

	// Set the state user data
	lua.SetHTTPUserData(s, w, r)

	// Set session user data
	lua.SetSessionMetaTableUserData(s, session)

	// Set language user data
	lua.SetI18nUserData(s, language)

	// Retrieve compiled proto
	proto, err := lua.CompiledPageList.Get(filepath.Join("pages", pageName, r.Method+".lua"))
	if err != nil {
		w.WriteHeader(404)
		util.Logger.Logger.Errorf("Cannot find lua proto, subtopic source: %v", pageName, err)
		return
	}

	// Execute compiled file
	if err := lua.DoCompiledFile(
		s,
		proto,
	); err != nil {
		w.WriteHeader(404)
		util.Logger.Logger.Errorf("Cannot get %v subtopic source: %v", pageName, err)
		return
	}

	if err := lua.ExecuteControllerPage(s, r.Method); err != nil {
		w.WriteHeader(500)
		util.Logger.Logger.Errorf("Cannot execute subtopic %v: %v", pageName, err)
	}
}
