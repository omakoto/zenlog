#!/bin/bash

set -e

cd "${0%/*}/.."

./docker/build-docker.sh
docker run -it --rm -t zenlog_docker /bin/bash -l -c '"go/src/github.com/omakoto/zenlog/scripts/e2e-test.sh"'