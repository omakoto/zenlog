package builtins

import (
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"time"
)

const (
	COMMAND_START_COMMAND = "start-command"
	COMMAND_END_COMMAND   = "end-command"

	READ_TIMEOUT = time.Second * 1
)

func MustSendToLogger(config *config.Config, vals []string) {
	util.Check(util.WriteToFile(config.LoggerIn, vals), "Failed to send to logger.")
}

func MustReceiveFromLogger(config *config.Config, predicate func(vals []string) bool) []string {
	ret, err := util.ReadFromFile(config.LoggerOut, predicate, READ_TIMEOUT)
	util.Check(err, "Failed to receive from logger")
	return ret
}
