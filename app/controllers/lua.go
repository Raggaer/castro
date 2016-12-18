package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"
)

// LuaPage executes the given lua page
func LuaPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	luaState := lua.Pool.Get()
	defer lua.Pool.Put(luaState)
	if err := luaState.DoFile("pages/" + ps.ByName("page") + ".lua"); err != nil {
		if util.Config.IsDev() {
			util.Logger.Errorf("Cannot execute %v: %v\n", ps.ByName("page"), err)
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(400)
		w.Write([]byte("Cannot execute the given subtopic"))
	}
	templateName := luaState.GetGlobal(lua.TemplateVarName).String()
	if templateName != "" {
		util.Template.RenderTemplate(w, r, templateName, nil)
	}
}
