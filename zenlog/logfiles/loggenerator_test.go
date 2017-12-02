package logfiles

import (
	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCreateLogFiles(t *testing.T) {
	config := config.Config{}

	config.LogDir = "/tmp/zenlog-test/log/"
	config.ZenlogPid = 111

	os.RemoveAll(config.LogDir)

	// TODO The input time is GMT but the logfilename uses a local time.
	// This needs to generate a local time. How?
	clock := util.NewInjectedClock(time.Unix(1319202062, 123*1000*1000))
	tests := []struct {
		commandLine string
		log         string
	}{
		{"/bin/echo ok", "/tmp/zenlog-test/log/SAN/2011/10/21/06-01-02.123-00111_+_bin_echo_ok.log"},
		{"/bin/echo ok # comment tag ", "/tmp/zenlog-test/log/SAN/2011/10/21/06-02-02.123-00111_+comment_tag__+_bin_echo_ok_comment_tag.log"},
	}
	for _, v := range tests {
		actual := OpenLogFiles(&config, clock.Now(), ParseCommandLine(&config, v.commandLine))
		defer actual.Close()
		util.AssertStringsEqual(t, v.commandLine, v.log, actual.SanFile)
		util.AssertStringsEqual(t, v.commandLine, strings.Replace(v.log, "SAN", "RAW", 1), actual.RawFile)
		util.AssertStringsEqual(t, v.commandLine, strings.Replace(v.log, "SAN", "ENV", 1), actual.EnvFile)
		util.AssertFileExist(t, actual.SanFile)
		util.AssertFileExist(t, actual.RawFile)
		util.AssertFileExist(t, actual.EnvFile)

		clock = util.NewInjectedClock(clock.Now().Add(time.Minute))
	}
}
