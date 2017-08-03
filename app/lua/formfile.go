package lua

import (
	"bytes"
	"github.com/kardianos/osext"
	"github.com/nfnt/resize"
	"github.com/yuin/gopher-lua"
	"image/png"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type formFileUserData struct {
	File   []byte
	Header *multipart.FileHeader
}

func createFormFileMetaTable(file []byte, header *multipart.FileHeader, L *lua.LState) *lua.LTable {
	// Create metatable
	table := L.NewTable()

	// Set metatable functions
	L.SetFuncs(table, formFileMethods)

	// Create new user data to hold file and header information
	u := L.NewUserData()

	// Set user value
	u.Value = &formFileUserData{
		File:   file,
		Header: header,
	}

	// Set user data as field
	L.SetField(table, "__file", u)

	return table
}

func getFormFileObject(L *lua.LState) *formFileUserData {
	// Retrieve metatable
	table := L.ToTable(1)

	// Get user data field
	data := L.GetField(table, "__file").(*lua.LUserData)

	// Return user value as form file pointer
	return data.Value.(*formFileUserData)
}

// FormFileIsValidPNG checks if the current form file is a valid png file
func FormFileIsValidPNG(L *lua.LState) int {
	// Get form file
	formFile := getFormFileObject(L)

	if http.DetectContentType(formFile.File) != "image/png" {

		// Push false
		L.Push(lua.LBool(false))

		return 1
	}

	// The file is a image png
	L.Push(lua.LBool(true))

	return 1
}

// FormFileIsValidExtension checks if the current form file is a valid extension
func FormFileIsValidExtension(L *lua.LState) int {
	// Get form file
	formFile := getFormFileObject(L)

	if http.DetectContentType(formFile.File) != L.ToString(2) {

		// Push false
		L.Push(lua.LBool(false))

		return 1
	}

	// The file is a image png
	L.Push(lua.LBool(true))

	return 1
}

// GetFormFileByteArray returns the form file byte array as a lua string
func GetFormFileByteArray(L *lua.LState) int {
	// Get form file
	formFile := getFormFileObject(L)

	// Push file as string
	L.Push(lua.LString(formFile.File))

	return 1
}

// SaveFormFile saves the current form file to the given destination
func SaveFormFile(L *lua.LState) int {
	// Get form file
	formFile := getFormFileObject(L)

	// Get executable folder
	f, err := osext.ExecutableFolder()

	if err != nil {
		L.RaiseError("Cannot get executable folder path: %v", err)
		return 0
	}

	// Get destination
	destination := L.ToString(2)

	// Write file to handle
	if err := ioutil.WriteFile(filepath.Join(f, destination), formFile.File, 0666); err != nil {
		L.RaiseError("Cannot save file to destination: %v", err)
	}

	return 0
}

// SaveFormFileAsPNG saves the current form file as a png with resizing optional
func SaveFormFileAsPNG(L *lua.LState) int {
	// Get form file
	formFile := getFormFileObject(L)

	// Get executable folder
	f, err := osext.ExecutableFolder()

	if err != nil {
		L.RaiseError("Cannot get executable folder path: %v", err)
		return 0
	}

	// Get destination
	destination := L.ToString(2)

	// Create file handle
	file, err := os.OpenFile(filepath.Join(f, destination), os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		L.RaiseError("Cannot get file handle: %v", err)
		return 0
	}

	// Close file handle
	defer file.Close()

	// Create png image from byte array
	pngImage, err := png.Decode(bytes.NewBuffer(formFile.File))

	if err != nil {
		L.RaiseError("Cannot decode image from byte array: %v", err)
		return 0
	}

	// Get desired image sizes
	imageWidth := L.ToInt(3)
	imageHeight := L.ToInt(4)

	// Resize if wanted using nearest neighbor algorithm
	if imageWidth > 0 && imageHeight > 0 {
		pngImage = resize.Resize(uint(imageWidth), uint(imageHeight), pngImage, resize.NearestNeighbor)
	}

	// Encode image to file handle
	if err := png.Encode(file, pngImage); err != nil {
		L.RaiseError("Cannot decode png image to handle: %v", err)
	}

	return 0
}
