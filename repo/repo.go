package repo

import (
	"errors"
	"strings"
)

// Repository describes GitHub repository metadata
type Repository struct {
	// Maintainer is the full name of the repository owner, or
	// username if the real name is unknown.
	Maintainer string
	// URL is the HTML URL of a repository that can be viewed in a
	// webbrowser.
	URL string
}

// Info returns information about source code repositories based on
// the import path. Only a few common paths are currently supported.
func Info(path string, gitHubAuth *GitHubAuth) (*Repository, error) {
	if strings.HasPrefix(path, "github.com/") {
		return GitHub(path, gitHubAuth)
	}

	if strings.HasPrefix(path, "honnef.co/") {
		return Honnef(path)
	}

	if strings.HasPrefix(path, "golang.org/") {
		return Golang(path)
	}

	return nil, errors.New("not a recognized repository path")
}
