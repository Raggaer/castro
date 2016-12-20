package controllers

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"

	glua "github.com/yuin/gopher-lua"
)

var (
	httpMethods = map[string]glua.LGFunction{
		"redirect": lua.Redirect,
		"render": lua.RenderTemplate,
	}
)

// LuaPage executes the given lua page
func LuaPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get state from the pool
	luaState := lua.Pool.Get()

	// Defer the state put method
	defer lua.Pool.Put(luaState)

	// Create and set HTTP metatable
	httpMetaTable := luaState.NewTypeMetatable(lua.HTTPMetaTableName)
	luaState.SetGlobal(lua.HTTPMetaTableName, httpMetaTable)
	luaState.SetField(httpMetaTable, lua.HTTPMetaTableMethodName, glua.LString(r.Method))
	httpW := luaState.NewUserData()
	httpW.Value = w
	luaState.SetField(httpMetaTable, lua.HTTPResponseWriterName, httpW)
	httpR := luaState.NewUserData()
	httpR.Value = r
	luaState.SetField(httpMetaTable, lua.HTTPRequestName, httpR)
	luaState.SetFuncs(httpMetaTable, httpMethods)

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

		// If AAC is running on development mode log error
		if util.Config.IsDev() {
			util.Logger.Errorf("Cannot execute %v: %v\n", ps.ByName("page"), err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(500)
		w.Write([]byte("Cannot execute the given subtopic"))
	}
}
