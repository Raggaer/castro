package lua

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	glua "github.com/yuin/gopher-lua"
)

type LuaState struct {
	L        *glua.LState
	Request  *http.Request
	Response http.ResponseWriter
	Params   httprouter.Params
}

// NewState creates and returns a new lua.LState
// pointer with some options
func NewState(req *http.Request, res http.ResponseWriter, params httprouter.Params) *LuaState {
	luaState := glua.NewState(
		glua.Options{
			IncludeGoStackTrace: true,
		},
	)
	state := &LuaState{
		L:        luaState,
		Request:  req,
		Response: res,
		Params:   params,
	}
	defineMethods(state.L)
	return state
}

func defineMethods(L *glua.LState) {

}
