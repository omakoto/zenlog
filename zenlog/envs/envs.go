package envs

const (
	// ZenlogConf is an environmental variable key for ~/.zenlog.toml.
	ZenlogConf = "ZENLOG_CONF"

	// ZenlogBin is an environmental variable key for the zenlog executable fullpath.
	ZenlogBin = "ZENLOG_BIN"

	// ZenlogSourceDir is an environmental variable key for the zenlog source top directory fullpath.
	ZenlogSourceDir = "ZENLOG_SRC_DIR"

	// ZenlogBinCtime is an environmental variable key for the zenlog binary timestamp.
	ZenlogBinCtime = "ZENLOG_BIN_CTIME"

	// ZenlogPid is an environmental variable key for the zenlog logger process PID.
	ZenlogPid = "ZENLOG_PID"

	// ZenlogSignature is an environmental variable key.
	ZenlogSignature = "ZENLOG_SIGNATURE"

	// ZenlogDir is an environmental variable key for the log directory.
	ZenlogDir = "ZENLOG_DIR"

	// ZenlogTemp is an environmental variable key for the temporary directory.
	ZenlogTemp = "ZENLOG_TEMP"

	// ZenlogAutoFlush is an environmental variable key.
	ZenlogAutoFlush = "ZENLOG_AUTO_FLUSH"

	// ZenlogUseExperimentalCommandParser is an environmental variable key.
	ZenlogUseExperimentalCommandParser = "ZENLOG_USE_EXPERIMENTAL_COMMAND_PARSER"

	// ZenlogOuterTty is an environmental variable key.
	ZenlogOuterTty = "ZENLOG_OUTER_TTY"

	// ZenlogLoggerIn is an environmental variable key.
	ZenlogLoggerIn = "ZENLOG_LOGGER_IN"

	// ZenlogLoggerOut is an environmental variable key.
	ZenlogLoggerOut = "ZENLOG_LOGGER_OUT"

	// If ZenlogUseSplice is true, zenlog uses splice(2) and tee(2) on linux.
	ZenlogUseSplice = "ZENLOG_USE_SPLICE"
)
