package cc

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/omakoto/go-common/src/fileutils"
	"github.com/omakoto/go-common/src/shell"
	"github.com/omakoto/zenlog/zenlog/builtins/history"
	"github.com/omakoto/zenlog/zenlog/config"
	"github.com/omakoto/zenlog/zenlog/util"
)

// CommandLineContext represents "command line context" to detect consecutive calls of the same command.
type CommandLineContext struct {
	FirstCommandLine string
	FirstCursorPos   int

	BeforeCommandLine string
	BeforeCursorPos   int

	AfterCommandLine string
	AfterCursorPos   int

	LastLog string

	Operation string
	NumRepeat int

	config *config.Config
}

func filenameForConfig(c *config.Config) string {
	return fmt.Sprintf("%s/zenlog_lastcommand_%d.json", c.TempDir, c.ZenlogPid)
}

func FromEnvironment(operation string, proxy shell.Proxy) *CommandLineContext {
	config := config.InitConfigForCommands()

	ret := CommandLineContext{}
	ret.config = config
	ret.Operation = operation
	ret.BeforeCommandLine, ret.BeforeCursorPos = proxy.GetCommandLine()
	ret.LastLog = history.NthLastLog(config, 0, 1, history.LogTypeSan)

	// Load the last cc.
	last := FromLastFile()

	if last.Operation == ret.Operation &&
		last.AfterCommandLine == ret.BeforeCommandLine &&
		last.AfterCursorPos == ret.BeforeCursorPos &&
		last.LastLog == ret.LastLog {

		ret.FirstCommandLine = last.FirstCommandLine
		ret.FirstCursorPos = last.FirstCursorPos
		ret.NumRepeat = last.NumRepeat + 1

	} else {
		ret.FirstCommandLine = ret.BeforeCommandLine
		ret.FirstCursorPos = ret.BeforeCursorPos
	}

	return &ret
}

func FromLastFile() *CommandLineContext {
	config := config.InitConfigForCommands()
	file := filenameForConfig(config)

	ret := CommandLineContext{config: config}
	if !fileutils.FileExists(file) {
		return &ret
	}
	dat, err := ioutil.ReadFile(file)
	util.Check(err, "ReadFile failed")

	util.MustUnmarshal(string(dat), &ret)
	return &ret
}

func (cc *CommandLineContext) Save() {
	file := filenameForConfig(cc.config)

	dat := []byte(util.MustMarshal(&cc))

	err := ioutil.WriteFile(file, dat, 0600)
	util.Warn(err, "WriteFile failed")
}

func (cc *CommandLineContext) ClearSaved() {
	os.Remove(filenameForConfig(cc.config))
}

func (cc *CommandLineContext) Config() *config.Config {
	return cc.config
}
