#!/bin/bash

set -e

cd "${0%/*}/.."

pwd
ls -l

out=bin
mkdir -p "$out"

go build -o "$out/zenlog" ./zenlog/cmd/zenlog
