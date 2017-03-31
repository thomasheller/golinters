package main

import (
	"flag"

	"github.com/thomasheller/golinters"
)

func main() {
	out := flag.String("write", "", "write HTML output to file instead of opening a browser")
	ghUser := flag.String("ghuser", "", "GitHub username (for API use)")
	ghToken := flag.String("ghtoken", "", "GitHub token (for API use)")
	remove := flag.Bool("remove", false, "delete all linters in GOPATH/src (be careful)")
	flag.Parse()

	if *remove {
		golinters.RemoveAllRepos()
	} else {
		golinters.Analyze(ghUser, ghToken, out)
	}
}
