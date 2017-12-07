#!/bin/bash

set -e

cd "${0%/*}/.." # CD to the top directory.

export ZENLOG_DEBUG=0
export ZENLOG_START_COMMAND="exec /bin/bash --noprofile --rcfile $(readlink -m scripts)/bashrc.manual-test"

while getopts "drz" opt; do
    case "$opt" in
        d)  # Debug mode.
            export ZENLOG_DEBUG=1
            ;;
        r)  # Real bash mode.
            export ZENLOG_START_COMMAND="exec /bin/bash -l"
            ;;
        z)  # Run zsh instead.
            export ZDOTDIR=./scripts/zshrc
            export ZENLOG_START_COMMAND="exec /usr/bin/zsh -l"
            ;;

    a) app=1 ;;
  esac
done
shift $(($OPTIND - 1))

./scripts/build.sh

export PATH="$(readlink -m bin):$PATH"

export ZENLOG_CONF=$(pwd)/dot_zenlog.toml

export ZENLOG_DIR="/tmp/zenlog-manual-test/"
export ZENLOG_PREFIX_COMMANDS="(?:builtin|time|sudo|command)"
export ZENLOG_ALWAYS_NO_LOG_COMMANDS="(?:vim?|nano|pico|emacs|zenlog.*)"

export ZENLOG_BIN="$(readlink -m bin/zenlog)"

export PS1_OVERRIDE="[zenlog-testing]"

unset ZENLOG_NO_DEFAULT_BINDING ZENLOG_NO_DEFAULT_PROMPT

mkdir -p "$ZENLOG_DIR"

"$ZENLOG_BIN"
