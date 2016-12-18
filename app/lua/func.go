package lua

import glua "github.com/yuin/gopher-lua"

// RenderTemplate sets the _template lua variable
// to later render that template on the controller
func RenderTemplate(L *glua.LState) int {
	L.SetGlobal(
		TemplateVarName,
		L.Get(1),
	)
	return 0
}
