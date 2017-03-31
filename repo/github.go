package repo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bndr/gopencils"
)

// GitHubAuth represents authentication data for the GitHub API.
type GitHubAuth struct {
	// Username is a GitHub username.
	Username string
	// Token is a GitHub access token.
	Token string
}

type repo struct {
	Owner    owner
	HTML_URL string
}

type owner struct {
	Login string
}

type user struct {
	Login string
	Name  string
}

// Repo fetches basic metadata of a GitHub repository identified by
// its import paths. Import paths that are not recognizable GitHub
// repositories return an error.
func GitHub(path string, githubAuth *GitHubAuth) (*Repository, error) {
	if !strings.HasPrefix(path, "github.com") {
		return nil, errors.New("not a GitHub repository")
	}

	var auth *gopencils.BasicAuth

	if githubAuth.Username != "" && githubAuth.Token != "" {
		auth = &gopencils.BasicAuth{githubAuth.Username, githubAuth.Token}
	}

	api := gopencils.Api("https://api.github.com", auth)

	repos := api.Res("repos")
	users := api.Res("users")

	parts := strings.Split(path, "/")
	repoName := parts[1] + "/" + parts[2]

	var res *gopencils.Resource
	var err error

	r := new(repo)

	if res, err = repos.Id(repoName, r).Get(); err != nil {
		return nil, err
	}
	if res.Raw.StatusCode >= 400 {
		return nil, errorMsg(res.Raw.StatusCode)
	}

	username := r.Owner.Login

	u := new(user)

	if res, err = users.Id(username, u).Get(); err != nil {
		return nil, err
	}
	if res.Raw.StatusCode >= 400 {
		return nil, errorMsg(res.Raw.StatusCode)
	}

	result := &Repository{u.Name, r.HTML_URL}

	if result.Maintainer == "" {
		result.Maintainer = u.Login
	}

	return result, nil
}

func errorMsg(statusCode int) error {
	if statusCode == 403 {
		return errors.New("Error 403 - possibly rate limit exceeded. Did you supply GitHub credentials?")
	} else {
		return fmt.Errorf("Error %d", statusCode)
	}
}
