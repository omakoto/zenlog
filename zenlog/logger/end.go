package logger

import (
	"fmt"
	"github.com/omakoto/go-common/src/utils"
	"github.com/omakoto/zenlog/zenlog/config"
	"github.com/omakoto/zenlog/zenlog/util"
	"time"
)

type StopRequest struct {
	ExitStatus int
	FinishTime time.Time
}

type StopReply struct {
	NumLines int
}

func EndCommand(exitStatus int, wantLineNumber bool, clock utils.Clock) {
	config := config.InitConfigForCommands()

	// Create a requst.
	req := StopRequest{exitStatus, clock.Now()}
	util.Dump("stopRequest=", req)

	fingerprint := util.Fingerprint()

	// Send it.
	MustSendToLogger(config, utils.StringSlice(CommandEndCommand, fingerprint, util.MustMarshal(req)))

	// Wait for reply.
	ret := MustReceiveFromLogger(config, func(args []string) bool {
		return (len(args) == 3) && (args[0] == CommandEndCommand) && (args[1] == fingerprint)
	})
	reply := StopReply{}
	util.MustUnmarshal(ret[2], &reply)

	if wantLineNumber {
		fmt.Println(reply.NumLines)
	}
}
