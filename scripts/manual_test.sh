#!/bin/bash

set -e

cd "${0%/*}/.."

./scripts/build.sh

export PATH="$(readlink -m bin):$PATH"

export ZENLOG_DEBUG=0
if [[ "$1" == "-d" ]] ;then
    export ZENLOG_DEBUG=1
fi

export ZENLOG_START_COMMAND="exec /bin/bash --noprofile --rcfile $(readlink -m scripts)/bashrc.manual-test"
export ZENLOG_DIR="/tmp/zenlog-manual-test/"
export ZENLOG_PREFIX_COMMANDS="(?:builtin|time|sudo|command)"
export ZENLOG_ALWAYS_NO_LOG_COMMANDS="(?:vi|vim|man|nano|pico|less|watch|emacs|ssh|zenlog.*)"

zenlog
