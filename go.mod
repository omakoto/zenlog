module github.com/omakoto/zenlog

go 1.13

//replace github.com/omakoto/go-common => ../go-common

// To update the "github.com/omakoto/go-common" line
// remove this line
// run `go get github.com/omakoto/go-common`
// run `go mod tidy`

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/creack/pty v1.1.11
	github.com/davecgh/go-spew v1.1.1
	github.com/mattn/go-isatty v0.0.9
	github.com/omakoto/go-common v0.0.0-20210223020755-bd49cd9ce44e
	github.com/pkg/term v0.0.0-20190109203006-aa71e9d9e942
	golang.org/x/sys v0.0.0-20190927073244-c990c680b611
)
