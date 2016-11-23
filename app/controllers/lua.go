package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/lua"
	"github.com/raggaer/castro/app/util"
)

// LuaPage executes the given lua page
func LuaPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	luaState := lua.NewState(r, w, ps)
	defer luaState.L.Close()
	if err := luaState.L.DoFile("pages/" + ps.ByName("page") + ".lual"); err != nil {
		util.Logger.Error(err.Error())
	}
}
