package lua

import (
	"bytes"
	"github.com/goincremental/negroni-sessions"
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/util"
	glua "github.com/yuin/gopher-lua"
	"html/template"
	"net/http"
	"sync"
)

// SetHTTPMetaTable sets the http metatable on the given
// lua state
func SetHTTPMetaTable(luaState *glua.LState, w http.ResponseWriter, r *http.Request) {
	// Create and set HTTP metatable
	httpMetaTable := luaState.NewTypeMetatable(HTTPMetaTableName)
	luaState.SetGlobal(HTTPMetaTableName, httpMetaTable)

	// Set all HTTP metatable functions
	luaState.SetFuncs(httpMetaTable, httpMethods)

	// Set HTTP method field
	luaState.SetField(httpMetaTable, HTTPMetaTableMethodName, glua.LString(r.Method))

	// Set HTTP response writer field
	httpW := luaState.NewUserData()
	httpW.Value = w
	luaState.SetField(httpMetaTable, HTTPResponseWriterName, httpW)

	// Set HTTP request field
	httpR := luaState.NewUserData()
	httpR.Value = r
	luaState.SetField(httpMetaTable, HTTPRequestName, httpR)

	// Set GET values as LUA table
	luaState.SetField(httpMetaTable, HTTPGetValuesName, URLValuesToTable(r.URL.Query()))

	// Check if request is POST
	if r.Method == http.MethodPost {

		// Set POST values as LUA table
		luaState.SetField(httpMetaTable, HTTPPostValuesName, URLValuesToTable(r.PostForm))
	}
}

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
		if err := util.Widgets.LoadWidgetList("widgets/"); err != nil {
			util.Logger.Fatalf("Cannot load widget list: %v", err)
		}

		// Create new template
		util.WidgetTemplate = util.NewTemplate("widget")

		// Set template FuncMap
		util.WidgetTemplate.Tmpl.Funcs(util.FuncMap)

		if err := util.LoadTemplates("widgets/", &util.WidgetTemplate); err != nil {
			L.RaiseError("Cannot execute widgets: %v", err)
			return 0
		}
	}

	// Get csrf token
	tkn, ok := req.Context().Value("csrf-token").(*models.CsrfToken)
	if !ok {
		L.RaiseError("Cannot get CSRF token")
		return 0
	}

	// Get session
	sess := getSessionData(L)

	// Create a wait group for the widgets
	wg := &sync.WaitGroup{}

	// Loop all widgets
	for _, widget := range util.Widgets.List {

		// Add a task
		wg.Add(1)

		// Execute widget to get the result
		go ex(widget, sess, wg, tkn.Token)
	}

	// Wait for the tasks
	wg.Wait()

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

func ex(widget *util.Widget, sess sessions.Session, wg *sync.WaitGroup, tkn string) {
	// End task
	defer wg.Done()

	// Get a new lua state
	L := Pool.Get()

	// Create XML metatable
	SetXMLMetaTable(L)

	// Create crypto metatable
	SetCryptoMetaTable(L)

	// Create validator metatable
	SetValidatorMetaTable(L)

	// Create database metatable
	SetDatabaseMetaTable(L)

	// Create config metatable
	SetConfigMetaTable(L)

	// Create session metatable
	SetSessionMetaTable(L, sess)

	// Create map metatable
	SetMapMetaTable(L)

	// Return state to the pool
	defer Pool.Put(L)

	// Execute widget
	result, err := widget.ExecuteWidget(L)

	if err != nil {

		widget.Result = template.HTML(
			"Cannot execute widget <b>" + widget.Name + "</b>: " + err.Error(),
		)
		return
	}

	// Hold template result
	tmplResult := &bytes.Buffer{}

	args := TableToMap(result)
	args["csrfToken"] = tkn

	// Execute widget template
	if err := util.WidgetTemplate.Tmpl.ExecuteTemplate(tmplResult, widget.Name+".html", args); err != nil {

		widget.Result = template.HTML(
			"Cannot execute widget <b>" + widget.Name + "</b>: " + err.Error(),
		)
		return
	}

	// Assign result to widget
	widget.Result = template.HTML(tmplResult.String())
}
