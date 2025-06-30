module github.com/omakoto/zenlog

go 1.24

//replace github.com/omakoto/go-common => ../go-common

// To update the "github.com/omakoto/go-common" line
// remove this line
// run `go get github.com/omakoto/go-common`
// run `go mod tidy`

require (
	github.com/BurntSushi/toml v1.4.0
	github.com/creack/pty v1.1.23
	github.com/davecgh/go-spew v1.1.1
	github.com/mattn/go-isatty v0.0.20
	github.com/omakoto/go-common v0.0.0-20210223020755-bd49cd9ce44e
	github.com/pkg/term v1.2.0-beta.2
	golang.org/x/sys v0.25.0
)
