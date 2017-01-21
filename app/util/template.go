package util

import (
	"fmt"
	"github.com/raggaer/castro/app/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Tmpl struct {
	Tmpl *template.Template
}

var (
	// Template holds all the app templates
	Template Tmpl

	// FuncMap holds the main FuncMap definition
	FuncMap template.FuncMap
)

// NewTemplate creates and returns a new tmpl instance
func NewTemplate(name string) Tmpl {
	return Tmpl{
		Tmpl: template.New(name),
	}
}

// LoadTemplates parses and loads all template into
// the given variable
func LoadTemplates(dir string, t *Tmpl) error {
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

func (t *Tmpl) FuncMap(f template.FuncMap) {
	t.Tmpl.Funcs(f)
}

func (t Tmpl) RenderTemplate(w http.ResponseWriter, req *http.Request, name string, args map[string]interface{}) {
	// Check if app is running on dev mode
	if Config.IsDev() {

		// Create new template
		t = NewTemplate("castro")

		// Set template FuncMap
		t.Tmpl.Funcs(FuncMap)

		// Reload all templates
		if err := LoadTemplates("views/", &t); err != nil {
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
	args["_csrf"] = tkn.Token

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

		// Reload all templates
		if err := LoadTemplates("views/", &t); err != nil {
			return err
		}
	}

	// Execute template and return error
	return t.Tmpl.ExecuteTemplate(wr, name, args)
}
