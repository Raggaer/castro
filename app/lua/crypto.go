package lua

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/hmac"
	"encoding/hex"
	"encoding/base64"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/skip2/go-qrcode"
	"github.com/yuin/gopher-lua"
)

// SetCryptoMetaTable sets the crypto metatable of the given state
func SetCryptoMetaTable(luaState *lua.LState) {
	// Create and set the crypto metatable
	cryptoMetaTable := luaState.NewTypeMetatable(CryptoMetaTableName)
	luaState.SetGlobal(CryptoMetaTableName, cryptoMetaTable)

	// Set all crypto metatable functions
	luaState.SetFuncs(cryptoMetaTable, cryptoMethods)
}

// Sha1Hash returns the sha1 hash of the given string
func Sha1Hash(L *lua.LState) int {
	// Get string to be hashed
	str := L.Get(2)

	// Check for valid string type
	if str.Type() != lua.LTString {

		L.ArgError(1, "Invalid string format. Expected string")
		return 0
	}

	// Hash string using sha1
	data := sha1.Sum([]byte(str.String()))

	// Convert byte array to string and push tu stack
	L.Push(
		lua.LString(
			fmt.Sprintf("%x", data),
		),
	)

	return 1
}

// Sha256Hash returns the sha256 hash of the given string
func Sha256Hash(L *lua.LState) int {
	// Get string to be hashed
	str := L.Get(2)

	// Check for valid string type
	if str.Type() != lua.LTString {

		L.ArgError(1, "Invalid string format. Expected string")
		return 0
	}

	// Hash string using sha1
	data := sha256.Sum256([]byte(str.String()))

	// Convert byte array to string and push to stack
	L.Push(
		lua.LString(
			fmt.Sprintf("%x", data),
		),
	)

	return 1
}

// HmacSha256 returns the hmac-sha256 for the given message + secret
func HmacSha256(L *lua.LState) int {
	// Get message string to be hashed
	message := L.Get(2)

	// Check for valid string type
	if message.Type() != lua.LTString {

		L.ArgError(1, "Invalid message format. Expected string")
		return 0
	}

	// Get secret string
	secret := L.Get(3)

	// Check for valid string type
	if secret.Type() != lua.LTString {

		L.ArgError(1, "Invalid secret format. Expected string")
		return 0
	}

	hash := hmac.New(sha256.New, []byte(secret.String()))
	hash.Write([]byte(message.String()))

	// Convert hash to hex string and push to stack
	L.Push(lua.LString(hex.EncodeToString(hash.Sum(nil))))

	return 1
}

// Md5Hash returns the md5 hash of the given string
func Md5Hash(L *lua.LState) int {
	// Get string to be hashed
	str := L.Get(2)

	// Check for valid string type
	if str.Type() != lua.LTString {

		L.ArgError(1, "Invalid string format. Expected string")
		return 0
	}

	// Hash string using sha1
	data := md5.Sum([]byte(str.String()))

	// Convert byte array to string and push tu stack
	L.Push(
		lua.LString(
			fmt.Sprintf("%x", data),
		),
	)

	return 1
}

// RandomString generates a random string with the given length
func RandomString(L *lua.LState) int {
	// Get length
	length := L.Get(2)

	// Valid length type
	if length.Type() != lua.LTNumber {
		L.ArgError(1, "Invalid length type. Expected number")
		return 0
	}

	// Push random string
	L.Push(lua.LString(uniuri.NewLenChars(L.ToInt(2), []byte("abcdefghijklmnropqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"))))

	return 1
}

// GenerateAuthSecretKey generates a valid authentication secret key
func GenerateAuthSecretKey(L *lua.LState) int {
	// Push random key
	L.Push(lua.LString(uniuri.NewLenChars(16, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"))))

	return 1
}

// GenerateQRCode generates a QR code for the given string and returns a base64 encoded image
func GenerateQRCode(L *lua.LState) int {
	// Get string to encode
	msg := L.ToString(2)

	// Create QR code
	code, err := qrcode.Encode(msg, qrcode.Medium, 256)

	if err != nil {
		L.RaiseError("Cannot create QR code: %v", err)
		return 0
	}

	// Base64 encode the byte array
	encoded := base64.StdEncoding.EncodeToString(code)

	// Push as string
	L.Push(lua.LString(encoded))

	return 1
}
