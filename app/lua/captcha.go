package lua

import (
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
)

// SetCaptchaMetaTable sets the captcha metatable of the given state
func SetCaptchaMetaTable(luaState *lua.LState) {
	// Create and set the captcha metatable
	captchaMetaTable := luaState.NewTypeMetatable(CaptchaMetaTableName)
	luaState.SetGlobal(CaptchaMetaTableName, captchaMetaTable)

	// Set all captcha metatable functions
	luaState.SetFuncs(captchaMetaTable, captchaMethods)
}

// IsEnabled checks if the captcha service is enabled
func IsEnabled(L *lua.LState) int {
	// Push captcha status
	L.Push(
		lua.LBool(util.Config.Captcha.Enabled),
	)

	return 1
}

// VerifyCaptcha checks if the given captcha response is valid
func VerifyCaptcha(L *lua.LState) int {
	// Get captcha response
	answer := L.Get(2)

	// Check for valid response type
	if answer.Type() != lua.LTString {

		L.ArgError(1, "Invalid captcha response format. Expected string")
		return 0
	}

	// Verify captcha answer
	check, err := util.VerifyCaptcha(answer.String())

	if err != nil {

		L.RaiseError("Cannot verify captcha answer: %v", err)
		return 0
	}

	// Push verification status to stack
	L.Push(lua.LBool(check))

	return 1
}
