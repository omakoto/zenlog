package builtins

import (
	"encoding/json"
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/logfiles"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"os"
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

	vals := make([]string, 2)
	vals[0] = COMMAND_START_COMMAND
	vals[1] = string(req.MustEncode())
	MustSendToLogger(config, vals[:])

	os.Exit(0)
}

func (s *StartRequest) MustEncode() []byte {
	dat, err := json.Marshal(s)
	util.Check(err, "Stringfy failed")
	return dat
}

func DecodeStartRequest(data string) (*StartRequest, error) {
	ret := StartRequest{}
	err := json.Unmarshal([]byte(data), &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
