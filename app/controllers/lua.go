package controllers

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"

	glua "github.com/yuin/gopher-lua"
)

// LuaPage executes the given lua page
func LuaPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get Castro cookie
	cookie, err := r.Cookie(util.Config.Cookies.Name)

	if err != nil {

		// If AAC is running on development mode log error
		if util.Config.IsDev() {
			util.Logger.Errorf("Cannot execute %v: %v\n", ps.ByName("page"), err)
		}

		// Throw error to user
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	// Get json web token claims from the cookie value
	_, err = util.ParseJWToken(cookie.Value)

	if err != nil {

		// If AAC is running on development mode log error
		if util.Config.IsDev() {
			util.Logger.Errorf("Cannot execute %v: %v\n", ps.ByName("page"), err)
		}

		// Throw error to user
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	sessionData, err := util.GetSession(cookie.Value)

	if err != nil {

		// If AAC is running on development mode log error
		if util.Config.IsDev() {
			util.Logger.Errorf("Cannot execute %v: %v\n", ps.ByName("page"), err)
		}

		// Throw error to user
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	// Get state from the pool
	luaState := lua.Pool.Get()

	// Defer the state put method
	defer lua.Pool.Put(luaState)

	// Get JWT metatable
	jwtMetaTable := luaState.GetTypeMetatable(lua.JWTMetaTable)

	// Set session field
	sess := luaState.NewUserData()
	sess.Value = sessionData
	luaState.SetField(jwtMetaTable, lua.JWTTokenName, sess)

	// Get HTTP metatable
	httpMetaTable := luaState.GetTypeMetatable(lua.HTTPMetaTableName)

	// Set HTTP method field
	luaState.SetField(httpMetaTable, lua.HTTPMetaTableMethodName, glua.LString(r.Method))

	// Set HTTP response writer field
	httpW := luaState.NewUserData()
	httpW.Value = w
	luaState.SetField(httpMetaTable, lua.HTTPResponseWriterName, httpW)

	// Set HTTP request field
	httpR := luaState.NewUserData()
	httpR.Value = r
	luaState.SetField(httpMetaTable, lua.HTTPRequestName, httpR)

	// Check if request is POST
	if r.Method == http.MethodPost {

		// Parse POST form
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// Set POST values as LUA table
		luaState.SetField(httpMetaTable, lua.HTTPPostValuesName, lua.URLValuesToTable(r.PostForm))
	}

	// Set LUA file name
	pageName := ps.ByName("page")

	// If there is no subtopic request index
	if pageName == "" {
		pageName = "index"
	}

	// Execute the requested page
	if err := luaState.DoFile("pages/" + pageName + "/" + r.Method + ".lua"); err != nil {

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
