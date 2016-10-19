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

	// GetStages returns the server stages
	GetStages() []Stage
}

// Stage struct used for server stages
type Stage struct {
	From       int
	To         int
	Multiplier int
}

// Current holds the runtime dialect
var Current Dialect

// SetDialect defines the dialect to use during
// runtime
func SetDialect(d Dialect) {
	Current = d
}
