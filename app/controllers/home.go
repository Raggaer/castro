package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/util"
)

// Home is the main aac homepage
func Home(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := util.Template.Render(w, "home.html", nil); err != nil {
		util.Logger.Fatal(err)
	}
}
