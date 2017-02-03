package lua

import (
	"github.com/kardianos/osext"
	"github.com/raggaer/castro/app/util"
	glua "github.com/yuin/gopher-lua"
	"path/filepath"
	"sync"
)

// luaStatePool struct used for lua state pooling
type luaStatePool struct {
	m     sync.Mutex
	saved []*glua.LState
}

var (
	// Pool saves all lua state pointers to create a sync.Pool
	Pool = &luaStatePool{
		saved: make([]*glua.LState, 0, 10),
	}

	cryptoMethods = map[string]glua.LGFunction{
		"sha1": Sha1Hash,
	}
	mysqlMethods = map[string]glua.LGFunction{
		"query":   Query,
		"execute": Execute,
	}
	configMethods = map[string]glua.LGFunction{
		"get": GetConfigLuaValue,
	}
	httpMethods = map[string]glua.LGFunction{
		"redirect": Redirect,
		"render":   RenderTemplate,
	}
	validatorMethods = map[string]glua.LGFunction{
		"validate":      Validate,
		"blackList":     BlackList,
		"validUsername": ValidUsername,
		"validTown":     ValidTown,
		"validVocation": ValidVocation,
	}
	sessionMethods = map[string]glua.LGFunction{
		"isLogged":      IsLogged,
		"getFlash":      GetFlash,
		"setFlash":      SetFlash,
		"set":           SetSessionData,
		"get":           GetSessionData,
		"destroy":       DestroySession,
		"loggedAccount": GetLoggedAccount,
	}
	captchaMethods = map[string]glua.LGFunction{
		"isEnabled": IsEnabled,
		"verify":    VerifyCaptcha,
	}
	mapMethods = map[string]glua.LGFunction{
		"houseList":  HouseList,
		"townList":   TownList,
		"townByID":   GetTownByID,
		"townByName": GetTownByName,
	}
	xmlMethods = map[string]glua.LGFunction{
		"vocationList":   VocationList,
		"vocationByName": GetVocationByName,
	}
	mailMethods = map[string]glua.LGFunction{
		"send": SendMail,
	}
	cacheMethods = map[string]glua.LGFunction{
		"get": GetCacheValue,
		"set": SetCacheValue,
	}
	debugMethods = map[string]glua.LGFunction{
		"value": DebugValue,
	}
)

// Get retrieves a lua state from the pool
// if no states are available we create one
func (p *luaStatePool) Get() *glua.LState {
	// Lock and unlock our mutex to prevent
	// data race
	p.m.Lock()
	defer p.m.Unlock()

	// If no states available create one
	if (len(p.saved)) == 0 {
		return p.New()
	}

	// Return last state from the pool
	x := p.saved[len(p.saved)-1]
	p.saved = p.saved[0 : len(p.saved)-1]

	return x
}

// Put saves a lua state back to the pool
func (p *luaStatePool) Put(state *glua.LState) {
	// Lock and unlock our mutex to prevent
	// data race
	p.m.Lock()
	defer p.m.Unlock()

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

	// Get executable folder
	f, err := osext.ExecutableFolder()

	if err != nil {
		util.Logger.Fatalf("Cannot get executable folder path: %v", err)
	}

	// Get package metatable
	pkg := luaState.GetGlobal("package")

	// Set path field
	luaState.SetField(
		pkg,
		"path",
		glua.LString(
			filepath.Join(f, "app", "lua", "engine", "?.lua"),
		),
	)

	// Return the lua state
	return luaState
}
