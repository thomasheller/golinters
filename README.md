# golinters

golinters generates HTML reports about Go linters. It fetches the
linters' source and does static analysis. It also queries the GitHub
API to figure out the maintainers.

```sh
$ go get github.com/thomasheller/golinters/cmd/golinters
$ golinters
[...]
```

## Example output

![HTML screenshot](https://raw.githubusercontent.com/thomasheller/golinters/master/examples/output-2017-03-31-214655-CEST.png)
