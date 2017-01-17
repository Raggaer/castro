package main

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/raggaer/castro/app/util"
	"path/filepath"
)

func serveStatic(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Check if the asset already has gzip compression
	stats, err := os.Stat(filepath.Join("public", ps.ByName("filepath")))
	if err != nil {
		if err := gzipCompress("public/" + ps.ByName("filepath")); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	}

	// Check if file needs to be refreshed
	if time.Since(stats.ModTime()) > time.Minute*4 || util.Config.IsDev() {
		if err := gzipCompress("public/" + ps.ByName("filepath")); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
	}

	// Set gzip encoding header
	w.Header().Set("Content-Encoding", "gzip")

	// Set content type and content length
	switch {
	case strings.HasSuffix(ps.ByName("filepath"), ".css"):
		w.Header().Add("Content-Type", "text/css")
	case strings.HasSuffix(ps.ByName("filepath"), ".js"):
		w.Header().Add("Content-Type", "application/javascript")
	case strings.HasSuffix(ps.ByName("filepath"), ".otf"):
		w.Header().Add("Content-Type", "application/x-font-otf")
	case strings.HasSuffix(ps.ByName("filepath"), ".ttf"):
		w.Header().Add("Content-Type", "application/x-font-ttf")
	case strings.HasSuffix(ps.ByName("filepath"), ".woff"):
		w.Header().Add("Content-Type", "application/x-font-woff")
	case strings.HasSuffix(ps.ByName("filepath"), ".woff2"):
		w.Header().Add("Content-Type", "font/woff2")
	default:
		w.Header().Add("Content-Type", "text/plain")
	}

	// Serve file
	http.ServeFile(w, r, filepath.Join("public", ps.ByName("filepath")+".gzip"))
}

func gzipCompress(path string) error {
	// Create new gzip file for the asset
	gzipFile, err := os.Create(path + ".gzip")
	if err != nil {
		return err
	}

	defer gzipFile.Close()

	// Read all bytes from the requested asset
	asset, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Create the gzip writter
	gzipWritter := gzip.NewWriter(gzipFile)

	defer gzipWritter.Close()

	// Write asset data
	if _, err := gzipWritter.Write(asset); err != nil {
		return err
	}

	return nil
}
