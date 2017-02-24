package controllers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"
	glua "github.com/yuin/gopher-lua"
	"net/http"
	"path/filepath"
	"strings"
)

// LuaPage executes the given lua page
func LuaPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check if request is POST
	if r.Method == http.MethodPost {

		// Parse POST form
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(500)
			return
		}
	}

	// If development mode reload pages and widgets
	if util.Config.IsDev() {

		// Reload pages
		if err := lua.PageList.Load("pages"); err != nil {
			// Set error header
			w.WriteHeader(500)

			// If AAC is running on development mode log error
			if util.Config.IsDev() || util.Config.IsLog() {
				util.Logger.Errorf("Cannot load subtopic %v: %v", ps.ByName("page"), err)
			}

			return
		}

		// Reload widgets
		if err := lua.WidgetList.Load("widgets"); err != nil {
			// Set error header
			w.WriteHeader(500)

			// If AAC is running on development mode log error
			if util.Config.IsDev() || util.Config.IsLog() {
				util.Logger.Errorf("Cannot load widgets when executing %v subtopic: %v", ps.ByName("page"), err)
			}

			return
		}
	}

	// Get session
	session, ok := r.Context().Value("session").(map[string]interface{})

	if !ok {
		// Set error header
		w.WriteHeader(500)

		// If AAC is running on development mode log error
		if util.Config.IsDev() || util.Config.IsLog() {
			util.Logger.Error("Cannot get session as map")
		}

		return
	}

	// Set LUA file name
	pageName := ps.ByName("filepath")

	// If there is no subtopic request index
	if pageName == "" {
		pageName = "index"
	}

	// Get state from the pool
	s, err := lua.PageList.Get(filepath.Join("pages", pageName, r.Method+".lua"))

	if err != nil {
		// Set error header
		w.WriteHeader(500)

		// If AAC is running on development mode log error
		if util.Config.IsDev() || util.Config.IsLog() {
			util.Logger.Errorf("Cannot get %v subtopic source code: %v", pageName, err)
		}

		return
	}

	// Create HTTP metatable
	lua.SetHTTPMetaTable(s)

	defer lua.PageList.Put(s, filepath.Join("pages", pageName, r.Method+".lua"))

	// Set the state user data
	lua.SetHTTPUserData(s, w, r)

	// Set session user data
	lua.SetSessionMetaTableUserData(s, session)

	if err := s.CallByParam(

		glua.P{
			Fn:      s.GetGlobal(strings.ToLower(r.Method)),
			NRet:    0,
			Protect: !util.Config.IsDev(),
		},
	); err != nil {

		// Set error header
		w.WriteHeader(500)

		// If AAC is running on development mode log error
		if util.Config.IsDev() || util.Config.IsLog() {
			util.Logger.Errorf("Cannot get %v subtopic source code: %v", pageName, err)
		}

		return
	}

	return

	/*
		// Get state from the pool
		luaState := lua.Pool.Get()

		// Defer the state put method
		defer lua.Pool.Put(luaState)

		// Set the state user data
		lua.SetHTTPUserData(luaState, w, r)

		// Set session user data
		lua.SetSessionMetaTableUserData(luaState, session)

		// Get function for the subtopic
		source, err := lua.Subtopics.Get("pages", pageName, r.Method)

		if err != nil {
			// Set error header
			w.WriteHeader(500)

			// If AAC is running on development mode log error
			if util.Config.IsDev() || util.Config.IsLog() {
				util.Logger.Errorf("Cannot get %v subtopic source code: %v", pageName, err)
			}

			return
		}

		// Execute function
		if err := luaState.DoString(source); err != nil {
			// Set error header
			w.WriteHeader(500)

			// If AAC is running on development mode log error
			if util.Config.IsDev() || util.Config.IsLog() {
				util.Logger.Errorf("Cannot execute %v subtopic: %v", pageName, err)
			}

			return
		}

		// Remove top stack value
		luaState.Pop(-1)*/
}
