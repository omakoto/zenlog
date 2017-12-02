#!/bin/bash

medir="${0%/*}"

TEST_NAME=03
. "$medir/zenlog-test-common"

clear_log

cd "$medir"

run_zenlog <<EOF
echo "|ok" #tag1
echo "a"#tag2
echo "a";#tag3
echo $_ZENLOG_E2E_EXIT_TIME >"$_ZENLOG_TIME_INJECTION_FILE"; exit
EOF

check_result