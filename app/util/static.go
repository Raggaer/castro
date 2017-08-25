package util

import (
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/raggaer/castro/app/database"
)

var (
	// ExtensionStatic holds all extension subtopic static folders
	ExtensionStatic = &StaticList{
		list: map[string]http.FileSystem{},
	}
)

// StaticList struct used to hold static lists
type StaticList struct {
	rw   sync.RWMutex
	list map[string]http.FileSystem
}

// FileExists checks if the given resource exists
func (e *StaticList) FileExists(id string) (http.FileSystem, bool) {
	// Read lock mutex
	e.rw.RLock()
	defer e.rw.RUnlock()

	// Get element from the map
	dir, ok := e.list[id]

	if !ok {
		return nil, false
	}

	return dir, true
}

// Load loads all the static resources from the enabled extensions
func (e *StaticList) Load(d string) error {
	// Lock and unlock mutexes
	e.rw.Lock()
	defer e.rw.Unlock()

	// Reset static list
	e.list = map[string]http.FileSystem{}

	// Get extensions from database
	rows, err := database.DB.Queryx("SELECT id FROM castro_extensions WHERE installed = 1")

	if err != nil {
		return err
	}

	// Close rows
	defer rows.Close()

	// Loop rows
	for rows.Next() {

		// Hold extension identifier
		var id string

		// Scan extension id
		if err := rows.Scan(&id); err != nil {
			return err
		}

		dir := filepath.Join(d, id, "static")

		// Make sure that directory exist
		if _, err = os.Stat(dir); err != nil {

			if !os.IsNotExist(err) {
				return err
			}

			continue
		}

		// Add http directory to map
		e.list[id] = http.Dir(dir)
	}

	return nil
}
