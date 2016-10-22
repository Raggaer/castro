package util

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type tmpl struct {
	tmpl *template.Template
}

// Template holds all the app templates
var Template tmpl

// LoadTemplates parses and loads all template into
// the given variable
func LoadTemplates(t *tmpl) error {
	// Declare new template castro
	t.tmpl = template.New("castro")

	// Walk over the views directory
	err := filepath.Walk("views/"+Config.Dialect+"/", func(path string, info os.FileInfo, err error) error {

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
