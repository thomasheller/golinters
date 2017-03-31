package gometalinter

import (
	"errors"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/loader"
)

// GometalinterAST attemts to find the current linter definitions in
// the gometalinter source. It spends too much time on this task and
// is very picky about the structure of gometalinter's source code.
// You don't want to use GometalinterAST unless you care to know how
// to dig through a particular AST manually. Note that there exists
// ast.Walk if you're OK with depth-first search.
type GometalinterAST struct {
	defs []string
}

func (g *GometalinterAST) GetLinterDefinitions() ([]string, error) {
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

	found := g.parseProg(lprog)

	if !found {
		return nil, errors.New("linter definitions not found")
	}

	return g.defs, nil
}

func (g *GometalinterAST) parseProg(lprog *loader.Program) bool {
	for _, pkg := range lprog.InitialPackages() {
		if g.parseFiles(pkg.Files) {
			return true
		}
	}
	return false
}

func (g *GometalinterAST) parseFiles(files []*ast.File) bool {
	for _, file := range files {
		if g.parseDecls(file.Decls) {
			return true
		}
	}
	return false
}

func (g *GometalinterAST) parseDecls(decls []ast.Decl) bool {
	for _, decl := range decls {
		switch t := decl.(type) {
		case *ast.GenDecl:
			if g.parseGenDecl(t) {
				return true
			}
		}
	}
	return false
}

func (g *GometalinterAST) parseGenDecl(decl *ast.GenDecl) bool {
	if decl.Tok == token.VAR {
		return g.parseSpecs(decl.Specs)
	}
	return false
}

func (g *GometalinterAST) parseSpecs(specs []ast.Spec) bool {
	for _, spec := range specs {
		switch t := spec.(type) {
		case *ast.ValueSpec:
			if g.parseValueSpec(t) {
				return true
			}
		}
	}
	return false
}

func (g *GometalinterAST) parseValueSpec(spec *ast.ValueSpec) bool {
	for _, ident := range spec.Names {
		if ident.Name == "linterDefinitions" {
			g.parseValues(spec.Values)
			return true
		}
	}
	return false
}

func (g *GometalinterAST) parseValues(values []ast.Expr) {
	for _, value := range values {
		switch t := value.(type) {
		case *ast.CompositeLit:
			g.parseElts(t.Elts)
		}
	}
}

func (g *GometalinterAST) parseElts(elts []ast.Expr) {
	g.defs = []string{}
	for _, elt := range elts {
		switch t := elt.(type) {
		case *ast.KeyValueExpr:
			found, s := g.parseKeyValue(t.Value)
			if found {
				g.defs = append(g.defs, s)
			}
		}
	}
}

func (g *GometalinterAST) parseKeyValue(value ast.Expr) (bool, string) {
	switch t := value.(type) {
	case *ast.BasicLit:
		return true, t.Value[1 : len(t.Value)-2]
	case *ast.BinaryExpr:
		return g.parseBinary(t)
	}
	return false, ""
}

func (g *GometalinterAST) parseBinary(expr *ast.BinaryExpr) (bool, string) {
	// handles the "go vet" definitions
	return g.parseKeyValue(expr.X)
}
