#!/bin/bash

set -e

go get github.com/golang/lint/golint                        # Linter
go get honnef.co/go/tools/cmd/megacheck                     # Badass static analyzer/linter
