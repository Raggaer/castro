package lua

import (
	"github.com/asaskevich/govalidator"
	"github.com/yuin/gopher-lua"
)

// methods holds all the validation methods related to
// govalidator
var methods = map[string]govalidator.Validator{
	"IsURL": govalidator.IsURL,
	"IsAlpha": govalidator.IsAlpha,
	"IsAlphanumeric": govalidator.IsAlphanumeric,
	"IsEmail": govalidator.IsEmail,

}

// Validate executes the given govalidator func
// and returns its output
func Validate(L *lua.LState) int {
	// Get function name
	name := L.Get(2)

	// Check for invalid name
	if name.Type() != lua.LTString {

		// Raise argument error
		L.ArgError(2, "Invalid validatior name")

		return 0
	}

	// Get main argument to validate
	arg := L.Get(3)

	// Check for invalid argument
	if arg.Type() != lua.LTString {

		// Raise argument error
		L.ArgError(3, "Invalid validator object")

		return 0
	}

	v, ok := methods[name.String()]

	// Check if validator exists
	if !ok {

		// Raise argument error
		L.ArgError(2, "Unkown validator name")

		return 0
	}

	L.Push(lua.LBool(v(arg.String())))

	return 1
}