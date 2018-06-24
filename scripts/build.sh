#!/bin/bash

set -e

cd "${0%/*}/.."

out=bin
mkdir -p "$out"

go build -o "$out/zenlog" ./zenlog/cmd/zenlog
