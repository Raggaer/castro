package lua

import (
	"github.com/yuin/gopher-lua"
	"reflect"
)

// Config holds the current lua configuration file
var Config = &Configuration{}

// Configuration struct used for LUA configuration
// file
type Configuration struct {
	WorldType                           string `lua:"worldType"`
	HotkeyAimbotEnabled                 bool   `lua:"hotkeyAimbotEnabled"`
	ProtectionLevel                     int    `lua:"protectionLevel"`
	RedSkullKills                       int    `lua:"killsToRedSkull"`
	BlackSkullKills                     int    `lua:"killsToBlackSkull"`
	PzLocked                            int    `lua:"pzLocked"`
	RemoveChargesFromRunes              bool   `lua:"removeChargesFromRunes"`
	TimeToDecreaseFrags                 int    `lua:"timeToDecreaseFrags"`
	WhiteSkullTime                      int    `lua:"whiteSkullTime"`
	StairJumpExhaustion                 int    `lua:"stairJumpExhaustion"`
	ExperienceByKillingPlayers          bool   `lua:"experienceByKillingPlayers"`
	ExpFromPlayersLevelRange            int    `lua:"expFromPlayersLevelRange"`
	IP                                  string `lua:"ip"`
	BindOnlyGlobalAddress               bool   `lua:"bindOnlyGlobalAddress"`
	LoginProtocolPort                   int    `lua:"loginProtocolPort"`
	GameProtocolPort                    int    `lua:"gameProtocolPort"`
	StatusProtocolPort                  int    `lua:"statusProtocolPort"`
	MaxPlayers                          int    `lua:"maxPlayers"`
	Motd                                string `lua:"motd"`
	OnePlayerOnlinePerAccount           bool   `lua:"onePlayerOnlinePerAccount"`
	AllowClones                         bool   `lua:"allowClones"`
	ServerName                          string `lua:"serverName"`
	StatusTimeOut                       int    `lua:"statusTimeout"`
	ReplaceKickOnLogin                  bool   `lua:"eplaceKickOnLogin"`
	MaxPacketsPerSecond                 int    `lua:"maxPacketsPerSecond"`
	DeathLosePercent                    int    `lua:"deathLosePercent"`
	HousePriceEachSQM                   int    `lua:"housePriceEachSQM"`
	HouseRentPeriod                     string `lua:"houseRentPeriod"`
	TimeBetweenActions                  int    `lua:"timeBetweenActions"`
	TimeBetweenExActions                int    `lua:"timeBetweenExActions"`
	MapName                             string `lua:"mapName"`
	MapAuthor                           string `lua:"mapAuthor"`
	MarketOfferDuration                 int    `lua:"marketOfferDuration"`
	PremiumToCreateMarketOffer          bool   `lua:"premiumToCreateMarketOffer"`
	CheckExpiredMarketOffersEachMinutes int    `lua:"checkExpiredMarketOffersEachMinutes"`
	MaxMarketOffersAtATimePerPlayer     int    `lua:"maxMarketOffersAtATimePerPlayer"`
	MySQLHost                           string `lua:"mysqlHost"`
	MySQLUser                           string `lua:"mysqlUser"`
	MySQLPass                           string `lua:"mysqlPass"`
	MySQLDatabase                       string `lua:"mysqlDatabase"`
	MySQLPort                           int    `lua:"mysqlPort"`
	MySQLSock                           string `lua:"mysqlSock"`
	PasswordType                        string `lua:"passwordType"`
	AllowChangeOutfit                   bool   `lua:"allowChangeOutfit"`
	FreePremium                         bool   `lua:"freePremium"`
	KickIdlePlayerAfterMinutes          int    `lua:"kickIdlePlayerAfterMinutes"`
	MaxMessageBuffer                    int    `lua:"maxMessageBuffer"`
	EmoteSpells                         bool   `lua:"emoteSpells"`
	ClassicEquipmentSlots               bool   `lua:"classicEquipmentSlots"`
	RateExp                             int    `lua:"rateExp"`
	RateSkill                           int    `lua:"rateSkill"`
	RateLoot                            int    `lua:"rateLoot"`
	RateMagic                           int    `lua:"rateMagic"`
	RateSpawn                           int    `lua:"rateSpawn"`
	DeSpawnRange                        int    `lua:"deSpawnRange"`
	DeSpawnRadius                       int    `lua:"deSpawnRadius"`
	StaminaSystem                       bool   `lua:"staminaSystem"`
	WarnUnsafeScripts                   bool   `lua:"warnUnsafeScripts"`
	ConvertUnsafeScripts                bool   `lua:"convertUnsafeScripts"`
	DefaultPriority                     string `lua:"defaultPriority"`
	StartupDatabaseOptimization         bool   `lua:"startupDatabaseOptimization"`
	OwnerName                           string `lua:"ownerName"`
	OwnerEmail                          string `lua:"ownerEmail"`
	URL                                 string `lua:"url"`
	Location                            string `lua:"location"`
}

// SetConfigMetaTable sets the config metatable of the given state
func SetConfigMetaTable(luaState *lua.LState) {
	// Create and set Config metatable
	configMetaTable := luaState.NewTypeMetatable(ConfigMetaTableName)
	luaState.SetGlobal(ConfigMetaTableName, configMetaTable)

	// Set all Config metatable functions
	luaState.SetFuncs(configMetaTable, configMethods)
}

// LoadConfig loads the lua configuration file using
// lua vm to get the global variables
func LoadConfig(path string, dest *Configuration) error {
	L := lua.NewState()
	defer L.Close()
	if err := L.DoFile(path + "/config.lua"); err != nil {
		return err
	}
	return GetStructVariables(dest, L)
}

// GetConfigLuaValue gets a value from the config struct using reflect
func GetConfigLuaValue(L *lua.LState) int {
	// Get value of Config struct
	r := reflect.ValueOf(Config)

	// Get field by its name
	f := reflect.Indirect(r).FieldByName(L.ToString(2))

	// Switch field type to push to stack
	switch f.Kind() {
	case reflect.String:
		L.Push(lua.LString(f.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		L.Push(lua.LNumber(f.Int()))
	case reflect.Bool:
		L.Push(lua.LBool(f.Bool()))
	}

	return 1
}
