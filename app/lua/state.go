package lua

import (
	"github.com/raggaer/castro/app/util"
	glua "github.com/yuin/gopher-lua"
	"sync"
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

		s.List[subtopic] = append(s.List[subtopic], state)
	}

	return nil
}

func (s *stateList) Get(path string) (*glua.LState, error) {
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
	// Lock mutex
	s.rw.Lock()
	defer s.rw.Unlock()

	// Save state
	s.List[path] = append(s.List[path], state)
}
