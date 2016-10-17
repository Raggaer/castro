package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

// Start the main execution point for Castro
func Start() {
	setLogger()
	mux := httprouter.New()
	n := negroni.Classic()
	n.UseHandler(mux)
	l.Error("hola")
	if err := http.ListenAndServe(":8080", n); err != nil {

	}
}
