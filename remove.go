package golinters

import (
	"log"
	"os"
	"strings"

	"github.com/thomasheller/gopath"
)

// RemoveAllRepos deletes the entire source code repository in
// GOPATH/src for all known linters. Be careful.
func RemoveAllRepos() {
	for _, linter := range list() {
		parts := strings.Split(linter.path, "/")

		p, err := gopath.Join("src", parts[0], parts[1], parts[2])
		if err != nil {
			log.Printf("Error getting path to linter %s: %v", linter.name, err)
			continue
		}

		log.Printf("removing %s\n", p)

		err = os.RemoveAll(p)
		if err != nil {
			log.Printf("Error removing linter %s: %v", linter.name, err)
		}
	}
}
