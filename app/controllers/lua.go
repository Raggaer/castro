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
	// Get state from the pool
	luaState := lua.Pool.Get()

	// Defer the state put method
	defer lua.Pool.Put(luaState)

	// Set some lua state values
	luaState.SetGlobal(
		lua.HTTPMethodName,
		glua.LString(r.Method),
	)

	// If method is POST convert all the values
	// to a LUA table
	if r.Method == http.MethodPost {

		// Parse form to limit maximum number of bytes
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		// Set POST values as LUA table
		luaState.SetGlobal(
			lua.PostValuesName,
			lua.URLValuesToTable(r.PostForm),
		)
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

	// Get redirect location
	redirectLocation := luaState.GetGlobal(lua.RedirectVarName).String()

	// If user needs to be redirected
	if redirectLocation != "" && redirectLocation != "nil" {

		// Redirect user to the location
		http.Redirect(w, r, redirectLocation, 302)

		// Stop execution
		return
	}

	// Get template name to render
	templateName := luaState.GetGlobal(lua.TemplateVarName).String()

	// If there is a template to be rendered
	if templateName != "" && templateName != "nil" {

		// Get template arguments
		globalArgs := luaState.GetGlobal(lua.TemplateArgsVarName)

		// Check if arguments is a LUA table
		if globalArgs.Type() != glua.LTTable {

			// Execute template without arguments
			util.Template.RenderTemplate(w, r, templateName, nil)
			return
		}

		// Convert LUA table to Go map
		args := lua.TableToMap(
			globalArgs.(*glua.LTable),
		)

		// Execute template with arguments
		util.Template.RenderTemplate(w, r, templateName, args)
	}
}
