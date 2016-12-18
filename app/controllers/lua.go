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

	// Execute the requested page
	if err := luaState.DoFile("pages/" + ps.ByName("page") + ".lua"); err != nil {

		// If AAC is running on development mode log error
		if util.Config.IsDev() {
			util.Logger.Errorf("Cannot execute %v: %v\n", ps.ByName("page"), err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(400)
		w.Write([]byte("Cannot execute the given subtopic"))
	}

	// Get redirect location
	redirectLocation := luaState.GetGlobal(lua.RedirectVarName).String()

	// If user needs to be redirected
	if redirectLocation != "" {

		// Redirect user to the location
		http.Redirect(w, r, redirectLocation, 302)

		// Stop execution
		return
	}

	// Get template name to render
	templateName := luaState.GetGlobal(lua.TemplateVarName).String()

	// If there is a template to be rendered
	if templateName != "" {

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
