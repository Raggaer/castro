package lua

import (
	glua "github.com/yuin/gopher-lua"
	"log"
	"net/http"
)

// RenderTemplate sets the _template lua variable
// to later render that template on the controller
func RenderTemplate(L *glua.LState) int {
	// Set template name
	L.SetGlobal(
		TemplateVarName,
		L.Get(1),
	)
	// Get template args if any
	args := L.Get(2)

	// If user does not send args proceed
	if args.Type() == glua.LTNil {
		return 0
	}

	// If args are not a table exit
	if args.Type() != glua.LTTable {
		return 0
	}

	// Save template args to global variable
	L.SetGlobal(TemplateArgsVarName, L.ToTable(2))

	// Dont return anything to LUA
	return 0
}

// Redirect redirects the user to the given
// location with a 302 header
func Redirect(L *glua.LState) int {
	// Get HTTP metatable
	metatable := L.GetTypeMetatable(HTTPMetaTableName)

	// Get HTTP request field
	req := L.GetField(metatable, HTTPRequestName).(*glua.LUserData).Value.(*http.Request)

	log.Println(req.Method)


	// Dont return anything to LUA
	return 0
}