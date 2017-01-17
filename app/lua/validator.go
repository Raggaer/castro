package lua

import (
	"github.com/asaskevich/govalidator"
	"github.com/raggaer/castro/app/util"
	"github.com/yuin/gopher-lua"
	"regexp"
)

// methods holds all the validation methods related to
// govalidator
var methods = map[string]govalidator.Validator{
	"IsURL":          govalidator.IsURL,
	"IsAlpha":        govalidator.IsAlpha,
	"IsAlphanumeric": govalidator.IsAlphanumeric,
	"IsEmail":        govalidator.IsEmail,
	"IsJson":         govalidator.IsJSON,
	"IsNull":         govalidator.IsNull,
	"IsEmpty":        govalidator.IsNull,
	"IsASCII":        govalidator.IsASCII,
	"IsUpperCase":    govalidator.IsUpperCase,
	"IsLowerCase":    govalidator.IsLowerCase,
	"IsInt":          govalidator.IsInt,
}

// SetValidatorMetaTable sets the validator metatable of the given state
func SetValidatorMetaTable(luaState *lua.LState) {
	// Create and set the validator metatable
	validMetaTable := luaState.NewTypeMetatable(ValidatorMetaTableName)
	luaState.SetGlobal(ValidatorMetaTableName, validMetaTable)

	// Set all validator metatable functions
	luaState.SetFuncs(validMetaTable, validatorMethods)
}

// ValidVocation checks if the given vocation exists
func ValidVocation(L *lua.LState) int {
	// Get vocation value
	voc := L.Get(2)

	// Check is users wants to get base vocation
	base := L.ToBool(3)

	// Check for valid vocation type
	if voc.Type() != lua.LTString && voc.Type() != lua.LTNumber {

		L.ArgError(1, "Invalid vocation format. Expected number or string")
		return 0
	}

	// If vocation is number we assume its the vocation id
	if voc.Type() == lua.LTNumber {

		// Convert vocation to int
		vocid := L.ToInt(2)

		// Loop vocation list
		for _, voc := range util.ServerVocationList.List.Vocations {

			// If we find the vocation we are looking for
			if voc.ID == vocid {

				if base {

					// If its a base vocation return true
					if voc.FromVoc == voc.ID {

						// Vocation is found push true
						L.Push(lua.LBool(true))

						return 1
					}

					L.Push(lua.LBool(false))

					return 1
				}

				// Vocation is found push true
				L.Push(lua.LBool(true))

				return 1
			}
		}

		L.Push(lua.LBool(false))

		return 1
	}

	// If vocation is string we assume its the vocation name
	vocname := L.ToString(2)

	// Loop vocation list
	for _, voc := range util.ServerVocationList.List.Vocations {

		// If we find the vocation we are looking for
		if voc.Name == vocname {

			if base {

				// If its a base vocation return true
				if voc.FromVoc == voc.ID {

					// Vocation is found push true
					L.Push(lua.LBool(true))

					return 1
				}

				L.Push(lua.LBool(false))

				return 1
			}

			// Vocation is found push true
			L.Push(lua.LBool(true))

			return 1
		}
	}

	L.Push(lua.LBool(false))

	return 1
}

// ValidTown checks if the given town exists
func ValidTown(L *lua.LState) int {
	// Get town value
	town := L.Get(2)

	// Check for valid town type
	if town.Type() != lua.LTString && town.Type() != lua.LTNumber {

		L.ArgError(1, "Invalid town format. Expected number or string")
		return 0
	}

	// If town is number we assume its the town id
	if town.Type() == lua.LTNumber {

		// Convert town id to uint32
		townid := uint32(L.ToInt(2))

		// Check if town exists
		for _, town := range util.OTBMap.Towns {

			// If its the town we are looking for
			if town.ID == townid {

				// Town is found push true
				L.Push(lua.LBool(true))

				return 1
			}
		}

		L.Push(lua.LBool(false))

		return 1
	}

	// If town is string we assume its the town name
	townName := L.ToString(2)

	// Check if town exists
	for _, town := range util.OTBMap.Towns {

		// If its the town we are looking for
		if town.Name == townName {

			// Town is found push true
			L.Push(lua.LBool(true))

			return 1
		}
	}

	L.Push(lua.LBool(false))

	return 1
}

// ValidUsername checks if the given username contains only
// letters and spaces
func ValidUsername(L *lua.LState) int {
	// Get string to validate
	v := L.Get(2)

	// Check for valid type
	if v.Type() != lua.LTString {

		L.ArgError(1, "Invalid string format. Expected string")
		return 0
	}

	// Check against regexp
	match, err := regexp.MatchString("^[a-zA-Z ]*$", v.String())

	if err != nil {
		L.RaiseError("Cannot compare string against regexp: %v", err)
	}

	// Push regexp result
	L.Push(lua.LBool(match))

	return 1
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
		L.ArgError(1, "Unknown validator name")
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
