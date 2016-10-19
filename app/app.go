package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/util"
	"github.com/urfave/negroni"
)

// Start the main execution point for Castro
func Start() {
	mux := httprouter.New()
	n := negroni.Classic()
	n.UseHandler(mux)
	util.Logger.Infof("Starting Castro http server on port: %v", 8080)
	if err := http.ListenAndServe(":8080", n); err != nil {
		util.Logger.Errorf("Cannot start Castro http server: %v", err)
	}
}
