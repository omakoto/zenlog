#!/bin/bash

set -e

. "${0%/*}"/build.sh

gofmt -s -d $(find . -type f -name '*.go') |& perl -pe 'END{exit($. > 0 ? 1 : 0)}'

go vet ./...

go test -v -race ./...
