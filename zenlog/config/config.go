package config

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/omakoto/zenlog-go/zenlog/envs"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"runtime"
)

// Config represents configuration parameters, read from ~/.zenlog.toml and overridden with the environmental variables.
type Config struct {
	LogDir              string `toml:"ZENLOG_DIR"`
	TempDir             string `toml:"ZENLOG_TEMP"`
	StartCommand        string `toml:"ZENLOG_START_COMMAND"`
	PrefixCommands      string `toml:"ZENLOG_PREFIX_COMMANDS"`
	AlwaysNoLogCommands string `toml:"ZENLOG_ALWAYS_NO_LOG_COMMANDS"`
	AutoFlush           bool   `toml:"ZENLOG_AUTO_FLUSH"`

	UseExperimentalCommandParser bool   `toml:"ZENLOG_USE_EXPERIMENTAL_COMMAND_PARSER"`
	CommandSplitter              string `toml:"ZENLOG_COMMAND_SPLITTER"`
	CommentSplitter              string `toml:"ZENLOG_COMMENT_SPLITTER"`

	Maxproc int `toml:"ZENLOG_MAXPROC"`

	ZenlogPid int

	OuterTty  string
	LoggerIn  string
	LoggerOut string
}

// If a string is empty, get from a given environmental variabne.
func overwriteWithEnviron(to *string, envKey string, def string) {
	if val := os.Getenv(envKey); val != "" {
		*to = val
	} else {
		*to = os.ExpandEnv(*to)
	}
	if *to == "" {
		*to = def
	}
}

func overwriteBoolWithEnviron(to *bool, envKey string) {
	if val := os.Getenv(envKey); val != "" {
		*to = (val == "1")
	}
}

// Make sure a string ends with a slash.
func ensureSlash(v *string) {
	if strings.HasSuffix(*v, "/") {
		return
	}
	*v = *v + "/"
}

// InitConfigForLogger returns a Config for a new session, loading from ~/.zenlog.toml and the environmental variables.
func InitConfigForLogger() *Config {
	file := os.Getenv(envs.ZenlogConf)
	if file == "" {
		file = os.ExpandEnv("$HOME/.zenlog.toml")
	}
	var c Config

	c.UseExperimentalCommandParser = true // Default to true.

	data, err := ioutil.ReadFile(file)
	if err == nil {
		if _, err := toml.Decode(string(data), &c); err != nil {
			util.Fatalf("Unable to load %s: %s", file, err)
		}
	} else if os.IsNotExist(err) {
		util.Warn(err, "%s doesn't exist; using the default instead", file)
	}

	overwriteWithEnviron(&c.StartCommand, "ZENLOG_START_COMMAND", "")
	overwriteWithEnviron(&c.LogDir, envs.ZenlogDir, os.ExpandEnv("$HOME/zenlog/"))
	overwriteWithEnviron(&c.PrefixCommands, "ZENLOG_PREFIX_COMMANDS", `(?:command|builtin|time|sudo|[a-zA-Z0-9_]+\=.*)`)
	overwriteWithEnviron(&c.AlwaysNoLogCommands, "ZENLOG_ALWAYS_NO_LOG_COMMANDS", `(?:vi|vim|man|nano|pico|emacs|zenlog.*)`)

	overwriteWithEnviron(&c.CommandSplitter, "ZENLOG_COMMAND_SPLITTER", "")
	overwriteWithEnviron(&c.CommentSplitter, "ZENLOG_COMMENT_SPLITTER", "")

	overwriteWithEnviron(&c.TempDir, envs.ZenlogTemp, "")

	overwriteBoolWithEnviron(&c.AutoFlush, envs.ZenlogAutoFlush)
	overwriteBoolWithEnviron(&c.UseExperimentalCommandParser, envs.ZenlogUseExperimentalCommandParser)

	if c.Maxproc < 1 {
		c.Maxproc = 1
	}

	if c.StartCommand == "" {
		shell := os.Getenv("SHELL")
		if shell != "" {
			c.StartCommand = "exec " + shell + " -l"
		} else {
			c.StartCommand = "exec /bin/sh -l"
		}
	}

	if c.TempDir == "" || !util.DirExists(c.TempDir) {
		c.TempDir = util.FirstNonEmpty(os.Getenv("TEMP"), os.Getenv("TMP"), "/tmp/")
	}

	ensureSlash(&c.LogDir)
	ensureSlash(&c.TempDir)

	// For E2E testing, override the PID with _ZENLOG_LOGGER_PID, if set.
	c.ZenlogPid = util.GetIntEnv("_ZENLOG_LOGGER_PID", os.Getpid())

	runtime.GOMAXPROCS(c.Maxproc)

	util.Dump("Logger config=", c)

	return &c
}

// InitConfigForCommands returns a Config for subcommands, loading from ~/.zenlog.toml and the environmental variables.
// Some of the parameters (such as ZenlogPid) will be inherited from the current zenlog session.
func InitConfigForCommands() *Config {
	var c Config

	// Take over some of the parameters from the logger.

	pid, err := strconv.Atoi(os.Getenv(envs.ZenlogPid))
	util.Check(err, "ZENLOG_PID not integer")
	c.ZenlogPid = pid

	c.LogDir = os.Getenv(envs.ZenlogDir)
	c.OuterTty = os.Getenv(envs.ZenlogOuterTty)
	c.LoggerIn = os.Getenv(envs.ZenlogLoggerIn)
	c.LoggerOut = os.Getenv(envs.ZenlogLoggerOut)

	if c.ZenlogPid == 0 {
		util.Fatalf(envs.ZenlogPid + " not set.")
	}
	if c.LogDir == "" {
		util.Fatalf(envs.ZenlogDir + " not set.")
	}
	if c.OuterTty == "" {
		util.Fatalf(envs.ZenlogOuterTty + " not set.")
	}
	if c.LoggerIn == "" {
		util.Fatalf(envs.ZenlogLoggerIn + " not set.")
	}
	if c.LoggerOut == "" {
		util.Fatalf(envs.ZenlogLoggerOut + " not set.")
	}

	// We still need to load certain parameters from TOML.
	lc := InitConfigForLogger()
	c.AlwaysNoLogCommands = lc.AlwaysNoLogCommands
	c.PrefixCommands = lc.PrefixCommands
	c.CommandSplitter = lc.CommandSplitter
	c.CommentSplitter = lc.CommentSplitter
	c.UseExperimentalCommandParser = lc.UseExperimentalCommandParser
	c.TempDir = lc.TempDir
	c.Maxproc = lc.Maxproc

	util.Dump("Command config=", c)
	return &c
}
