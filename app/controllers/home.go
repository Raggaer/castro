package controllers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/util"
)

func Home(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	util.Template.Render(w, "test.html", nil)
}
