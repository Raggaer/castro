package database

// Dialect interface used to define new
// database dialects
type Dialect interface {
	// Name returns the dialect name
	Name() string

	// Version returns the dialect version
	Version() string
}

// CurrentDialect holds the runtime dialect
var CurrentDialect Dialect

// SetDialect defines the dialect to use during
// runtime
func SetDialect(d Dialect) {
	CurrentDialect = d
}
