package zenlog

import (
	"fmt"
	"github.com/omakoto/zenlog-go/zenlog/util"
)

// Version of the logger - command communication protocol.
const PROTOCOL_VERSION = 1

// Return the "signature" of the zenlog executable.
func Signature() string {
	return fmt.Sprintf("%s:[%d]", util.FindSelf(), PROTOCOL_VERSION)
}
