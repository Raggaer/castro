package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/controllers"
)

// PageNotFound executes a 404 lua page or a simple 404 page
func PageNotFound(w http.ResponseWriter, r *http.Request) {
	controllers.LuaPage(w, r, httprouter.Params{
		{
			"filepath",
			"404",
		},
	})
}
