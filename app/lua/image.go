package lua

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/raggaer/goimage"
	"github.com/yuin/gopher-lua"
)

// SetImageMetaTable sets the image metatable of the given state
func SetImageMetaTable(luaState *lua.LState) {
	// Create and set the json metatable
	imgMetaTable := luaState.NewTypeMetatable(ImageMetaTableName)
	luaState.SetGlobal(ImageMetaTableName, imgMetaTable)

	// Set all image metatable functions
	luaState.SetFuncs(imgMetaTable, imgMethods)
}

// getGoImage retrieves the user data goimage from the given state
func getGoImage(luaState *lua.LState) goimage.Image {
	// Get metatable
	meta := luaState.Get(1)

	// Get user data
	data := luaState.GetField(meta, "__img").(*lua.LUserData)

	return data.Value.(goimage.Image)
}

// NewGoImage creates and returns a new goimage image
func NewGoImage(L *lua.LState) int {
	// Create image
	img := goimage.NewImage(
		L.ToInt(2),
		L.ToInt(3),
	)

	// Create metatable
	tbl := L.NewTable()

	// Create image user data
	imgUserData := L.NewUserData()
	imgUserData.Value = img

	// Set the user data field
	L.SetField(tbl, "__img", imgUserData)

	// Set the metatable methods
	L.SetFuncs(tbl, goimageMethods)

	// Push metatable
	L.Push(tbl)

	return 1
}

// WriteGoImageText writes text to the given goimage
func WriteGoImageText(L *lua.LState) int {
	// Get goimage
	img := getGoImage(L)

	// Convert color string to go color
	textColor, err := colorful.Hex(L.ToString(3))

	if err != nil {
		L.RaiseError("Cannot convert string to color: %v", err)
		return 0
	}

	// Write text to the image
	if err := img.WriteText(
		L.ToString(2),
		textColor,
		float64(L.ToInt(4)),
		L.ToInt(5),
		L.ToInt(6),
	); err != nil {
		L.RaiseError("Cannot write string to image: %v", err)
		return 0
	}

	return 0
}

// SetBackgroundGoImage sets the background of a goimage
func SetBackgroundGoImage(L *lua.LState) int {
	// Get goimage
	img := getGoImage(L)

	if err := img.SetBackGroundImage(L.ToString(2)); err != nil {
		L.RaiseError("Cannot set background image: %v", err)
		return 0
	}

	return 1
}

// SaveGoImage saves the given goimage
func SaveGoImage(L *lua.LState) int {
	// Get goimage
	img := getGoImage(L)

	if err := img.Save(L.ToString(2)); err != nil {
		L.RaiseError("Invalid image save location: %v", err)
		return 0
	}

	return 0
}
