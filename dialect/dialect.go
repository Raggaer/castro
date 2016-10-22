package dialect

import (
	"fmt"
	"sync"
)

// Dialect interface used to define new
// database dialects
type Dialect interface {
	// Start function called when a dialect is choosen
	// useful to load stuff
	Start() error

	// Name returns the dialect name
	Name() string

	// Version returns the dialect version
	Version() string

	// LoadStages parses the stages xml file
	LoadStages() error

	// GetStages returns the server stages
	GetStages() []Stage
}

// Map struct used to save all registered dialects
// to later load the desired one
type Map struct {
	rw *sync.RWMutex
	m  map[string]Dialect
}

// Stage struct used for server stages
type Stage struct {
	From       int
	To         int
	Multiplier int
}

var (
	// Current holds the runtime dialect
	Current Dialect

	// List holds all registered dialects
	List = Map{
		rw: &sync.RWMutex{},
		m:  map[string]Dialect{},
	}
)

// Register saves a dialect into the map
func (m Map) Register(name string, d Dialect) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.m[name] = d
}

// Get returns the given dialect
func (m Map) Get(name string) (Dialect, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	if d, ok := m.m[name]; ok {
		return d, nil
	}
	return nil, fmt.Errorf("Dialect %v not found", name)
}

// SetDialect defines the dialect to use during
// runtime
func SetDialect(d Dialect) {
	Current = d
}
