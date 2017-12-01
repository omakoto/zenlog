#!/bin/bash

set -e

cd "${0%/*}/.."

go test "${@}" ./zenlog/...
