package golinters

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/tools/go/loader"

	"github.com/skratchdot/open-golang/open"

	"github.com/thomasheller/golinters/gometalinter"
	"github.com/thomasheller/golinters/repo"
)

var (
	linters          []linter
	gometalinterDefs []string
	metalintPkgs     []string
)

type result struct {
	Name         string
	Repo         *repo.Repository
	GoParser     bool
	GoLoader     bool
	GoSSA        bool
	Gometalinter bool
	Metalint     bool
	Checker      bool
	Notes        string
}

func Analyze(ghUser *string, ghToken *string, out *string) {
	linters = list()

	a := &repo.GitHubAuth{*ghUser, *ghToken}

	log.Println("Fetching missing linters, if required...")

	for _, linter := range linters {
		log.Println(linter.name)

		if err := install(linter.path); err != nil {
			log.Printf("Error installing %s: %v\n", linter.name, err)
		}
	}

	if err := install("github.com/alecthomas/gometalinter"); err != nil {
		log.Printf("Error installing gometalinter: %v\n", err)
	}

	var err error

	g := &gometalinter.GometalinterSource{}
	gometalinterDefs, err = g.GetLinterDefinitions()
	if err != nil {
		log.Fatalf("Error finding gometalinter's linter definitions: %v")
	}

	metalintPkgs, err = imports("github.com/mvdan/lint/cmd/metalint")
	if err != nil {
		log.Fatalf("Couldn't find imports of metalint: %v")
	}

	var results []result

	for _, linter := range linters {
		r, err := details(linter, a)
		if err != nil {
			log.Printf("Error analzying %s: %v\n", linter.name, err)
			continue
		}
		results = append(results, r)
	}

	writeHTML(*out, results)
}

// install downloads a package through go get.
func install(path string) error {
	c := exec.Command("go", "get", "-d", path)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		return err
	}

	return nil
}

// details reports a linter's metadata, requirements and capabilities
// based on its package path, imports and GitHub API data.
func details(l linter, a *repo.GitHubAuth) (result, error) {
	log.Printf("Analyzing %s...", l.name)

	var pkgs []string
	var err error
	if pkgs, err = imports(l.path); err != nil {
		return result{}, err
	}

	var r result

	r.Name = l.name
	r.Repo, err = repo.Info(l.path, a)
	if err != nil {
		log.Printf("%s: could not get repository info: %v", l.name, err)
	}
	r.Notes = l.comment

	for _, pkg := range pkgs {
		switch pkg {
		case "go/parser":
			r.GoParser = true
		case "golang.org/x/tools/go/loader":
			r.GoLoader = true
		case "golang.org/x/tools/go/ssa":
			r.GoSSA = true
		case "github.com/mvdan/lint":
			r.Checker = true
		}
	}

	for _, def := range gometalinterDefs {
		if strings.HasPrefix(def, l.cmd) {
			r.Gometalinter = true
			break
		}
	}

	for _, metalintPkg := range metalintPkgs {
		if l.path == metalintPkg {
			r.Metalint = true
			break
		}
	}

	return r, nil
}

// imports returns all imports for the given package.
func imports(path string) ([]string, error) {
	args := []string{path}

	var conf loader.Config
	if _, err := conf.FromArgs(args, false); err != nil {
		return nil, err
	}

	var lprog *loader.Program
	var err error
	if lprog, err = conf.Load(); err != nil {
		return nil, err
	}

	var imports []string
	for p := range lprog.AllPackages {
		imports = append(imports, p.Path())
	}

	return imports, nil
}

// writeHTML generates a HTML report and writes it to a file. If no
// filename is given, a temporary file is chosen and the report opens
// in the default browser.
func writeHTML(file string, results []result) error {
	browser := file == ""

	tmpl := template.Must(template.New("html").Parse(htmlTemplate))

	if browser {
		dir, err := ioutil.TempDir("", "golinters")
		if err != nil {
			return err
		}

		file = filepath.Join(dir, "golinters.html")
	}

	out, err := os.Create(file)
	if err != nil {
		return err
	}

	defer out.Close()

	data := TemplateData{
		Timestamp: time.Now().Format(time.RFC1123),
		Results:   results,
	}

	err = tmpl.Execute(out, data)
	if err != nil {
		return err
	}

	if browser {
		open.Run("file://" + out.Name())
	}

	return nil
}

type TemplateData struct {
	Timestamp string
	Results   []result
}

const htmlTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<style>
			html, body {
				font-family: Arial, sans-serif;
			}
			tt {
				font-family: Menlo, monospace;
			}
			table, th, td {
				border: 1px solid #000;
				border-collapse: collapse;
			}
			th, td {
				padding: .33em;
			}
			td.t, td.f {
				text-align: center;
			}
			td.t {
				background-color: #5bd64a;
			}
			td.f {
				background-color: #d64a4a;
			}
			td.notes, .timestamp {
				font-size: small;
			}
		</style>
	</head>
	<body>
		<table>
			<thead>
				<tr>
					<th colspan="3">General info</th>
					<th colspan="3">Input</th>
					<th colspan="3">Metalinter support</th>
					<th rowspan="2">Notes</th>
				</tr>
				<tr>
					<th>Name</th>
					<th>Maintainer</th>
					<th>Repository URL</th>
					<th><tt>go/parser</tt></th>
					<th><tt>go/loader</tt></th>
					<th><tt>go/ssa</tt></th>
					<th><tt>gometalinter</tt></th>
					<th><tt>metalint</tt></th>
					<th><tt>Checker</tt></th>
				</tr>
			</thead>
			<tbody>
				{{ range .Results }}<tr>
					<td>{{ .Name }}</td>
					<td>{{ if .Repo }}{{ .Repo.Maintainer }}{{ end }}</td>
					<td>{{ if .Repo }}<a href="{{ .Repo.URL }}">{{ .Repo.URL }}</a>{{ end }}</td>
					{{ if .GoParser }}<td class="t">Y</td>{{ else }}<td class="f">N</td>{{ end }}
					{{ if .GoLoader }}<td class="t">Y</td>{{ else }}<td class="f">N</td>{{ end }}
					{{ if .GoSSA }}<td class="t">Y</td>{{ else }}<td class="f">N</td>{{ end }}
					{{ if .Gometalinter }}<td class="t">Y</td>{{ else }}<td class="f">N</td>{{ end }}
					{{ if .Metalint }}<td class="t">Y</td>{{ else }}<td class="f">N</td>{{ end }}
					{{ if .Checker }}<td class="t">Y</td>{{ else }}<td class="f">N</td>{{ end }}
					<td class="notes">{{ .Notes }}</td>
				</tr>{{ end }}
			</tbody>
		</table>
		<p class="timestamp">{{ .Timestamp }}</p>
	</body>
</html>`
