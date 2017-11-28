package logger

import (
	"encoding/json"
	"fmt"
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"os"
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
	vals := make([]string, 3)
	vals[0] = COMMAND_END_COMMAND
	vals[1] = fingerprint
	vals[2] = string(req.MustEncode())
	MustSendToLogger(config, vals[:])

	// Wait for reply.
	ret := MustReceiveFromLogger(config, func(args []string) bool {
		return (len(args) == 3) && (args[0] == COMMAND_END_COMMAND) && (args[1] == fingerprint)
	})
	reply, err := DecodeStopReply(ret[2])
	util.Check(err, "DecodeStopReply failed")

	if wantLineNumber {
		fmt.Println(reply.NumLines)
	}

	os.Exit(0)
}

func (s *StopRequest) MustEncode() []byte {
	dat, err := json.Marshal(s)
	util.Check(err, "Stringfy failed")
	return dat
}

func DecodeStopRequest(data string) (*StopRequest, error) {
	ret := StopRequest{}
	err := json.Unmarshal([]byte(data), &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (s *StopReply) MustEncode() []byte {
	dat, err := json.Marshal(s)
	util.Check(err, "Stringfy failed")
	return dat
}

func DecodeStopReply(data string) (*StopReply, error) {
	ret := StopReply{}
	err := json.Unmarshal([]byte(data), &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
