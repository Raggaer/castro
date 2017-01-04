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

var (
	// Pool saves all lua state pointers to create a sync.Pool
	Pool = &luaStatePool{
		saved: make([]*glua.LState, 0, 10),
	}

	cryptoMethods = map[string]glua.LGFunction{
		"sha1": Sha1Hash,
	}
	mysqlMethods = map[string]glua.LGFunction{
		"query": Query,
		"execute": Execute,
	}
	configMethods = map[string]glua.LGFunction{
		"get": GetConfigValue,
	}
	httpMethods = map[string]glua.LGFunction{
		"redirect": Redirect,
		"render": RenderTemplate,
	}
	validatorMethods = map[string]glua.LGFunction{
		"validate": Validate,
		"blackList": BlackList,
	}
	sessionMethods = map[string]glua.LGFunction{
		"isLogged": IsLogged,
		"getFlash": GetFlash,
		"setFlash": SetFlash,
		"set": SetSessionData,
		"get": GetSessionData,
		"destroy": DestroySession,
	}
	captchaMethods = map[string]glua.LGFunction{
		"isEnabled": IsEnabled,
		"verify": VerifyCaptcha,
	}
	mapMethods = map[string]glua.LGFunction{
		"houseList": HouseList,
		"townList": TownList,
		"townByID": GetTownByID,
	}
)

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

	// Remove HTTP metatable
	//state.SetGlobal(HTTPMetaTableName, glua.LNil)

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

	// Create and set the captcha metatable
	captchaMetaTable := luaState.NewTypeMetatable(CaptchaMetaTableName)
	luaState.SetGlobal(CaptchaMetaTableName, captchaMetaTable)

	// Set all captcha metatable functions
	luaState.SetFuncs(captchaMetaTable, captchaMethods)

	// Create and set the crypto metatable
	cryptoMetaTable := luaState.NewTypeMetatable(CryptoMetaTableName)
	luaState.SetGlobal(CryptoMetaTableName, cryptoMetaTable)

	// Set all crypto metatable functions
	luaState.SetFuncs(cryptoMetaTable, cryptoMethods)

	// Create and set the validator metatable
	validMetaTable := luaState.NewTypeMetatable(ValidatorMetaTableName)
	luaState.SetGlobal(ValidatorMetaTableName, validMetaTable)

	// Set all validator metatable functions
	luaState.SetFuncs(validMetaTable, validatorMethods)

	// Create and set the MySQL metatable
	mysqlMetaTable := luaState.NewTypeMetatable(MySQLMetaTableName)
	luaState.SetGlobal(MySQLMetaTableName, mysqlMetaTable)

	// Set all MySQL metatable functions
	luaState.SetFuncs(mysqlMetaTable, mysqlMethods)

	// Create and set the json web token metatable
	jwtMetaTable := luaState.NewTypeMetatable(JWTMetaTable)
	luaState.SetGlobal(JWTMetaTable, jwtMetaTable)

	// Set all jwt functions
	luaState.SetFuncs(jwtMetaTable, sessionMethods)

	// Create and set Config metatable
	configMetaTable := luaState.NewTypeMetatable(ConfigMetaTableName)
	luaState.SetGlobal(ConfigMetaTableName, configMetaTable)

	// Set all Config metatable functions
	luaState.SetFuncs(configMetaTable, configMethods)

	// Create and set HTTP metatable
	httpMetaTable := luaState.NewTypeMetatable(HTTPMetaTableName)
	luaState.SetGlobal(HTTPMetaTableName, httpMetaTable)

	// Set all HTTP metatable functions
	luaState.SetFuncs(httpMetaTable, httpMethods)

	// Create and set map metatable
	mapMetaTable := luaState.NewTypeMetatable(MapMetaTableName)
	luaState.SetGlobal(MapMetaTableName, mapMetaTable)

	// Set all map metatable functions
	luaState.SetFuncs(mapMetaTable, mapMethods)

	// Return the lua state
	return luaState
}
