package lua

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yuin/gopher-lua/parse"

	"github.com/kardianos/osext"
	"github.com/raggaer/castro/app/util"
	glua "github.com/yuin/gopher-lua"
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

	globalFuncList = map[string]func(l *glua.LState) int{
		"sleep":   ThreadSleep,
		"Player":  PlayerConstructor,
		"Guild":   GuildConstructor,
		"try":     TryCatch,
		"ternary": Ternary,
	}
	cryptoMethods = map[string]glua.LGFunction{
		"sha1":         Sha1Hash,
		"sha256":       Sha256Hash,
		"hmacsha256":   HmacSha256,
		"md5":          Md5Hash,
		"randomString": RandomString,
		"qr":           GenerateQRCode,
		"qrKey":        GenerateAuthSecretKey,
	}
	base64Methods = map[string]glua.LGFunction{
		"encode": Base64Encode,
		"decode": Base64Decode,
	}
	mysqlMethods = map[string]glua.LGFunction{
		"query":       Query,
		"execute":     Execute,
		"singleQuery": SingleQuery,
	}
	configMethods = map[string]glua.LGFunction{
		"get":       GetConfigLuaValue,
		"setCustom": SetConfigCustomValue,
	}
	httpMethods = map[string]glua.LGFunction{
		"setCookie":          SetCookie,
		"getCookie":          GetCookie,
		"redirect":           Redirect,
		"render":             RenderTemplate,
		"write":              WriteResponse,
		"serveFile":          ServeFile,
		"get":                GetRequest,
		"setHeader":          SetHeader,
		"postForm":           PostFormRequest,
		"getHeader":          GetHeader,
		"getRemoteAddress":   GetRemoteAddress,
		"curl":               CreateRequestClient,
		"formFile":           GetFormFile,
		"parseMultiPartForm": ParseMultiPartForm,
		"GetRelativeURL":     GetRelativeURL,
	}
	httpRegularMethods = map[string]glua.LGFunction{
		"curl":     CreateRequestClient,
		"postForm": PostFormRequest,
		"get":      GetRequest,
	}
	validatorMethods = map[string]glua.LGFunction{
		"validate":       Validate,
		"blackList":      BlackList,
		"validUsername":  ValidUsername,
		"validTown":      ValidTown,
		"validVocation":  ValidVocation,
		"validGuildName": ValidGuildName,
		"validGuildRank": ValidGuildRank,
		"validQRToken":   CheckQRCode,
		"validGender":    ValidGender,
		"escapeString":   EscapeString,
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
		"encode":     EncodeMap,
	}
	xmlMethods = map[string]glua.LGFunction{
		"vocationList":   VocationList,
		"vocationByName": GetVocationByName,
		"vocationByID":   GetVocationByID,
		"marshal":        MarshalXML,
		"unmarshal":      UnmarshalXML,
		"unmarshalFile":  UnmarshalXMLFile,
		"monsterList":    MonsterList,
		"monsterByName":  MonsterByName,
	}
	mailMethods = map[string]glua.LGFunction{
		"send": SendMail,
	}
	cacheMethods = map[string]glua.LGFunction{
		"get":    GetCacheValue,
		"set":    SetCacheValue,
		"delete": DeleteCacheValue,
	}
	debugMethods = map[string]glua.LGFunction{
		"value": DebugValue,
	}
	urlMethods = map[string]glua.LGFunction{
		"decode": DecodeURL,
		"encode": EncodeURL,
	}
	timeMethods = map[string]glua.LGFunction{
		"parseUnix":     ParseUnixTimestamp,
		"parseDuration": ParseDurationString,
		"parseDate":     ParseDate,
		"newDuration":   NewDuration,
	}
	reflectMethods = map[string]glua.LGFunction{
		"getGlobal": nil,
	}
	jsonMethods = map[string]glua.LGFunction{
		"marshal":       MarshalJSON,
		"unmarshal":     UnmarshalJSON,
		"unmarshalFile": UnmarshalJSONFile,
	}
	storageMethods = map[string]glua.LGFunction{
		"get": GetStorageValue,
		"set": SetStorageValue,
	}
	playerMethods = map[string]glua.LGFunction{
		"getAccountId":     GetPlayerAccountID,
		"isOnline":         IsPlayerOnline,
		"getBankBalance":   GetPlayerBankBalance,
		"setBankBalance":   SetPlayerBankBalance,
		"getStorageValue":  GetPlayerStorageValue,
		"setStorageValue":  SetPlayerStorageValue,
		"getVocation":      GetPlayerVocation,
		"getTown":          GetPlayerTown,
		"getGender":        GetPlayerGender,
		"getLevel":         GetPlayerLevel,
		"getPremiumEndsAt": GetPlayerPremiumEndsAt,
		"getPremiumTime":   GetPlayerPremiumTime,
		"getPremiumDays":   GetPlayerPremiumDays,
		"getName":          GetPlayerName,
		"getExperience":    GetPlayerExperience,
		"getCapacity":      GetPlayerCapacity,
		"getCustomField":   GetPlayerCustomField,
		"setCustomField":   SetPlayerCustomField,
		"getGuild":         GetPlayerGuild,
	}
	guildMethods = map[string]glua.LGFunction{
		"getOwner":   GetGuildOwner,
		"getMembers": GetGuildMembers,
		"getLeader":  GetGuildLeader,
	}
	widgetMethods = map[string]glua.LGFunction{
		"render": RenderWidgetTemplate,
	}
	eventsMethods = map[string]glua.LGFunction{
		"new": BackgroundEvent,
	}
	paypalMethods = map[string]glua.LGFunction{
		"createPayment":      CreatePaypalPayment,
		"paymentInformation": GetPaypalPayment,
		"executePayment":     ExecutePaypalPayment,
	}
	imgMethods = map[string]glua.LGFunction{
		"new": NewGoImage,
	}
	goimageMethods = map[string]glua.LGFunction{
		"writeText":     WriteGoImageText,
		"save":          SaveGoImage,
		"setBackground": SetBackgroundGoImage,
		"encode":        GetGoImageAsString,
	}
	fileMethods = map[string]glua.LGFunction{
		"mod":             GetFileModTime,
		"exists":          CheckFileExists,
		"getDirectories":  GetDirectories,
		"getFiles":        GetFiles,
		"createDirectory": CreateDirectory,
		"unzip":           UnzipFile,
	}
	envMethods = map[string]glua.LGFunction{
		"set": SetEnvVariable,
		"get": GetEnvVariable,
	}
	logMethods = map[string]glua.LGFunction{
		"error": LogError,
		"fatal": LogFatal,
		"info":  LogInfo,
	}
	globalMethods = map[string]glua.LGFunction{
		"set":    SetGlobalLuaValue,
		"get":    GetGlobalLuaValue,
		"delete": DeleteGlobalLuaValue,
	}
	formFileMethods = map[string]glua.LGFunction{
		"isValidPNG":       FormFileIsValidPNG,
		"isValidExtension": FormFileIsValidExtension,
		"contentType":      FormFileDetectContentType,
		"getFile":          GetFormFileByteArray,
		"saveFile":         SaveFormFile,
		"saveFileAsPNG":    SaveFormFileAsPNG,
		"SaveFileAsJPEG":   SaveFormFileAsJPEG,
	}
	outfitMethods = map[string]glua.LGFunction{
		"generate": GenerateOutfit,
	}
	extensionMethods = map[string]glua.LGFunction{
		"reload": ReloadExtensions,
	}
	i18nMethods = map[string]glua.LGFunction{
		"get": GetLanguageIndex,
	}
)

