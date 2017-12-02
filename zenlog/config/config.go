package config

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/omakoto/zenlog-go/zenlog/envs"
	"github.com/omakoto/zenlog-go/zenlog/util"
)

// Represents configuration parameters.
type Config struct {
	LogDir              string `toml:"ZENLOG_DIR"`
	TempDir             string `toml:"ZENLOG_TEMP"`
	StartCommand        string `toml:"ZENLOG_START_COMMAND"`
	PrefixCommands      string `toml:"ZENLOG_PREFIX_COMMANDS"`
	AlwaysNoLogCommands string `toml:"ZENLOG_ALWAYS_NO_LOG_COMMANDS"`
	AutoFlush           bool   `toml:"ZENLOG_AUTO_FLUSH"`
	CommandSplitter     string `toml:"ZENLOG_COMMAND_SPLITTER"`
	CommentSplitter     string `toml:"ZENLOG_COMMENT_SPLITTER"`

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

// Initialize a config, loading from ~/.zenlog.toml and the environmental variables.
func InitConfigiForLogger() *Config {
	file := os.Getenv(envs.ZENLOG_CONF)
	if file == "" {
		file = os.ExpandEnv("$HOME/.zenlog.toml")
	}
	var c Config
	data, err := ioutil.ReadFile(file)
	if err == nil {
		if _, err := toml.Decode(string(data), &c); err != nil {
			util.Fatalf("Unable to load %s: %s", file, err)
		}
	} else if os.IsNotExist(err) {
		util.Warn(err, "%s doesn't exist; using the default instead", file)
	}

	overwriteWithEnviron(&c.StartCommand, "ZENLOG_START_COMMAND", "")
	overwriteWithEnviron(&c.LogDir, envs.ZENLOG_DIR, os.ExpandEnv("$HOME/zenlog/"))
	overwriteWithEnviron(&c.PrefixCommands, "ZENLOG_PREFIX_COMMANDS", `(?:command|builtin|time|sudo|[a-zA-Z0-9_]+\=.*)`)
	overwriteWithEnviron(&c.AlwaysNoLogCommands, "ZENLOG_ALWAYS_NO_LOG_COMMANDS", `(?:vi|vim|man|nano|pico|emacs|zenlog.*)`)
	overwriteWithEnviron(&c.CommandSplitter, "ZENLOG_COMMAND_SPLITTER", "")
	overwriteWithEnviron(&c.CommentSplitter, "ZENLOG_COMMENT_SPLITTER", "")
	overwriteWithEnviron(&c.TempDir, envs.ZENLOG_TEMP, "")
	overwriteBoolWithEnviron(&c.AutoFlush, envs.ZENLOG_AUTOFLUSH)

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

	c.ZenlogPid = util.GetLoggerPid()

	util.Dump("Logger config=", c)

	return &c
}

func InitConfigForCommands() *Config {
	var c Config

	// Take over some of the parameters from the logger.

	pid, err := strconv.Atoi(os.Getenv(envs.ZENLOG_PID))
	util.Check(err, "ZENLOG_PID not integer")
	c.ZenlogPid = pid

	c.LogDir = os.Getenv(envs.ZENLOG_DIR)
	c.OuterTty = os.Getenv(envs.ZENLOG_OUTER_TTY)
	c.LoggerIn = os.Getenv(envs.ZENLOG_LOGGER_IN)
	c.LoggerOut = os.Getenv(envs.ZENLOG_LOGGER_OUT)

	if c.ZenlogPid == 0 {
		util.Fatalf(envs.ZENLOG_PID + " not set.")
	}
	if c.LogDir == "" {
		util.Fatalf(envs.ZENLOG_DIR + " not set.")
	}
	if c.OuterTty == "" {
		util.Fatalf(envs.ZENLOG_OUTER_TTY + " not set.")
	}
	if c.LoggerIn == "" {
		util.Fatalf(envs.ZENLOG_LOGGER_IN + " not set.")
	}
	if c.LoggerOut == "" {
		util.Fatalf(envs.ZENLOG_LOGGER_OUT + " not set.")
	}

	// We still need to load certain parameters from TOML.
	lc := InitConfigiForLogger()
	c.AlwaysNoLogCommands = lc.AlwaysNoLogCommands
	c.PrefixCommands = lc.PrefixCommands
	c.CommandSplitter = lc.CommandSplitter
	c.CommentSplitter = lc.CommentSplitter

	util.Dump("Command config=", c)
	return &c
}
