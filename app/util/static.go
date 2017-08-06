package util

import (
	"sync"
	"github.com/raggaer/castro/app/database"
	"path/filepath"
	"os"
	"strings"
)

var (
	// ExtensionStatic holds all extension subtopic static folders
	ExtensionStatic = &StaticList{
		list: []string{},
	}
)

// StaticList struct used to hold static lists
type StaticList struct {
	rw sync.RWMutex
	list []string
}

// FileExists checks if the given resource exists
func (e *StaticList) FileExists(url string) bool {
	// Read lock mutex
	e.rw.RLock()
	defer e.rw.RUnlock()

	// Loop list
	for _, e := range e.list {

		if strings.HasPrefix(strings.Replace(url, "/", "\\", -1), filepath.Join(e)) {
			return true
		}
	}

	return false
}

// Load loads all the static resources from the enabled extensions
func (e *StaticList) Load(d string) error {
	// Lock and unlock mutexes
	e.rw.Lock()
	defer e.rw.Unlock()

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

		e.list = append(e.list, strings.Replace(dir, d, "", 1))
	}

	return nil
}
