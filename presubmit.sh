#!/bin/bash

set -e

cd "${0%/*}"

run() {
  echo "Running: $*"
  "$@"
}

./scripts/presubmit.sh
./scripts/e2e-test.sh
