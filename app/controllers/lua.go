package controllers

import (
	"github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"
	"net/http"
)

// LuaPage executes the given lua page
func LuaPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check if request is POST
	if r.Method == http.MethodPost {

		// Parse POST form
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	}

	// If development mode reload pages and widgets
	if util.Config.IsDev() {

		// Reload pages
		if err := lua.Subtopics.Load("pages"); err != nil {
			// Set error header
			w.WriteHeader(500)

			// If AAC is running on development mode log error
			if util.Config.IsDev() || util.Config.IsLog() {
				util.Logger.Errorf("Cannot execute %v: %v", ps.ByName("page"), err)
				w.Write([]byte(err.Error()))
				return
			}

			w.Write([]byte("Cannot execute the given subtopic"))
			return
		}

		// Reload widgets
		if err := lua.Widgets.Load("widgets"); err != nil {
			// Set error header
			w.WriteHeader(500)

			// If AAC is running on development mode log error
			if util.Config.IsDev() || util.Config.IsLog() {
				util.Logger.Errorf("Cannot execute %v: %v", ps.ByName("page"), err)
				w.Write([]byte(err.Error()))
				return
			}

			w.Write([]byte("Cannot execute the given subtopic"))
			return
		}
	}

	// Get session
	session := sessions.GetSession(r)

	// Get state from the pool
	luaState := lua.Pool.GetApplicationState()

	// Set the state user data
	lua.SetHTTPUserData(luaState, w, r)

	// Set session user data
	lua.SetSessionMetaTableUserData(luaState, session)

	// Defer the state put method
	defer lua.Pool.Put(luaState)

	// Set LUA file name
	pageName := ps.ByName("filepath")

	// If there is no subtopic request index
	if pageName == "" {
		pageName = "index"
	}

	// Get function for the subtopic
	source, err := lua.Subtopics.Get("pages", pageName, r.Method)

	if err != nil {
		// Set error header
		w.WriteHeader(500)

		// If AAC is running on development mode log error
		if util.Config.IsDev() || util.Config.IsLog() {
			util.Logger.Errorf("Cannot execute %v: %v\n", ps.ByName("page"), err)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("Cannot execute the given subtopic"))
		return
	}

	// Execute function
	if err := luaState.DoString(source); err != nil {
		// Set error header
		w.WriteHeader(500)

		// If AAC is running on development mode log error
		if util.Config.IsDev() || util.Config.IsLog() {
			util.Logger.Errorf("Cannot execute %v: %v\n", ps.ByName("page"), err)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("Cannot execute the given subtopic"))
		return
	}

	// Remove top stack value
	luaState.Pop(-1)
}
