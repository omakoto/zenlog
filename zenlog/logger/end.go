package logger

import (
	"fmt"
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"time"
)

type StopRequest struct {
	ExitStatus int
	FinishTime time.Time
}

type StopReply struct {
	NumLines int
}

func EndCommand(exitStatus int, wantLineNumber bool, clock util.Clock) {
	config := config.InitConfigForCommands()

	// Create a requst.
	req := StopRequest{exitStatus, clock.Now()}
	util.Dump("stopRequest=", req)

	fingerprint := util.Fingerprint()

	// Send it.
	MustSendToLogger(config, util.StringSlice(COMMAND_END_COMMAND, fingerprint, util.MustMarshal(req)))

	// Wait for reply.
	ret := MustReceiveFromLogger(config, func(args []string) bool {
		return (len(args) == 3) && (args[0] == COMMAND_END_COMMAND) && (args[1] == fingerprint)
	})
	reply := StopReply{}
	util.MustUnmarshal(ret[2], &reply)

	if wantLineNumber {
		fmt.Println(reply.NumLines)
	}

	util.ExitSuccess()
}
