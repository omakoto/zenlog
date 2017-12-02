#!/bin/bash

medir="${0%/*}"

TEST_NAME=04
. "$medir/zenlog-test-common"

clear_log

cd "$medir"

run_zenlog <<EOF
echo ok|cat
echo a;cat data/fstab
echo a#tag
echo $_ZENLOG_E2E_EXIT_TIME >"$_ZENLOG_TIME_INJECTION_FILE"; exit
EOF

check_result