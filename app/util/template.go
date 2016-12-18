package util

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type tmpl struct {
	tmpl *template.Template
}

var (
	// Template holds all the app templates
	Template tmpl

	// FuncMap holds the main FuncMap definition
	FuncMap template.FuncMap
)

// NewTemplate creates and returns a new tmpl instance
func NewTemplate(name string) tmpl {
	return tmpl{
		tmpl: template.New(name),
	}
}

// LoadTemplates parses and loads all template into
// the given variable
func LoadTemplates(t *tmpl) error {
	// Walk over the views directory
	err := filepath.Walk("views/", func(path string, info os.FileInfo, err error) error {

		// Check if file has .html extesnion
		if strings.HasSuffix(info.Name(), ".html") {
			if t.tmpl, err = t.tmpl.ParseFiles(path); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (t *tmpl) FuncMap(f template.FuncMap) {
	t.tmpl.Funcs(f)
}

func (t tmpl) RenderTemplate(w http.ResponseWriter, req *http.Request, name string, args map[interface{}]interface{}) {
	// Check if app is running on dev mode
	if Config.IsDev() {

		// Create new template
		t = NewTemplate("castro")

		// Set template FuncMap
		t.tmpl.Funcs(FuncMap)

		// Reload all templates
		if err := LoadTemplates(&t); err != nil {
			Logger.Error(err)
			return
		}
	}

	// Check if args is a valid map
	if args == nil {
		args = map[interface{}]interface{}{}
	}

	// Load microtime from the microtimeHandler
	microtime, ok := req.Context().Value("microtime").(time.Time)
	if !ok {
		w.WriteHeader(500)
		w.Write([]byte("Cannot read microtime value"))
		return
	}

	// Set microtime value
	args["microtime"] = time.Since(microtime)

	// Render template and log error
	if err := t.tmpl.ExecuteTemplate(w, name, args); err != nil {
		Logger.Error(err)
	}
}

// Render executes the given template. if the app is running
// on dev mode all the templates will be reloaded
func (t tmpl) Render(wr io.Writer, name string, args interface{}) error {
	// Check if app is running on dev mode
	if Config.IsDev() {

		// Reload all templates
		if err := LoadTemplates(&t); err != nil {
			return err
		}
	}

	// Execute template and return error
	return t.tmpl.ExecuteTemplate(wr, name, args)
}
