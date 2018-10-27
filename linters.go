package golinters

type linter struct {
	name    string
	cmd     string
	path    string
	comment string
}

func list() []linter {
	return []linter{
		{"aligncheck", "aligncheck", "github.com/opennota/check/cmd/aligncheck", ""},
		{"deadcode", "deadcode", "github.com/tsenart/deadcode", ""},
		{"dupl", "dupl", "github.com/mibk/dupl", ""},
		{"errcheck", "errcheck", "github.com/kisielk/errcheck", ""},
		{"gas", "gas", "github.com/GoASTScanner/gas", ""},
		{"goconst", "goconst", "github.com/jgautheron/goconst/cmd/goconst", ""},
		{"gocyclo", "gocyclo", "github.com/fzipp/gocyclo", "gometalinter uses a fork: github.com/alecthomas/gocyclo"},
		{"gofmt", "gofmt -l -s", "github.com/golang/go/src/cmd/gofmt", ""},
		{"goimports", "goimports", "golang.org/x/tools/cmd/goimports", ""},
		{"golint", "golint", "golang.org/x/lint/golint", ""},
		{"gosimple", "gosimple", "honnef.co/go/tools/cmd/gosimple", ""},
		{"gotype", "gotype", "golang.org/x/tools/cmd/gotype", ""},
		{"ineffassign", "ineffassign", "github.com/gordonklaus/ineffassign", ""},
		{"interfacer", "interfacer", "github.com/mvdan/interfacer/cmd/interfacer", ""},
		{"lll", "lll", "github.com/walle/lll/cmd/lll", ""},
		{"misspell", "misspell", "github.com/client9/misspell/cmd/misspell", ""},
		{"safesql", "safesql", "github.com/stripe/safesql", ""},
		{"staticcheck", "staticcheck", "honnef.co/go/tools/cmd/staticcheck", ""},
		{"structcheck", "structcheck", "github.com/opennota/check/cmd/structcheck", ""},
		// {"test", "go test {path}:^--- FAIL:", "github.com/golang/go/src/cmd/go/internal/test", ""}, // TODO
		// {"testify", "go test {path}:Location:", "github.com/golang/go/src/cmd/go/internal/test", "essentially the same as test, gometalinter parses the output differently"}, // TODO
		{"unconvert", "unconvert", "github.com/mdempsky/unconvert", ""},
		{"unparam", "unparam", "github.com/mvdan/unparam", ""},
		{"unused", "unused", "honnef.co/go/tools/cmd/unused", ""},
		{"varcheck", "varcheck", "github.com/opennota/check/cmd/varcheck", ""},
		{"vet", "go tool vet {path}", "github.com/golang/go/src/cmd/vet", ""}, // include "{path}" so it's different from "--shadow"
		{"vetshadow", "go tool vet --shadow", "github.com/golang/go/src/cmd/vet", "same linter as vet, just run with --shadow"},
	}
}
