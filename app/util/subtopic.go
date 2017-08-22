package util

import (
	"os"
	"path/filepath"
	"strings"
)

// GetLuaFiles walks the given path and returns all lua files
func GetLuaFiles(dir string) ([]string, error) {
	// Data holder
	list := []string{}

	// Walk over the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		// Check if file has .lua extension
		if strings.HasSuffix(info.Name(), ".lua") && info.Name() != "config.lua" {

			// Append file
			list = append(list, filepath.Join(
				strings.ToLower(path),
			))
		}

		return nil
	})

	return list, err
}
