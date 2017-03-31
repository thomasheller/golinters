# golinters

golinters generates HTML reports about Go linters. It fetches the
linters' source and does static analysis. It also queries the GitHub
API to figure out the maintainers.

```sh
$ go get github.com/thomasheller/golinters/cmd/golinters
$ golinters
[...]
```

If a popular linter is missing, please file an issue!

## Example output

![HTML screenshot](https://raw.githubusercontent.com/thomasheller/golinters/master/examples/output-2017-03-31-214655-CEST.png)

[Download corresponding HTML file](https://raw.githubusercontent.com/thomasheller/golinters/master/examples/output-2017-03-31-214655-CEST.html)

## Options

By default, golinters will create a temporary file and open the report
in the standard browser.

You can specify `-write somefile.html` though, if you want golinters
to just write to a specific file and not open any browser.

Because golinters uses the GitHub API to figure out the maintainers'
names, you might want to supply a GitHub username and API token via
`-ghuser` and `-ghtoken` so that you don't run into rate limit
problems.

If you want to start over, you can use `-remove` to delete the
linters' source in your GOPATH. Be careful, as this deletes entire
repositories, even if the linter is just one part of it.
