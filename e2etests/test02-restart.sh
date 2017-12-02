#!/bin/bash

medir="${0%/*}"

TEST_NAME=02
. "$medir/zenlog-test-common"

# This is the default now.
# export ZENLOG_USE_EXPERIMENTAL_COMMAND_PARSER=1


clear_log

cd "$medir"

run_zenlog <<EOF
echo 'start'
zenlog_restart
EOF

check_result 02
