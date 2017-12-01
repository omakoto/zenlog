#!/bin/bash

medir="${0%/*}"

. "$medir/zenlog-test-common"

clear_log

cd "$medir"

#export ZENLOG_DEBUG=1 ; echo $'command cat /etc/fstab\nexit' | run_zenlog; exit 0

run_zenlog <<EOF
echo ok; tick 3
cat data/fstab | grep -v -- '^#'
man man
q
zenlog history # history 1
echo ok | cat # tag test abc  def <>/
zenlog current-log # com current log
zenlog last-log # com last log
zenlog current-log -r # com r current log
zenlog last-log -r # com r last log
true && echo "and test" # and test
false || echo "or test" # or test
cat data/fstab | fgrep dev
fgrep dev < data/fstab
186 fgrep dev < data/fstab
184 cat data/fstab
command cat data/fstab
/usr/bin/cat data/fstab
ABC="1 2 3" cat data/fstab
cat data/fstab | command sed -ne '1p'
in_zenlog && echo "in zenlog"
zenlog in_zenlog && echo "in zenlog"
zenlog_current_log # fun current log
zenlog_last_log # fun last log
zenlog_current_log -r # fun r current log
zenlog_last_log -r # fun r last log
zenlog history # history 2
zenlog history -r # history raw
echo $'a\xffb' # broken utf8
export ZENLOG_PID=10000
zenlog current-log -r
zenlog last-log -r
zenlog history
zenlog history -r
zenlog current-log -r -p $_ZENLOG_LOGGER_PID
zenlog last-log -r -p $_ZENLOG_LOGGER_PID
zenlog history -p $_ZENLOG_LOGGER_PID
zenlog history -r -p $_ZENLOG_LOGGER_PID
zenlog_current_log -e # fun e current log
zenlog_last_log -e # fun e last log
zenlog current-log -e -p $_ZENLOG_LOGGER_PID
zenlog last-log -e -p $_ZENLOG_LOGGER_PID
zenlog history -e -p $_ZENLOG_LOGGER_PID
export ZENLOG_PID=$_ZENLOG_LOGGER_PID
cat data/fstab|cat -n|cat -E
cat data/* #Wildcard
echo $_ZENLOG_E2E_EXIT_TIME >"$_ZENLOG_TIME_INJECTION_FILE"; exit
EOF

fail=0
echo "Checking tree..."
diff --color=always -c $medir/results/tree1.txt <($TREE -aF "$ZENLOG_DIR")
fail=$(( $fail || $?))

echo "Checking log files..."
diff --color=always  -X $medir/diff-ignore-files.txt \
    -ur $medir/results/files "$ZENLOG_DIR"
fail=$(( $fail || $?))

if (( $fail )); then
    echo "Some tests failed."
else
    echo "All tests passed."
fi
