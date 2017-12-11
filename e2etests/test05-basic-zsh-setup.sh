#!/bin/bash

medir="${0%/*}"

TEST_NAME=05
export ZDOTDIR=./e2etests/files/05/zshrc
export ZENLOG_START_COMMAND="exec /usr/bin/zsh -l"
. "$medir/zenlog-test-common"

clear_log

cd "$medir"

run_zenlog <<EOF
cat data/fstab | grep -v -- '^#' #TAG
echo $_ZENLOG_E2E_EXIT_TIME >"$_ZENLOG_TIME_INJECTION_FILE"; exit
EOF

check_result