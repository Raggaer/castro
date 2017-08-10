package controllers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/util"
	"net/http"
)

func ExtensionStatic(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Get extension identifier
	id := ps.ByName("id")

	// Check if static file exists
	dir, exists := util.ExtensionStatic.FileExists(id)

	if !exists {
		w.WriteHeader(404)
		return
	}

	// Open desired file
	f, err := dir.Open(ps.ByName("filepath"))

	if err != nil {
		w.WriteHeader(404)
		return
	}

	// Close file handle
	defer f.Close()

	// Get file information
	fi, err := f.Stat()

	if err != nil {
		w.WriteHeader(404)
		return
	}

	// Check if file is directory
	if fi.IsDir() {
		w.WriteHeader(404)
		return
	}

	// Serve file
	http.ServeContent(w, req, fi.Name(), fi.ModTime(), f)
}
