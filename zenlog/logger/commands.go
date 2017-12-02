package logger

import (
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"time"
)

const (
	CloseCommand = "close"
	FlushCommand = "flush"

	CommandStartCommand = "start-command"
	CommandEndCommand   = "end-command"

	// This is a command sent by the signal handler on SIGCHLD.
	ChildDiedCommand = "child-died"

	readTimeout = time.Second * 1
)

func MustSendToLogger(config *config.Config, vals []string) {
	util.Check(util.WriteToFile(config.LoggerIn, vals), "Failed to send to logger.")
}

func MustReceiveFromLogger(config *config.Config, predicate func(vals []string) bool) []string {
	ret, err := util.ReadFromFile(config.LoggerOut, predicate, readTimeout)
	util.Check(err, "Failed to receive from logger")
	return ret
}
