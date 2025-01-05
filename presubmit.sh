#!/bin/bash

# Extract commands from .travis.yml and execute them.

set -e

cd "${0%/*}"

run() {
  echo "Running: $*"
  "$@"
}

./scripts/presubmit.sh
./scripts/e2e-test.sh
