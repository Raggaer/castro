package lua

import (
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
)

// SetI18nMetaTable sets the i18n metatable of the given state
func SetI18nMetaTable(luaState *lua.LState) {
	// Create and set the file metatable
	i18nMetaTable := luaState.NewTypeMetatable(I18nMetaTableName)
	luaState.SetGlobal(I18nMetaTableName, i18nMetaTable)

	// Set all global metatable functions
	luaState.SetFuncs(i18nMetaTable, i18nMethods)
}

// SetI18nUserData sets the i18n language value
func SetI18nUserData(luaState *lua.LState, lang []string) {
	// Get metatable
	i18nMetatable := luaState.GetTypeMetatable(I18nMetaTableName)

	// Set language field
	luaState.SetField(i18nMetatable, "Language", StringSliceToTable(lang))
}

// GetLanguageIndex retrieves the given language index
func GetLanguageIndex(L *lua.LState) int {
	// Language file
	lang := L.ToString(2)

	// Index
	i := L.ToString(3)

	// Retrieve language file
	langFile, ok := util.LanguageFiles.Get(lang)
	if !ok {
		L.Push(lua.LNil)
		return 1
	}

	// Retrieve language index
	langStr, ok := langFile.Data[i]
	if !ok {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(lua.LString(langStr))
	return 1
}
