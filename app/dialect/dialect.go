package dialect

// Dialect interface used to define new
// database dialects
type Dialect interface {
	// Name returns the dialect name
	Name() string

	// Version returns the dialect version
	Version() string

	// LoadStages parses the stages xml file
	LoadStages() error
}

// Current holds the runtime dialect
var Current Dialect

// SetDialect defines the dialect to use during
// runtime
func SetDialect(d Dialect) {
	Current = d
}
