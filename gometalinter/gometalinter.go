package gometalinter

type Gometalinter interface {
	// GetLinterDefinitions returns gometalinter's linter
	// definition strings.
	GetLinterDefinitions() ([]string, error)
}
