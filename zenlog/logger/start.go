package logger

import (
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/logfiles"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"strings"
	"time"
)

type StartRequest struct {
	Command   logfiles.Command
	LogFiles  logfiles.LogFiles
	StartTime time.Time
}

func StartCommand(envs string, commandLineArray []string, clock util.Clock) {
	config := config.InitConfigForCommands()

	commandLine := strings.Join(commandLineArray, " ")
	command := logfiles.ParseCommandLine(config, commandLine)

	// Open the log file.
	now := clock.Now()
	logFiles := logfiles.OpenLogFiles(config, now, command)
	defer logFiles.Close()

	logFiles.WriteEnv(command, envs, now)

	// Send the start request to the logger.
	req := StartRequest{*command, logFiles, clock.Now()}
	util.Dump("startRequest=", req)

	MustSendToLogger(config, util.StringSlice(COMMAND_START_COMMAND, util.MustMarshal(req)))

	util.ExitSuccess()
}
