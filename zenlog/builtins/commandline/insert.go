package commandline

import (
	"github.com/omakoto/zenlog-go/zenlog/builtins/commandline/cc"
	"github.com/omakoto/zenlog-go/zenlog/builtins/history"
	"github.com/omakoto/zenlog-go/zenlog/shell"
	"github.com/omakoto/zenlog-go/zenlog/util"
)

// InsertLogBash handles ALT-L -- first call will insert the last log in the command line, and the subsequent calls
// will replace it with a previous log.
func InsertLogBash(args []string) {

	proxy := shell.GetProxy()

	cc := cc.FromEnvironment("insert-log-san", proxy)
	util.Dump("cc=", cc)

	log := history.NthLastLog(cc.Config(), 0, cc.NumRepeat+1, history.LogTypeSan)
	util.Debugf("log=%s", log)

	if log == "" {
		if cc.NumRepeat > 0 {
			proxy.PrintUpdateCommandLineEvalStr(cc.FirstCommandLine, cc.FirstCursorPos)
		}
		cc.ClearSaved()
		return
	}
	replacement := shell.Escape(log) + " "
	cl := cc.FirstCommandLine[0:cc.FirstCursorPos] + replacement + cc.FirstCommandLine[cc.FirstCursorPos:]
	cp := cc.FirstCursorPos + len(replacement)

	cc.AfterCommandLine = cl
	cc.AfterCursorPos = cp

	proxy.PrintUpdateCommandLineEvalStr(cl, cp)

	cc.Save()
}
