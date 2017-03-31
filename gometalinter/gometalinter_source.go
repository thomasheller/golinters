package gometalinter

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/thomasheller/gopath"
)

// GometalinterSource parses the gometalinter source (plain text) to
// find the linter definitions. This depends on no major changes to
// gometalinter's source. For the sake of speed and simplicity, the
// returned strings are inaccurate, but this is sufficient for what
// golinters wants to check.
type GometalinterSource struct {
	s    *bufio.Scanner
	defs []string
}

type stateFn func(*GometalinterSource) (stateFn, error)

func (g *GometalinterSource) GetLinterDefinitions() ([]string, error) {
	path, err := gopath.Join("src/github.com/alecthomas/gometalinter/config.go")
	if err != nil {
		return nil, err
	}

	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	g.s = bufio.NewScanner(r)

	var state stateFn
	for state, err = parseBefore, nil; state != nil; {
		state, err = state(g)
	}

	return g.defs, err
}

func parseBefore(g *GometalinterSource) (stateFn, error) {
	for g.s.Scan() {
		t := g.s.Text()

		if strings.HasPrefix(t, "\tlinterDefinitions = map[string]string{") {
			return parseDefs, nil
		}
	}

	return nil, errors.New("parse error: linterDefinitions not found")
}

func parseDefs(g *GometalinterSource) (stateFn, error) {
	for g.s.Scan() {
		t := g.s.Text()

		if t == "\t}" {
			return nil, nil // done
		}

		def, err := g.parseDef(t)
		if err != nil {
			return nil, err
		}

		g.defs = append(g.defs, def)
	}

	return nil, errors.New("parse error: unexpected EOF")
}

func (g *GometalinterSource) parseDef(line string) (string, error) {
	quote := g.findNthQuote(line, 3)
	if quote == -1 {
		return "", errors.New("parse error: unexpected quotes in linterDefinitions")
	}

	def := line[quote+1 : len(line)-2]

	return def, nil
}

func (g *GometalinterSource) findNthQuote(line string, n int) int {
	num := 0
	for i, c := range line {
		if c == '`' || c == '"' {
			num++
		}
		if num == n {
			return i
		}
	}
	return -1
}
