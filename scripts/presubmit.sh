#!/bin/bash

set -e

go get github.com/golang/lint/golint                        # Linter
go get honnef.co/go/tools/cmd/megacheck                     # Badass static analyzer/linter

NOT_FORMATTED=$(gofmt -s -d $(find . -type f -name '*.go'))

if [[ -n "$NOT_FORMATTED" ]] ; then
    echo $'Wrong format:\n'"$NOT_FORMATTED"
    exit 1
fi

go test -v -race ./...                   # Run all the tests with the race detector enabled
go vet ./...                             # go vet is the official Go static analyzer
megacheck ./...                          # "go vet on steroids" + linter
# golint -set_exit_status $(go list ./...) # one last linter
