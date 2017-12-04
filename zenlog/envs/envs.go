package envs

const (
	// ZenlogConf is an environmental variable key for ~/.zenlog.toml.
	ZenlogConf = "ZENLOG_CONF"

	// ZenlogBin is an environmental variable key for the zenlog executable fullpath.
	ZenlogBin = "ZENLOG_BIN"

	// ZenlogSourceTop is an environmental variable key for the zenlog source top directory fullpath.
	ZenlogSourceTop = "ZENLOG_SRC_TOP"

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

	// ZenlogAutoflush is an environmental variable key.
	ZenlogAutoflush = "ZENLOG_AUTO_FLUSH"

	// ZenlogUseExperimentalCommandParser is an environmental variable key.
	ZenlogUseExperimentalCommandParser = "ZENLOG_USE_EXPERIMENTAL_COMMAND_PARSER"

	// ZenlogOuterTty is an environmental variable key.
	ZenlogOuterTty = "ZENLOG_OUTER_TTY"

	// ZenlogLoggerIn is an environmental variable key.
	ZenlogLoggerIn = "ZENLOG_LOGGER_IN"

	// ZenlogLoggerOut is an environmental variable key.
	ZenlogLoggerOut = "ZENLOG_LOGGER_OUT"
)
