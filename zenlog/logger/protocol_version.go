package logger

import (
	"fmt"

	"github.com/omakoto/zenlog/zenlog/util"
)

// Version of the logger - command communication protocol.
const protocolVersion = 1

// Signature returns the "signature" of the zenlog executable.
func Signature() string {
	return fmt.Sprintf("%s:[%d]", util.FindZenlogBin(), protocolVersion)
}
