package lua

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/gopher-lua"
)

// SetFileMetaTable sets the file metatable of the given state
func SetFileMetaTable(luaState *lua.LState) {
	// Create and set the file metatable
	fileMetaTable := luaState.NewTypeMetatable(FileMetaTableName)
	luaState.SetGlobal(FileMetaTableName, fileMetaTable)

	// Set all file metatable functions
	luaState.SetFuncs(fileMetaTable, fileMethods)
}

// CheckFileExists checks if the given file exists
func CheckFileExists(L *lua.LState) int {
	// Get file info
	_, err := os.Stat(L.ToString(2))

	// If file exists push true
	if err == nil {
		L.Push(lua.LBool(true))
		return 1
	}

	// Check if file does not exists then push false
	if os.IsNotExist(err) {
		L.Push(lua.LBool(false))
		return 1
	}

	L.RaiseError("Cannot get file information: %v", err)

	return 0
}

// GetFileModTime gets the file unix timestamp
func GetFileModTime(L *lua.LState) int {
	// Get file info
	info, err := os.Stat(L.ToString(2))

	if err != nil {
		L.RaiseError("Cannot get file information: %v", err)
	}

	// Push file mod time
	L.Push(lua.LNumber(info.ModTime().Unix()))

	return 1
}

// GetDirectories gets any directories in the provided path
func GetDirectories(L *lua.LState) int {
	// Get files
	files, err := ioutil.ReadDir(L.ToString(2))

	if err != nil {
		L.Push(lua.LNil)
		return 1
	}

	// Result table
	tbl := L.NewTable()

	for _, f := range files {
		if f.IsDir() {
			// Append directory name
			tbl.Append(lua.LString(f.Name()))
		}
	}

	// Push directory list
	L.Push(tbl)

	return 1
}

// GetFiles gets a list of files for the given directory
func GetFiles(L *lua.LState) int {
	// Get files
	files, err := ioutil.ReadDir(L.ToString(2))

	if err != nil {
		L.Push(lua.LNil)
		return 1
	}

	// Result table
	tbl := L.NewTable()

	for _, f := range files {
		if !f.IsDir() {
			// Append file name
			tbl.Append(lua.LString(f.Name()))
		}
	}

	// Push file list
	L.Push(tbl)

	return 1
}

// CreateDirectory creates the given directory
func CreateDirectory(L *lua.LState) int {
	// Get directory path
	path := L.ToString(2)

	// Create any missing directory
	if err := os.MkdirAll(path, os.ModeDir); err != nil {
		L.RaiseError("Cannot create missing directories: %v", err)
	}

	return 0
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file to the given directory
func UnzipFile(L *lua.LState) int {

	var filenames []string
	dest := L.ToString(3)

	// Open zip archive
	r, err := zip.OpenReader(L.ToString(2))
	if err != nil {
		L.RaiseError("Cannot open file to unzip: %v", err)
	}

	defer r.Close()

	for _, f := range r.File {

		// Store filename/path
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. https://snyk.io/research/zip-slip-vulnerability#go
		if !strings.HasPrefix(fpath, filepath.Clean(dest) + string(os.PathSeparator)) {
			// Push nil + error
			L.Push(lua.LNil)
			L.Push(lua.LString("Illegal file path"))
			return 2
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Create folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create directories
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			// Push nil + error
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		// Open output file in write mode
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			// Push nil + error
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		// Open source file
		rc, err := f.Open()
		if err != nil {
			// Push nil + error
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		// Copy file contents
		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			// Push nil + error
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	}

	// Push success
	L.Push(lua.LBool(true))

	return 1
}
