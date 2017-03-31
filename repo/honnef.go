package repo

import (
	"errors"
	"strings"
)

const honnefToolsBase = "honnef.co/go/tools/"

func Honnef(url string) (*Repository, error) {
	if !strings.HasPrefix(url, honnefToolsBase) {
		return nil, errors.New("don't know about this URL")
	}

	return &Repository{
		"Dominik Honnef",
		"https://github.com/dominikh/go-tools",
	}, nil
}
