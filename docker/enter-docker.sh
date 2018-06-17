#!/bin/bash

set -e
cd "${0%/*}/../"

name=zenlog_docker

docker run $DOCKER_RUN_OPTS -it --rm -t $name .

