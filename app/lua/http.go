package lua

import (
	glua "github.com/yuin/gopher-lua"
	"net/http"
	"github.com/raggaer/castro/app/util"
	"bytes"
	"html/template"
)

func getRequestAndResponseWriter(L *glua.LState) (*http.Request, http.ResponseWriter) {
	// Get HTTP metatable
	metatable := L.GetTypeMetatable(HTTPMetaTableName)

	// Get HTTP request field
	req := L.GetField(metatable, HTTPRequestName).(*glua.LUserData).Value.(*http.Request)

	// Get HTTP response writer field
	w := L.GetField(metatable, HTTPResponseWriterName).(*glua.LUserData).Value.(http.ResponseWriter)

	return req, w
}

// RenderTemplate renders the given template
// with the given data as a LUA table
func RenderTemplate(L *glua.LState) int {
	// Get HTTP request and HTTP response writer
	req, w := getRequestAndResponseWriter(L)

	// If development mode compile all widget templates
	if util.Config.IsDev() {

		// Load widget list
		wList, err := util.LoadWidgetList("widgets/")

		if err != nil {
			util.Logger.Fatalf("Cannot load widget list: %v", err)
		}

		// Assign widget list to global variable
		util.WidgetList = wList

		// Create new template
		util.WidgetTemplate = util.NewTemplate("widget")

		// Set template FuncMap
		util.WidgetTemplate.Tmpl.Funcs(util.FuncMap)

		if err := util.LoadTemplates("widgets/", &util.WidgetTemplate); err != nil {
			L.RaiseError("Cannot execute widgets: %v", err)
			return 0
		}
	}

	// Loop all widgets
	for _, widget := range util.WidgetList {

		// Execute widget
		result, err := widget.ExecuteWidget(w, req, L)

		if err != nil {

			// Raise error if needed
			L.RaiseError("Cannot execute widgets: %v", err)
			return 0
		}

		// Hold template result
		tmplResult := &bytes.Buffer{}

		// Execute widget template
		if err := util.WidgetTemplate.Tmpl.ExecuteTemplate(tmplResult, widget.Name + ".html", TableToMap(result)); err != nil {
			L.RaiseError("Cannot execute widgets: %v", err)
			return 0
		}

		// Assign result to widget
		widget.Result = template.HTML(tmplResult.String())
	}

	// Get args table as LUA value
	tableValue := L.Get(3)

	// Check if args is set
	if tableValue.Type() == glua.LTTable {

		args := TableToMap(tableValue.(*glua.LTable))

		// Render template with args
		util.Template.RenderTemplate(w, req, L.ToString(2), args)
		return 0
	}

	// Render template without args
	util.Template.RenderTemplate(w, req, L.ToString(2), nil)

	return 0
}

// Redirect redirects the user to the given
// location with a 302 header
func Redirect(L *glua.LState) int {
	// Get HTTP request and HTTP response writer
	req, w := getRequestAndResponseWriter(L)

	// Redirect to the desired location
	http.Redirect(w, req, L.ToString(2), 302)

	return 0
}