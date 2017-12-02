package logger

import (
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"time"
)

const (
	// CloseSessionCommand tells the logger to finish successfully.
	CloseSessionCommand = "close"

	// FlushCommand tells the logger to flush the currently open log files.
	FlushCommand = "flush"

	// CommandStartCommand tells the logger when a command starts. The argument is StartRequest.
	CommandStartCommand = "start-command"

	// CommandEndCommand tells the logger when a command finishes. The argument is StopRequest.
	CommandEndCommand = "end-command"

	readTimeout = time.Second * 1
)

// MustSendToLogger writes a command (which is basically just a set of strings) to the named pipe to the logger.
func MustSendToLogger(config *config.Config, vals []string) {
	util.Check(util.WriteToFile(config.LoggerIn, vals), "Failed to send to logger.")
}

// MustReceiveFromLogger reads a command (which is basically just a set of strings) from the named pipe from the logger.
// acceptor decides if a received command is what the caller is waiting for or not.
func MustReceiveFromLogger(config *config.Config, acceptor func(vals []string) bool) []string {
	ret, err := util.ReadFromFile(config.LoggerOut, acceptor, readTimeout)
	util.Check(err, "Failed to receive from logger")
	return ret
}
