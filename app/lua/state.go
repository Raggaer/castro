package lua

import (
	"github.com/raggaer/castro/app/util"
	"github.com/raggaer/castro/app/database"
	glua "github.com/yuin/gopher-lua"
	"strings"
	"sync"
	"path/filepath"
	"os"
)

var (
	// PageList list of subtopic states
	PageList = &stateList{
		List: make(map[string][]*glua.LState),
		Type: "page",
	}

	// WidgetList list of widget states
	WidgetList = &stateList{
		List: make(map[string][]*glua.LState),
		Type: "widget",
	}
)

type stateList struct {
	rw   sync.Mutex
	List map[string][]*glua.LState
	Type string
}

// Load loads the given state list
func (s *stateList) Load(dir string) error {
	// Lock mutex
	s.rw.Lock()
	defer s.rw.Unlock()

	// Set list
	s.List = make(map[string][]*glua.LState)

	// Get subtopic list
	subtopicList, err := util.GetLuaFiles(dir)

	if err != nil {
		return err
	}

	// Loop subtopic list
	for _, subtopic := range subtopicList {

		// Create state
		state := glua.NewState()

		// Set castro metatables
		getApplicationState(state)

		if err := state.DoFile(subtopic); err != nil {
			return err
		}

		// Set lowercase path
		path := strings.ToLower(subtopic)

		// Add state to the pool
		s.List[path] = append(s.List[path], state)
	}

	return nil
}

// LoadExtensions loads the given state list
func (s *stateList) LoadExtensions() error {
	// Lock mutex
	s.rw.Lock()
	defer s.rw.Unlock()

	// Set list
	s.List = make(map[string][]*glua.LState)

	// Set extension type
	extType := s.Type + "s"

	// Get extensions from database
	rows, err := database.DB.Queryx(strings.Replace("SELECT extension_id FROM castro_extension_? WHERE enabled = 1", "?", extType, -1))

	if err != nil {
		return err
	}

	// Close rows
	defer rows.Close()

	// Loop rows
	for rows.Next() {

		// Hold extension id
		var extension_id string

		if err := rows.Scan(&extension_id); err != nil {
			return err
		}

		dir := filepath.Join("extensions", extension_id, extType)

		// Make sure that directory exist
		if _, err = os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				util.Logger.Logger.Errorf("Missing %v directory in extension %v", extType, extension_id)
			}
			continue
		}

		// Get subtopic list
		subtopicList, err := util.GetLuaFiles(dir)

		if err != nil {
			return err
		}

		// Loop subtopic list
		for _, subtopic := range subtopicList {

			// Create state
			state := glua.NewState()

			// Set castro metatables
			getApplicationState(state)

			if err := state.DoFile(subtopic); err != nil {
				return err
			}

			// Set lowercase path
			path := strings.ToLower(strings.Replace(subtopic, dir, extType, -1))

			// Add state to the pool
			s.List[path] = append(s.List[path], state)
		}
	}

	return nil
}

func (s *stateList) Get(path string) (*glua.LState, error) {
	// Set path as lowercase
	path = strings.ToLower(path)

	// Lock mutex
	s.rw.Lock()
	defer s.rw.Unlock()

	if len(s.List[path]) == 0 {

		// Create new state
		state := glua.NewState()

		// Set castro metatables
		getApplicationState(state)

		if err := state.DoFile(path); err != nil {
			return nil, err
		}

		return state, nil
	}

	// Return last state from the pool
	x := s.List[path][len(s.List[path])-1]
	s.List[path] = s.List[path][0 : len(s.List[path])-1]

	return x, nil
}

func (s *stateList) Put(state *glua.LState, path string) {
	// Set path as lowercase
	path = strings.ToLower(path)

	// Lock mutex
	s.rw.Lock()
	defer s.rw.Unlock()

	// Save state
	s.List[path] = append(s.List[path], state)
}
