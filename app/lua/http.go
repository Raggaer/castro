package lua

import (
	"github.com/raggaer/castro/app/util"
	glua "github.com/yuin/gopher-lua"
	"net/http"
)

// SetHTTPMetaTable sets the http metatable on the given lua state
func SetHTTPMetaTable(luaState *glua.LState) {
	// Create and set HTTP metatable
	httpMetaTable := luaState.NewTypeMetatable(HTTPMetaTableName)
	luaState.SetGlobal(HTTPMetaTableName, httpMetaTable)

	// Set all HTTP metatable functions
	luaState.SetFuncs(httpMetaTable, httpMethods)
}

// SetWidgetHTTPMetaTable sets the widget http metatable on the given lua state
func SetWidgetHTTPMetaTable(luaState *glua.LState) {
	// Create and set HTTP metatable
	httpMetaTable := luaState.NewTypeMetatable(HTTPMetaTableName)
	luaState.SetGlobal(HTTPMetaTableName, httpMetaTable)
}

// SetHTTPUserData sets the http metatable user data
func SetHTTPUserData(luaState *glua.LState, w http.ResponseWriter, r *http.Request) {
	// Get metatable
	httpMetaTable := luaState.GetTypeMetatable(HTTPMetaTableName)

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

	// Set current subtopic
	luaState.SetField(httpMetaTable, HTTPCurrentSubtopic, glua.LString(r.RequestURI))
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

// WriteResponse writes string to the response writer
func WriteResponse(L *glua.LState) int {
	// Get HTTP request and HTTP response writer
	_, w := getRequestAndResponseWriter(L)

	// Get data
	data := L.Get(2)

	// Check valid data type
	if data.Type() != glua.LTString {
		L.ArgError(1, "Invalid data type. Expected string")
		return 0
	}

	// Write to response writer
	w.Write([]byte(data.String()))

	return 0
}

// RenderTemplate renders the given template with the given data as a LUA table
func RenderTemplate(L *glua.LState) int {
	// Get HTTP request and HTTP response writer
	req, w := getRequestAndResponseWriter(L)

	// Get session
	session := getSessionData(L)

	templateName := L.ToString(2)

	// Get args table as LUA value
	tableValue := L.Get(3)

	// Compile widget list
	widgets, err := compileWidgetList(req, w, session)

	if err != nil {
		util.Logger.Errorf("Cannot compile widget list: %v", err)
	}

	// Check if args is set
	if tableValue.Type() == glua.LTTable {

		// Convert table to map
		args := TableToMap(tableValue.(*glua.LTable))

		args["widgets"] = widgets

		// Render template with args
		util.Template.RenderTemplate(w, req, templateName, args)
		return 0
	}

	// Render template without args
	util.Template.RenderTemplate(w, req, templateName, map[string]interface{}{
		"widgets": widgets,
	})

	return 0
}

// Redirect redirects the user to the given location with a header
func Redirect(L *glua.LState) int {
	// Get HTTP request and HTTP response writer
	req, w := getRequestAndResponseWriter(L)

	// Get destination
	dest := L.Get(2)

	// If there is no destination redirect to current subtopic
	if dest.Type() == glua.LTNil {
		http.Redirect(w, req, req.RequestURI, 302)
		return 0
	}

	// Get status code
	header := L.ToInt(3)

	// Set default header
	if header == 0 {
		header = 302
	}

	// Redirect to the desired location
	http.Redirect(w, req, dest.String(), header)

	return 0
}

// ServeFile serves the given file
func ServeFile(L *glua.LState) int {
	// Get file path
	path := L.Get(2)

	// Check valid path type
	if path.Type() != glua.LTString {
		L.ArgError(1, "Invalid path type. Expected string")
		return 0
	}

	// Get request and response
	req, w := getRequestAndResponseWriter(L)

	// Serve file
	http.ServeFile(w, req, path.String())

	return 0
}
