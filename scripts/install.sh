#!/bin/bash

set -e

cd "${0%/*}/.."

out=bin
mkdir -p "$out"

go install ./zenlog/cmd/zenlog
