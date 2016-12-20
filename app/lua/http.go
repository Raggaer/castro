package lua

import (
	glua "github.com/yuin/gopher-lua"
	"net/http"
	"github.com/raggaer/castro/app/util"
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

	// Get args table as LUA value
	tableValue := L.Get(3)

	// Check if args is set
	if tableValue.Type() == glua.LTTable {

		// Render template with args
		util.Template.RenderTemplate(w, req, L.ToString(2), TableToMap(L.ToTable(3)))

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