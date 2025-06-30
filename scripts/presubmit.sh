#!/bin/bash

set -e

. "${0%/*}"/build.sh

gofmt -s -d $(find . -type f -name '*.go') 2>&1 | perl -pe 'END{exit($. > 0 ? 1 : 0)}'

go test -v -race ./...

# Static analysis
go vet ./...
staticcheck ./...
