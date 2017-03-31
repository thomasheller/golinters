package gometalinter

import (
	"errors"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa/ssautil"
)

type GometalinterSSA struct {
	defs []string
}

func (g *GometalinterSSA) GetLinterDefinitions() ([]string, error) {
	args := []string{"github.com/alecthomas/gometalinter"}

	var conf loader.Config
	if _, err := conf.FromArgs(args, false); err != nil {
		return nil, err
	}

	var lprog *loader.Program
	var err error
	if lprog, err = conf.Load(); err != nil {
		return nil, err
	}

	found := g.parseSSA(lprog)

	if !found {
		return nil, errors.New("linter definitions not found")
	}

	return g.defs, nil
}

func (g *GometalinterSSA) parseSSA(lprog *loader.Program) bool {
	prog := ssautil.CreateProgram(lprog, 0)
	prog.Build()

	for _, pkg := range ssautil.MainPackages(prog.AllPackages()) {
		pkg.Var("linterDefinitions")

		// TODO: How to extract the actual string value?
	}

	return false
}
