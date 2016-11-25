package lua

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	glua "github.com/yuin/gopher-lua"
	"sync"
)

type LuaState struct {
	L        *glua.LState
	Request  *http.Request
	Response http.ResponseWriter
	Params   httprouter.Params
}

type luaStatePool struct {
	m sync.Mutex
	saved []*glua.LState
}

// Pool saves all lua state pointers to create a sync.Pool
var Pool = &luaStatePool{
	saved: make([]*glua.LState, 0, 10),
}

func (p *luaStatePool) Get() *glua.LState {
	p.m.Lock()
	defer p.m.Unlock()
	if (len(p.saved)) == 0 {
		return p.New()
	}
	x := p.saved[len(p.saved) - 1]
	p.saved = p.saved[0:len(p.saved) - 1]
	return x
}

func (p *luaStatePool) Put(state *glua.LState) {
	p.m.Lock()
	defer p.m.Unlock()
	p.saved = append(p.saved, state)
}

func (p *luaStatePool) New() *glua.LState {
	return glua.NewState(
		glua.Options{
			IncludeGoStackTrace: true,
		},
	)
}
