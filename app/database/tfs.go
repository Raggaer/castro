package database

// TFS struct that defines a dialect for
// the forgotten server
type TFS struct {
}

// Name returns the current dialect name
func (t TFS) Name() string {
	return "the forgotten server"
}
