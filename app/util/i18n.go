package util

import (
	"os"
	"path/filepath"
	"strings"

	"sync"

	"github.com/BurntSushi/toml"
)

// LanguageFiles global language holder
var LanguageFiles = &LanguageHolder{
	List: map[string]*Language{},
	rw:   sync.RWMutex{},
}

type LanguageHolder struct {
	rw   sync.RWMutex
	List map[string]*Language
}

// Language represents a i18n language file
type Language struct {
	Name string
	Data map[string]string
}

func Loadi18n(path string) error {
	// Lock language files mutex
	LanguageFiles.rw.Lock()
	defer LanguageFiles.rw.Unlock()

	// Walk over i18n directory
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// Check if valid language file
		if !strings.HasSuffix(info.Name(), ".i18n") || info.IsDir() {
			return nil
		}

		// Decode language file
		langData := map[string]string{}
		if _, err := toml.DecodeFile(path, &langData); err != nil {
			return err
		}

		// Append language file
		LanguageFiles.List[strings.TrimSuffix(info.Name(), ".i18n")] = &Language{
			Name: strings.TrimSuffix(info.Name(), ".i18n"),
			Data: langData,
		}

		return nil
	})
}

// Get retrieves the given language
func (l *LanguageHolder) Get(lang string) (*Language, bool) {
	// Lock mutex
	l.rw.RLock()
	defer l.rw.RUnlock()

	// Load language from the list
	ng, ok := l.List[lang]
	if !ok {
		return nil, false
	}

	return ng, true
}