// CompileLua reads the passed lua file from disk and compiles it.
func CompileLua(filePath string) (*glua.FunctionProto, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}
	proto, err := glua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

// DoCompiledFile takes a FunctionProto, as returned by CompileLua, and runs it in the LState. It is equivalent
// to calling DoFile on the LState with the original source file.
func DoCompiledFile(state *glua.LState, proto *glua.FunctionProto) error {
	lfunc := state.NewFunctionFromProto(proto)
	state.Push(lfunc)
	return state.PCall(0, glua.MultRet, nil)
}

// OverwriteConfigFile gathers all external config file and pushes globals
func OverwriteConfigFile() error {
	// Load external config files
	list, err := util.LoadExternalConfigFiles()

	if err != nil {
		return err
	}

	// Create config state
	configState := glua.NewState()

	// Set castro metatables
	GetApplicationState(configState)

	// Loop list
	for _, config := range list {

		// Execute config file
		if err := configState.DoFile(config); err != nil {
			return err
		}
	}

	// Get app table
	appTable, ok := configState.GetGlobal("app").(*glua.LTable)

	if !ok {
		return errors.New("Cannot get app global as table")
	}

	// Get custom field
	customField := appTable.RawGetString("Custom").(*glua.LTable)

	if !ok {
		return errors.New("Cannot get app.Custom global as table")
	}

	// Convert table back to a map
	util.Config.Configuration.Custom = TableToMap(customField)

	return nil
}

