package util

import (
	"bytes"
	glua "github.com/yuin/gopher-lua"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sync"
)

var (
	// WidgetTemplate holds all the widget templates
	WidgetTemplate Tmpl

	// Widgets holds all the widget results
	Widgets = WidgetList{
		rw: &sync.RWMutex{},
	}
)

type WidgetList struct {
	rw   *sync.RWMutex
	List []*Widget
}

type Widget struct {
	rw     *sync.RWMutex
	Name   string
	Result template.HTML
}

// Load loads all the widgets from the given directory
func (w *WidgetList) Load(path string) error {
	// Lock widget list
	w.rw.Lock()

	// Unlock widget list
	defer w.rw.Unlock()

	// Load all directories of the widget folder
	widgets, err := ioutil.ReadDir(path)

	if err != nil {
		return err
	}

	// Set slice
	w.List = []*Widget{}

	// Loop folder items
	for _, file := range widgets {

		// Check if file is a directory
		if file.IsDir() {

			// Append widget
			w.List = append(w.List, &Widget{
				Name: file.Name(),
				rw:   &sync.RWMutex{},
			})
		}
	}

	return nil
}

// Execute gets the result of the given widget
func (w *Widget) Execute(luaState *glua.LState) (*glua.LTable, error) {
	// Execute lua file
	if err := luaState.DoFile(filepath.Join("widgets", w.Name, w.Name+".lua")); err != nil {
		return nil, err
	}

	// Get value of the top of the stack
	v := luaState.Get(-1)

	// Check for valid returned type
	if v.Type() != glua.LTTable {
		return nil, nil
	}

	// Return value as table
	return v.(*glua.LTable), nil
}

// SetResult changes a widget result
func (w *Widget) SetResult(buff *bytes.Buffer) {
	// Lock mutex
	w.rw.Lock()
	defer w.rw.Unlock()

	// Set result
	w.Result = template.HTML(buff.String())
}
