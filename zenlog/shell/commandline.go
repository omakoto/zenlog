package shell

import (
	"os"
	"path/filepath"
)

type Proxy interface {
	GetCommandLine() (commandLine string, cursorPos int)
	PrintUpdateCommandLineEvalStr(commandLine string, cursorPos int)
}

type nullShellProxy struct {
}

func (s *nullShellProxy) GetCommandLine() (string, int) {
	return "", 0
}

func (s *nullShellProxy) PrintUpdateCommandLineEvalStr(commandLine string, cursorPos int) {
}

func GetProxy() Proxy {
	shell := filepath.Base(os.Getenv("SHELL"))

	if shell == "bash" {
		return &BashProxy{}
	}

	return &nullShellProxy{}
}
