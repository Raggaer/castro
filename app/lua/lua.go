package lua

import (
	"fmt"
	"github.com/kardianos/osext"
	"github.com/kataras/go-errors"
	"github.com/raggaer/castro/app/util"
	glua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

// luaStatePool struct used for lua state pooling
type luaStatePool struct {
	m     sync.Mutex
	saved []*glua.LState
}

// FunctionList list of lua source files
type FunctionList struct {
	rw   *sync.RWMutex
	List map[string]string
}

var (
	// Pool saves all lua state pointers to create a sync.Pool
	Pool = &luaStatePool{
		saved: make([]*glua.LState, 0, 10),
	}

	// Subtopics is the application list of lua subtopics
	Subtopics = FunctionList{
		rw: &sync.RWMutex{},
	}

	// Widgets is the application list of lua widgets
	Widgets = FunctionList{
		rw: &sync.RWMutex{},
	}

	cryptoMethods = map[string]glua.LGFunction{
		"sha1": Sha1Hash,
	}
	mysqlMethods = map[string]glua.LGFunction{
		"query":       Query,
		"execute":     Execute,
		"singleQuery": SingleQuery,
	}
	configMethods = map[string]glua.LGFunction{
		"get": GetConfigLuaValue,
	}
	httpMethods = map[string]glua.LGFunction{
		"redirect": Redirect,
		"render":   RenderTemplate,
	}
	validatorMethods = map[string]glua.LGFunction{
		"validate":       Validate,
		"blackList":      BlackList,
		"validUsername":  ValidUsername,
		"validTown":      ValidTown,
		"validVocation":  ValidVocation,
		"validGuildName": ValidGuildName,
		"validGuildRank": ValidGuildRank,
	}
	sessionMethods = map[string]glua.LGFunction{
		"isLogged":      IsLogged,
		"isAdmin":       IsAdmin,
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
		"vocationByID":   GetVocationByID,
		"marshal":        MarshalXML,
		"unmarshal":      UnmarshalXML,
		"unmarshalFile":  UnmarshalXMLFile,
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
	urlMethods = map[string]glua.LGFunction{
		"decode": DecodeURL,
		"encode": EncodeURL,
	}
	timeMethods = map[string]glua.LGFunction{
		"parseUnix": ParseUnixTimestamp,
	}
	reflectMethods = map[string]glua.LGFunction{
		"getGlobal": GetGlobal,
	}
	jsonMethods = map[string]glua.LGFunction{
		"marshal":   MarshalJSON,
		"unmarshal": UnmarshalJSON,
	}
	storageMethods = map[string]glua.LGFunction{
		"get": GetStorageValue,
		"set": SetStorageValue,
	}
	playerMethods = map[string]glua.LGFunction{
		"getAccountId":    GetPlayerAccountID,
		"isOnline":        IsPlayerOnline,
		"getBankBalance":  GetPlayerBankBalance,
		"getStorageValue": GetPlayerStorageValue,
		"setStorageValue": SetPlayerStorageValue,
		"getVocation":     GetPlayerVocation,
	}
)

// Load loads all lua source files
func (s *FunctionList) Load(dir string) error {
	// Lock mutex
	s.rw.Lock()
	defer s.rw.Unlock()

	// Set list
	s.List = make(map[string]string)

	// Get a state from the pool
	L := Pool.Get()

	// Return the state
	defer Pool.Put(L)

	// Get subtopic list
	subtopicList, err := util.GetLuaFiles(dir)

	if err != nil {
		return err
	}

	// Loop subtopic list
	for _, subtopic := range subtopicList {

		// Load file
		f, err := ioutil.ReadFile(subtopic)

		if err != nil {
			return err
		}

		// Push result
		s.List[subtopic] = string(f)
	}

	return nil
}

// Get retrieves the source of the given lua file
func (s *FunctionList) Get(place, name, method string) (string, error) {
	// Lock mutex
	s.rw.Lock()
	defer s.rw.Unlock()

	// Build path
	path := filepath.Join(place, strings.ToLower(name), strings.ToLower(method)+".lua")

	// Check if path exists
	f, ok := s.List[path]

	if !ok {
		return "", errors.New("Subtopic not found")
	}

	return f, nil
}

// Get retrieves a lua state from the pool if no states are available we create one
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

// GetPageState returns a page configured lua state
func (p *luaStatePool) GetApplicationState() *glua.LState {
	// Get state from the pool
	luaState := Pool.Get()

	// Create storage metatable
	SetStorageMetaTable(luaState)

	// Create time metatable
	SetTimeMetaTable(luaState)

	// Create url metatable
	SetURLMetaTable(luaState)

	// Create debug metatable
	SetDebugMetaTable(luaState)

	// Create XML metatable
	SetXMLMetaTable(luaState)

	// Create captcha metatable
	SetCaptchaMetaTable(luaState)

	// Create crypto metatable
	SetCryptoMetaTable(luaState)

	// Create validator metatable
	SetValidatorMetaTable(luaState)

	// Create session metatable
	SetSessionMetaTable(luaState)

	// Create database metatable
	SetDatabaseMetaTable(luaState)

	// Create config metatable
	SetConfigMetaTable(luaState)

	// Create HTTP metatable
	SetHTTPMetaTable(luaState)

	// Create map metatable
	SetMapMetaTable(luaState)

	// Create mail metatable
	SetMailMetaTable(luaState)

	// Create cache metatable
	SetCacheMetaTable(luaState)

	// Create reflect metatable
	SetReflectMetaTable(luaState)

	// Create json metatable
	SetJSONMetaTable(luaState)

	// Set player global
	luaState.SetGlobal("Player", luaState.NewFunction(PlayerConstructor))

	// Set last log file name
	luaState.SetGlobal("logFile", glua.LString(
		fmt.Sprintf("%v-%v-%v.json", util.LastLoggerDay.Year(), util.LastLoggerDay.Month(), util.LastLoggerDay.Day()),
	))

	// Set server path
	luaState.SetGlobal("serverPath", glua.LString(util.Config.Datapack))

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

	return luaState
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

	// Return the lua state
	return luaState
}
