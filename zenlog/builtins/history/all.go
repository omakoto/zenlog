package history

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/omakoto/zenlog-go/zenlog/config"
	"github.com/omakoto/zenlog-go/zenlog/logfiles"
	"github.com/omakoto/zenlog-go/zenlog/util"
)

// AllCommandsAndLogCommand implements "all-commands" subcommand, which lists all log files with their command lines.
func AllCommandsAndLogCommand(args []string) {
	config := config.InitConfigForCommands()

	flags := flag.NewFlagSet("zenlog all-log", flag.ExitOnError)
	r := flags.Bool("r", false, "Print RAW filename instead")
	e := flags.Bool("e", false, "Print ENV filename instead")
	n := flags.Float64("n", 30, "Only print log within last n days. [default=30]")
	//c := flags.Bool("c", false, "Limit to current zenlog session")

	flags.Parse(args)

	now := util.GetInjectedNow(util.NewClock())
	wf := func(path string, info os.FileInfo, err error) error {
		if now.Sub(info.ModTime()).Hours() > *n {
			return nil
		}
		if info.Mode().IsRegular() {
			i, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer i.Close()
			br := bufio.NewReader(i)
			first, err := br.ReadString('\n')
			if err != nil {
				return nil
			}
			first = strings.TrimLeft(first, "$")
			first = strings.TrimLeft(first, " ")

			if *r {
				path = strings.Replace(path, logfiles.SanDir, logfiles.RawDir, 1)
			} else if *e {
				path = strings.Replace(path, logfiles.SanDir, logfiles.EnvDir, 1)
			}

			fmt.Print(path)
			fmt.Print(" ")
			fmt.Print(first)
		}
		return nil
	}

	filepath.Walk(config.LogDir+logfiles.SanDir, wf)

	util.Exit(true)
}
