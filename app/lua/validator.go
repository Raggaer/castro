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
	"IsJson": govalidator.IsJSON,
	"IsNull": govalidator.IsNull,
	"IsEmpty": govalidator.IsNull,
	"IsASCII": govalidator.IsASCII,
	"IsUpperCase": govalidator.IsUpperCase,
	"IsLowerCase": govalidator.IsLowerCase,
	"IsInt": govalidator.IsInt,

}

// Validate executes the given govalidator func
// and returns its output
func Validate(L *lua.LState) int {
	// Get function name
	name := L.Get(2)

	// Check for invalid name
	if name.Type() != lua.LTString {

		// Raise argument error
		L.ArgError(1, "Invalid validatior name")
		return 0
	}

	// Get main argument to validate
	arg := L.Get(3)

	// Check for invalid argument
	if arg.Type() != lua.LTString {

		// Raise argument error
		L.ArgError(2, "Invalid validator object")
		return 0
	}

	v, ok := methods[name.String()]

	// Check if validator exists
	if !ok {

		// Raise argument error
		L.ArgError(1, "Unkown validator name")
		return 0
	}

	L.Push(lua.LBool(v(arg.String())))

	return 1
}

// BlackList runs govalidator blacklist func
func BlackList(L *lua.LState) int {
	// Get main string to compare
	line := L.Get(2)

	// Check for invalid line
	if line.Type() != lua.LTString {

		// Raise argument error
		L.ArgError(1, "Invalid object type. Expected string")
		return 0
	}

	// Get words for blacklist
	words := L.Get(3)

	// Check for invalid type of word
	if words.Type() != lua.LTString {

		// Raise argument error
		L.ArgError(2, "Invalid table of words. Expected string")
		return 0
	}

	// Call govalidator method and push result to stack
	L.Push(
		lua.LString(
			govalidator.BlackList(
				line.String(),
				words.String(),
			),
		),
	)

	return 1
}