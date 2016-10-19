package tfs

// TFS struct that defines a dialect for
// the forgotten server
type TFS struct {
}

// Name shows the dialect name
func (t TFS) Name() string {
	return "the forgotten server"
}

// Version shows the dialect version
func (t TFS) Version() string {
	return "0.1 alpha-preview"
}

// LoadStages loads server xml stages
func (t TFS) LoadStages() error {
	return nil
}
