package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/util"
)

// Home is the main aac homepage
func Home(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	util.Template.RenderTemplate(w, r, "home.html", nil)
}
