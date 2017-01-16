package util

import (
	glua "github.com/yuin/gopher-lua"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
)

var (
	// Widgets holds all the AAC widgets
	Widgets = WidgetList{
		rw: &sync.RWMutex{},
	}

	// WidgetTemplate holds all the widget templates
	WidgetTemplate Tmpl
)

// WidgetList struct that holds a list of application widgets
type WidgetList struct {
	List []*Widget
	rw   *sync.RWMutex
}

// Widget struct used to hold widget information
type Widget struct {
	Name   string
	Result template.HTML
}

// LoadWidgetList parses and returns a list with widget names
func (w *WidgetList) LoadWidgetList(dir string) error {
	w.rw.Lock()
	defer w.rw.Unlock()

	// Load all directories of the widget folder
	widgets, err := ioutil.ReadDir(dir)

	if err != nil {
		return err
	}

	w.List = []*Widget{}

	// Loop folder items
	for _, file := range widgets {

		// Check if file is a directory
		if file.IsDir() {

			// Append file
			w.List = append(w.List, &Widget{
				Name: file.Name(),
			})
		}
	}

	return nil
}

// ExecuteWidget runs the given widget and returns its output
func (w *Widget) ExecuteWidget(res http.ResponseWriter, req *http.Request, luaState *glua.LState) (*glua.LTable, error) {
	// Execute widget LUA page
	if err := luaState.DoFile(
		"widgets/" + w.Name + "/" + w.Name + ".lua",
	); err != nil {
		return nil, err
	}

	// Get widget returning table if any
	args := luaState.Get(-1)

	// If user does not return anything stop
	if args.Type() != glua.LTTable {
		return nil, nil
	}

	return args.(*glua.LTable), nil
}
