package lua

import (
	glua "github.com/yuin/gopher-lua"
	"sync"
)

// luaStatePool struct used for lua state pooling
type luaStatePool struct {
	m sync.Mutex
	saved []*glua.LState
}

// Pool saves all lua state pointers to create a sync.Pool
var Pool = &luaStatePool{
	saved: make([]*glua.LState, 0, 10),
}

// Get retrieves a lua state from the pool
// if no states are available we create one
func (p *luaStatePool) Get() *glua.LState {
	// Lock and unlock our mutex to prevent
	// data race
	p.m.Lock()
	defer p.m.Unlock()

	// If no states available create onw
	if (len(p.saved)) == 0 {
		return p.New()
	}

	// Return last state from the pool
	x := p.saved[len(p.saved) - 1]
	p.saved = p.saved[0:len(p.saved) - 1]
	return x
}

// Put saves a lua state back to the pool
func (p *luaStatePool) Put(state *glua.LState) {
	// Lock and unlock our mutex to prevent
	// data race
	p.m.Lock()
	defer p.m.Unlock()

	// Remove POST values
	state.SetGlobal(
		PostValuesName,
		glua.LNil,
	)

	// Remove redirect location
	state.SetGlobal(
		RedirectVarName,
		glua.LNil,
	)

	// Remove template name
	state.SetGlobal(
		TemplateVarName,
		glua.LNil,
	)

	// Remove template args
	state.SetGlobal(
		TemplateArgsVarName,
		glua.LNil,
	)

	// Append to the pool
	p.saved = append(p.saved, state)
}

// New creates and returns a lua state
func (p *luaStatePool) New() *glua.LState {
	// Create a new lua state
	luaState := glua.NewState(
		glua.Options{
			IncludeGoStackTrace: true,
		},
	)

	// Set lua state methods
	luaState.SetGlobal(TemplateFuncName, luaState.NewFunction(RenderTemplate))
	luaState.SetGlobal(RedirectFuncName, luaState.NewFunction(Redirect))

	// Return the lua state
	return luaState
}
