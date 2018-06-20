#!/bin/bash

# Extract commands from .travis.yml and execute them.

set -e

cd "${0%/*}"

run() {
  echo "Running: $*"
  "$@"
}

. <(sed -ne 's/^ *- *\(.*\)#presubmit/run \1/p' .travis.yml)

echo '(Run [docker rmi $(docker images -f "dangling=true" -q)] to remove dangling images if needed.)'