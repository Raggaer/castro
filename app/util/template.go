package util

import (
	"bytes"
	"fmt"
	"github.com/kataras/go-errors"
	"github.com/raggaer/castro/app/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	// Template holds all the app templates
	Template Tmpl

	// FuncMap holds the main FuncMap definition
	FuncMap template.FuncMap
)

// Tmpl struct that holds an application template wrapper for the Go template used in the lua bindings
type Tmpl struct {
	rw   *sync.RWMutex
	Tmpl *template.Template
}

// NewTemplate creates and returns a new tmpl instance
func NewTemplate(name string) Tmpl {
	return Tmpl{
		rw:   &sync.RWMutex{},
		Tmpl: template.New(name),
	}
}

// LoadTemplates parses and loads all template into the given variable
func (t *Tmpl) LoadTemplates(dir string) error {
	// Lock mutex
	t.rw.Lock()
	defer t.rw.Unlock()

	// Walk over the views directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		// Check if file has .html extension
		if strings.HasSuffix(info.Name(), ".html") {
			if t.Tmpl, err = t.Tmpl.ParseFiles(path); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// FuncMap returns the template map of functions
func (t *Tmpl) FuncMap(f template.FuncMap) {
	// Lock mutex
	t.rw.Lock()
	defer t.rw.Unlock()

	t.Tmpl.Funcs(f)
}

// RenderWidget renders the given widget template
func (t Tmpl) RenderWidget(req *http.Request, name string, args map[string]interface{}) (*bytes.Buffer, error) {
	// Check if app is running on dev mode
	if Config.IsDev() {

		// Lock mutex
		t.rw.Lock()
		defer t.rw.Unlock()

		// Create new template
		t = NewTemplate("widget")

		// Set template FuncMap
		t.Tmpl.Funcs(FuncMap)

		// Reload all templates
		if err := t.LoadTemplates("widgets/"); err != nil {
			return nil, err
		}
	}

	// Get csrf token
	tkn, ok := req.Context().Value("csrf-token").(*models.CsrfToken)
	if !ok {
		return nil, errors.New("Cannot get CSRF token")
	}

	// Set token value
	args["csrfToken"] = tkn.Token

	// Data holder
	buff := &bytes.Buffer{}

	// Render template to buffer
	if err := t.Tmpl.ExecuteTemplate(buff, name, args); err != nil {
		return nil, err
	}

	return buff, nil
}

// RenderTemplate render the given template passing some values and loading all templates if in development mode
func (t Tmpl) RenderTemplate(w http.ResponseWriter, req *http.Request, name string, args map[string]interface{}) {
	// Check if app is running on dev mode
	if Config.IsDev() {

		// Lock mutex
		t.rw.Lock()
		defer t.rw.Unlock()

		// Create new template
		t = NewTemplate("castro")

		// Set template FuncMap
		t.Tmpl.Funcs(FuncMap)

		// Reload all templates
		if err := t.LoadTemplates("views/"); err != nil {
			Logger.Error(err.Error())
			return
		}

		// Reload all templates
		if err := t.LoadTemplates("pages/"); err != nil {
			Logger.Error(err.Error())
			return
		}
	}

	// Check if args is a valid map
	if args == nil {
		args = map[string]interface{}{}
	}

	// Load microtime from the microtimeHandler
	microtime, ok := req.Context().Value("microtime").(time.Time)
	if !ok {
		w.WriteHeader(500)
		w.Write([]byte("Cannot read microtime value"))
		return
	}

	// Get csrf token
	tkn, ok := req.Context().Value("csrf-token").(*models.CsrfToken)
	if !ok {
		w.WriteHeader(500)
		w.Write([]byte("Cannot read csrf token value"))
		return
	}

	// Set token value
	args["csrfToken"] = tkn.Token

	// Set microtime value
	args["microtime"] = fmt.Sprintf("%9.4f seconds", time.Since(microtime).Seconds())

	// Render template and log error
	if err := t.Tmpl.ExecuteTemplate(w, name, args); err != nil {
		Logger.Error(err.Error())
	}
}

// Render executes the given template. if the app is running
// on dev mode all the templates will be reloaded
func (t Tmpl) Render(wr io.Writer, name string, args interface{}) error {
	// Check if app is running on dev mode
	if Config.IsDev() {

		// Lock mutex
		t.rw.Lock()
		defer t.rw.Unlock()

		// Reload all templates
		if err := t.LoadTemplates("views/"); err != nil {
			return err
		}
	}

	// Execute template and return error
	return t.Tmpl.ExecuteTemplate(wr, name, args)
}
