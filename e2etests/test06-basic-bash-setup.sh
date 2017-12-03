#!/bin/bash

medir="${0%/*}"

TEST_NAME=06
export ZENLOG_START_COMMAND="exec /bin/bash --noprofile --rcfile $medir/files/06/bashrc"
. "$medir/zenlog-test-common"

clear_log

cd "$medir"

run_zenlog <<EOF
cat data/fstab | grep -v -- '^#' #TAG
echo $_ZENLOG_E2E_EXIT_TIME >"$_ZENLOG_TIME_INJECTION_FILE"; exit
EOF

check_result