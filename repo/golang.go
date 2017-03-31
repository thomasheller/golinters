package repo

import (
	"errors"
	"strings"
)

const golangToolsBase = "golang.org/x/tools/"

func Golang(url string) (*Repository, error) {
	if !strings.HasPrefix(url, golangToolsBase) {
		return nil, errors.New("don't know about this URL")
	}

	return &Repository{
		"Go",
		"https://github.com/golang/tools",
	}, nil
}
