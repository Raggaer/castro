package controllers

import (
	"github.com/goincremental/negroni-sessions"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"
	"net/http"
	"path/filepath"
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

	// Get session
	session := sessions.GetSession(r)

	// Get state from the pool
	luaState := lua.Pool.Get()

	// Defer the state put method
	defer lua.Pool.Put(luaState)

	// Create XML metatable
	lua.SetXMLMetaTable(luaState)

	// Create captcha metatable
	lua.SetCaptchaMetaTable(luaState)

	// Create crypto metatable
	lua.SetCryptoMetaTable(luaState)

	// Create validator metatable
	lua.SetValidatorMetaTable(luaState)

	// Create session metatable
	lua.SetSessionMetaTable(luaState, session)

	// Create database metatable
	lua.SetDatabaseMetaTable(luaState)

	// Create config metatable
	lua.SetConfigMetaTable(luaState)

	// Create HTTP metatable
	lua.SetHTTPMetaTable(luaState, w, r)

	// Create map metatable
	lua.SetMapMetaTable(luaState)

	// Create mail metatable
	lua.SetMailMetaTable(luaState)

	// Create cache metatable
	lua.SetCacheMetaTable(luaState)

	// Set LUA file name
	pageName := ps.ByName("page")

	// If there is no subtopic request index
	if pageName == "" {
		pageName = "index"
	}

	// Execute the requested page
	if err := luaState.DoFile(filepath.Join("pages", pageName, r.Method+".lua")); err != nil {

		// Set error header
		w.WriteHeader(500)

		// If AAC is running on development mode log error
		if util.Config.IsDev() {
			util.Logger.Errorf("Cannot execute %v: %v\n", ps.ByName("page"), err)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("Cannot execute the given subtopic"))
		return
	}
}
