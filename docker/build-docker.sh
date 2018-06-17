#!/bin/bash

set -e
cd "${0%/*}/../"

name=zenlog_docker

docker build $DOCKER_BUILD_OPTS -t $name .
