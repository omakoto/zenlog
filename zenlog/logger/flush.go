package logger

import (
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/util"
)

func FlushCommand() {
	config := config.InitConfigForCommands()

	MustSendToLogger(config, util.StringSlice(FLUSH_COMMAND))

	util.ExitSuccess()
}
