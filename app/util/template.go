package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/database"
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

// LoadExtensionTemplates parses and loads all extension templates into the given variable
func (t *Tmpl) LoadExtensionTemplates(extType string) error {
	// Lock mutex
	t.rw.Lock()
	defer t.rw.Unlock()

	// Get extensions from database
	rows, err := database.DB.Queryx(strings.Replace("SELECT extension_id FROM castro_extension_? WHERE enabled = 1", "?", extType, -1))

	if err != nil {
		return err
	}

	// Close rows
	defer rows.Close()

	// Loop rows
	for rows.Next() {

		// Hold extension id
		var extension_id string

		if err := rows.Scan(&extension_id); err != nil {
			return err
		}

		dir := filepath.Join("extensions", extension_id, extType)

		// Make sure that directory exist
		if _, err = os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				Logger.Logger.Errorf("Missing %v directory in extension %v", extType, extension_id)
			}
			continue
		}

		// Walk over the directory
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

			// Check if file has .html extension
			if strings.HasSuffix(info.Name(), ".html") {
				if t.Tmpl, err = t.Tmpl.ParseFiles(path); err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			// Log error
			Logger.Logger.Errorf("Cannot load extension template: %v", err)
			continue
		}
	}

	return nil
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
	if Config.Configuration.IsDev() {

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

		// Reload extension templates
		if err := t.LoadExtensionTemplates("widgets"); err != nil {
			Logger.Logger.Error(err.Error())
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
	if Config.Configuration.IsDev() {

		// Lock mutex
		t.rw.Lock()
		defer t.rw.Unlock()

		// Create new template
		t = NewTemplate("castro")

		// Set template FuncMap
		t.Tmpl.Funcs(FuncMap)

		// Reload all templates
		if err := t.LoadTemplates("views/"); err != nil {
			Logger.Logger.Error(err.Error())
			return
		}

		// Reload all templates
		if err := t.LoadTemplates("pages/"); err != nil {
			Logger.Logger.Error(err.Error())
			return
		}

		// Reload all extension templates
		if err := t.LoadExtensionTemplates("pages"); err != nil {
			Logger.Logger.Error(err.Error())
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
		Logger.Logger.Error(err.Error())
	}
}

// Render executes the given template. if the app is running on dev mode all the templates will be reloaded
func (t Tmpl) Render(wr io.Writer, name string, args interface{}) error {
	// Check if app is running on dev mode
	if Config.Configuration.IsDev() {

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
