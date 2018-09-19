package logfiles

import (
	"github.com/omakoto/go-common/src/utils"
	"github.com/omakoto/zenlog/zenlog/config"
	"github.com/omakoto/zenlog/zenlog/util"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCreateLogFiles(t *testing.T) {
	config := config.Config{}

	os.Setenv("TZ", "America/Los_Angeles")

	config.LogDir = "/tmp/zenlog-test/log/"
	config.ZenlogPid = 111
	config.PrefixCommands = "(?:time|sudo|[a-zA-Z0-9_]+=.*)"

	os.RemoveAll(config.LogDir)

	// TODO The input time is GMT but the logfilename uses a local time.
	// This needs to generate a local time. How?
	clock := utils.NewInjectedClock(time.Unix(1319202062, 123*1000*1000))
	tests := []struct {
		commandLine string
		log         string
	}{
		{"/bin/echo ok", "/tmp/zenlog-test/log/SAN/2011/10/21/06-01-02.123-00111_+echo_ok.log"},
		{"/bin/echo ok # comment tag ", "/tmp/zenlog-test/log/SAN/2011/10/21/06-02-02.123-00111_+comment_tag_+echo_ok_comment_tag.log"},
		{"echo ok", "/tmp/zenlog-test/log/SAN/2011/10/21/06-03-02.123-00111_+echo_ok.log"},
		{"./echo ok", "/tmp/zenlog-test/log/SAN/2011/10/21/06-04-02.123-00111_+echo_ok.log"},

		// Note for log filenames, we do *not* use PREFIX_COMMAND.
		// PREFIX_COMMAND is only used to decide the directory name under $ZENLOG_DIR/cmds/.
		{"time echo ok", "/tmp/zenlog-test/log/SAN/2011/10/21/06-05-02.123-00111_+time_echo_ok.log"},
	}
	for _, v := range tests {
		actual := CreateAndOpenLogFiles(&config, clock.Now(), ParseCommandLine(&config, v.commandLine))
		defer actual.Close()
		util.AssertStringsEqual(t, v.commandLine, v.log, actual.SanFile)
		util.AssertStringsEqual(t, v.commandLine, strings.Replace(v.log, "SAN", "RAW", 1), actual.RawFile)
		util.AssertStringsEqual(t, v.commandLine, strings.Replace(v.log, "SAN", "ENV", 1), actual.EnvFile)
		util.AssertFileExist(t, actual.SanFile)
		util.AssertFileExist(t, actual.RawFile)
		util.AssertFileExist(t, actual.EnvFile)

		clock = utils.NewInjectedClock(clock.Now().Add(time.Minute))
	}
}
