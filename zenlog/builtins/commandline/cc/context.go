package cc

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/omakoto/zenlog-go/zenlog/builtins/history"
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/shell"
	"github.com/omakoto/zenlog-go/zenlog/util"
)

// CC represents "command line context" to detect consecutive calls of the same command.
type CC struct {
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

func FromEnvironment(operation string, proxy shell.Proxy) *CC {
	config := config.InitConfigForCommands()

	ret := CC{}
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

func FromLastFile() *CC {
	config := config.InitConfigForCommands()
	file := filenameForConfig(config)

	ret := CC{config: config}
	if !util.FileExists(file) {
		return &ret
	}
	dat, err := ioutil.ReadFile(file)
	util.Check(err, "ReadFile failed")

	util.MustUnmarshal(string(dat), &ret)
	return &ret
}

func (cc *CC) Save() {
	file := filenameForConfig(cc.config)

	dat := []byte(util.MustMarshal(&cc))

	err := ioutil.WriteFile(file, dat, 0600)
	util.Warn(err, "WriteFile failed")
}

func (cc *CC) ClearSaved() {
	os.Remove(filenameForConfig(cc.config))
}

func (cc *CC) Config() *config.Config {
	return cc.config
}
