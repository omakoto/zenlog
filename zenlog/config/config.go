package config

import (
	"github.com/BurntSushi/toml"
	"github.com/davecgh/go-spew/spew"
	"github.com/omakoto/zenlog-go/zenlog/envs"
	"github.com/omakoto/zenlog-go/zenlog/util"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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
	var config Config
	data, err := ioutil.ReadFile(file)
	if err == nil {
		if _, err := toml.Decode(string(data), &config); err != nil {
			util.Fatalf("Unable to load %s: %s", file, err)
		}
	} else if os.IsNotExist(err) {
		util.Warn(err, "%s doesn't exist; using the default instead", file)
	}

	overwriteWithEnviron(&config.StartCommand, "ZENLOG_START_COMMAND", "")
	overwriteWithEnviron(&config.LogDir, envs.ZENLOG_DIR, "/tmp/zenlog/")
	overwriteWithEnviron(&config.PrefixCommands, "ZENLOG_PREFIX_COMMANDS", `(?:command|builtin|time|sudo|[a-zA-Z0-9_]+\=.*)`)
	overwriteWithEnviron(&config.AlwaysNoLogCommands, "ZENLOG_ALWAYS_NO_LOG_COMMANDS", `(?:vi|vim|man|nano|pico|less|watch|emacs|zenlog.*)`)
	overwriteWithEnviron(&config.CommandSplitter, "ZENLOG_COMMAND_SPLITTER", "")
	overwriteWithEnviron(&config.CommentSplitter, "ZENLOG_COMMENT_SPLITTER", "")
	overwriteWithEnviron(&config.TempDir, envs.ZENLOG_TEMP, "")
	overwriteBoolWithEnviron(&config.AutoFlush, envs.ZENLOG_AUTOFLUSH)

	if config.StartCommand == "" {
		shell := os.Getenv("SHELL")
		if shell != "" {
			config.StartCommand = "exec " + shell + " -l"
		} else {
			config.StartCommand = "exec /bin/sh -l"
		}
	}

	if config.TempDir == "" {
		config.TempDir = util.FirstNonEmpty(os.Getenv("TEMP"), os.Getenv("TMP"), "/tmp/")
	}

	ensureSlash(&config.LogDir)
	ensureSlash(&config.TempDir)

	config.ZenlogPid = os.Getpid()

	if util.Debug {
		util.Debugf("Config=%s", spew.Sdump(&config))
	}

	return &config
}

func InitConfigForCommands() *Config {
	var config Config

	pid, err := strconv.Atoi(os.Getenv(envs.ZENLOG_PID))
	util.Check(err, "ZENLOG_PID not integer")
	config.ZenlogPid = pid

	config.LogDir = os.Getenv(envs.ZENLOG_DIR)
	config.OuterTty = os.Getenv(envs.ZENLOG_OUTER_TTY)
	config.LoggerIn = os.Getenv(envs.ZENLOG_LOGGER_IN)
	config.LoggerOut = os.Getenv(envs.ZENLOG_LOGGER_OUT)

	if config.ZenlogPid == 0 {
		util.Fatalf(envs.ZENLOG_PID + " not set.")
	}
	if config.LogDir == "" {
		util.Fatalf(envs.ZENLOG_DIR + " not set.")
	}
	if config.OuterTty == "" {
		util.Fatalf(envs.ZENLOG_OUTER_TTY + " not set.")
	}
	if config.LoggerIn == "" {
		util.Fatalf(envs.ZENLOG_LOGGER_IN + " not set.")
	}
	if config.LoggerOut == "" {
		util.Fatalf(envs.ZENLOG_LOGGER_OUT + " not set.")
	}
	return &config
}