// Get retrieves a lua state from the pool if no states are available we create one
func (p *luaStatePool) Get() *glua.LState {
	// Lock and unlock our mutex to prevent data race
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

// GetApplicationState returns a page configured lua state
func GetApplicationState(luaState *glua.LState) {
	// Create i18n metatable
	SetI18nMetaTable(luaState)

	// Create outfit metatable
	SetOutfitMetaTable(luaState)

	// Create global metatable
	SetGlobalMetaTable(luaState)

	// Create http regular metatable
	SetRegularHTTPMetaTable(luaState)

	// Create log metatable
	SetLogMetaTable(luaState)

	// Create env metatable
	SetEnvMetaTable(luaState)

	// Create file metatable
	SetFileMetaTable(luaState)

	// Create image metatable
	SetImageMetaTable(luaState)

	// Create paypal metatable
	SetPayPalMetaTable(luaState)

	// Create events metatable
	SetEventsMetaTable(luaState)

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

	// Create base64 metatable
	SetBase64MetaTable(luaState)

	// Create validator metatable
	SetValidatorMetaTable(luaState)

	// Create session metatable
	SetSessionMetaTable(luaState)

	// Create database metatable
	SetDatabaseMetaTable(luaState)

	// Create config metatable
	SetConfigMetaTable(luaState)

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

	// Loop global functions map
	for funcName, luaFunc := range globalFuncList {

		// Set global function
		luaState.SetGlobal(funcName, luaState.NewFunction(luaFunc))
	}

	// Set global variables
	luaState.SetGlobal("serverPath", glua.LString(util.Config.Configuration.Datapack))
	luaState.SetGlobal("logFile",
		glua.LString(
			fmt.Sprintf("%v-%v-%v.json", util.Logger.LastLoggerDay.Year(), util.Logger.LastLoggerDay.Month(), util.Logger.LastLoggerDay.Day()),
		),
	)

	// Get executable folder
	f, err := osext.ExecutableFolder()

	if err != nil {
		util.Logger.Logger.Fatalf("Cannot get executable folder path: %v", err)
	}

	// Get package metatable
	pkg := luaState.GetGlobal("package")

	// Set path field
	luaState.SetField(
		pkg,
		"path",
		glua.LString(
			filepath.Join(f, "engine", "?.lua"),
		),
	)

	// Set config field
	SetConfigGlobal(luaState)
}

// SetConfigGlobal sets the config global value
func SetConfigGlobal(L *glua.LState) {
	// Create table
	tbl := L.NewTable()

	// Create Security table
	secTable := StructToTable(&util.Config.Configuration.Security)

	// Create CSP table
	cspTable := StructToTable(&util.Config.Configuration.Security.CSP)

	// Set CSP Frame table
	L.SetField(cspTable, "Frame", StructToTable(&util.Config.Configuration.Security.CSP.Frame))

	// Set CSP Script table
	L.SetField(cspTable, "Script", StructToTable(&util.Config.Configuration.Security.CSP.Script))

	// Set CSP Font table
	L.SetField(cspTable, "Font", StructToTable(&util.Config.Configuration.Security.CSP.Font))

	// Set CSP Connect table
	L.SetField(cspTable, "Connect", StructToTable(&util.Config.Configuration.Security.CSP.Connect))

	// Set CSP Style table
	L.SetField(cspTable, "Style", StructToTable(&util.Config.Configuration.Security.CSP.Style))

	// Set CSP Image table
	L.SetField(cspTable, "Image", StructToTable(&util.Config.Configuration.Security.CSP.Image))

	// Set CSP table inside Security table
	L.SetField(secTable, "CSP", cspTable)

	// Set Security table
	L.SetField(tbl, "Security", secTable)

	// Set Shop table
	L.SetField(tbl, "Shop", StructToTable(&util.Config.Configuration.Shop))

	// Set Plugin value
	L.SetField(tbl, "Plugin", StructToTable(&util.Config.Configuration.Plugin))

	// Set main value
	L.SetField(tbl, "Main", StructToTable(util.Config.Configuration))

	// Set PayPal value
	L.SetField(tbl, "PayPal", StructToTable(&util.Config.Configuration.PayPal))

	// Set Fortumo value
	L.SetField(tbl, "Fortumo", StructToTable(&util.Config.Configuration.Fortumo))

	// Set Captcha value
	L.SetField(tbl, "Captcha", StructToTable(&util.Config.Configuration.Captcha))

	// Set Mail value
	L.SetField(tbl, "Mail", StructToTable(&util.Config.Configuration.Mail))

	// Set Custom value
	L.SetField(tbl, "Custom", MapToTable(util.Config.Configuration.Custom))

	// Set PayGol value
	L.SetField(tbl, "PayGol", StructToTable(&util.Config.Configuration.PayGol))

	// Set SSL value
	L.SetField(tbl, "SSL", StructToTable(&util.Config.Configuration.SSL))

	// Set global value
	L.SetGlobal("app", tbl)

	// Set default fields
	L.SetField(tbl, "Version", glua.LString(util.VERSION))
	L.SetField(tbl, "BuildDate", glua.LString(util.BUILD_DATE))
	L.SetField(tbl, "CheckUpdates", glua.LBool(util.Config.Configuration.CheckUpdates))
	L.SetField(tbl, "URL", glua.LString(util.Config.Configuration.URL))
	L.SetField(tbl, "Port", glua.LNumber(util.Config.Configuration.Port))
	L.SetField(tbl, "Mode", glua.LString(util.Config.Configuration.Mode))
	L.SetField(tbl, "Datapack", glua.LString(util.Config.Configuration.Datapack))
}

// Put saves a lua state back to the pool
func (p *luaStatePool) Put(state *glua.LState) {
	// Lock and unlock our mutex to prevent data race
	p.m.Lock()
	defer p.m.Unlock()

	// Append to the pool
	p.saved = append(p.saved, state)
}

// New creates and returns a lua state
func (p *luaStatePool) New() *glua.LState {
	// Create a new lua state
	state := glua.NewState(
		glua.Options{
			IncludeGoStackTrace: util.Config.Configuration.IsDev(),
		},
	)

	// Set castro metatables
	GetApplicationState(state)

	// Return the lua state
	return state
}

// NewState creates and returns a new lua state
func NewState() *glua.LState {
	// Create a new lua state
	state := glua.NewState(
		glua.Options{
			IncludeGoStackTrace: util.Config.Configuration.IsDev(),
		},
	)

	// Set castro metatables
	GetApplicationState(state)

	// Return the lua state
	return state
}

// ExecuteFile calls lua dofile
func ExecuteFile(luaState *glua.LState, path string) error {
	// Execute lua file
	if err := luaState.DoFile(path); err != nil {
		return err
	}

	return nil
}

// ExecuteControllerPage executes the given subtopic using call by param
func ExecuteControllerPage(luaState *glua.LState, method string) error {
	// Call file function
	if err := luaState.CallByParam(
		glua.P{
			Fn:      luaState.GetGlobal(strings.ToLower(method)),
			NRet:    0,
			Protect: !util.Config.Configuration.IsDev(),
		},
	); err != nil {
		return err
	}

	return nil
}
