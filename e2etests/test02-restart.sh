#!/bin/bash

medir="${0%/*}"

TEST_NAME=02
. "$medir/zenlog-test-common"

clear_log

cd "$medir"

# This is tricky.
# We want to run 3 commands:
#   echo start
#   zenlog_restart
#   echo post-restart
# However, before sending the last echo, we need to make sure zenlog has
# restarted. Otherwise the last command would just go to the first session's
# pipe buffer, which will be cleared with the re-exec, so the second session
# will never see it.
# So we use a coproc to write to zenlog and also read from zenlog.
# To send the output to coprop, normally you'd use >&"${ioproc[1]}", but
# that wouldn't work because script(1) would buffer the output if stdout
# is not a terminal, and stdbuf(1) didn't work either.
# So here, we pass the coproc's stdin FD to script(1) as a logfile.
# and also use -f to make sure script would flush the output.

coproc ioproc {
    # Send to the first session.
    echo 'echo start'
    echo 'zenlog_restart'

    # Wait until we see the zenlog start message *twice*.
    count=0
    while read line; do
        if [[ "$line" = *ZENLOG_DIR* ]] ; then
            count=$(( $count + 1 ))

            if (( $count == 2 )) ; then
                break
            fi
        fi
    done

    # Send this to the second session.
    echo 'echo post-restart'
    echo exit
}
script -qefc "$ZENLOG_BIN" /proc/$$/fd/${ioproc[1]} <&"${ioproc[0]}"

check_result 02
