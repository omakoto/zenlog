package logger

import (
	"io"
	"os"

	"github.com/omakoto/zenlog-go/zenlog/util"
)

func forward(in, out *os.File) error {
	_, err := io.Copy(out, in)
	return err
}

func tee(in, out1, out2 *os.File) error {
	buf := make([]byte, 32*1024)

	var err error
	for {
		nr, err := in.Read(buf)
		if nr > 0 {
			// First, write to stdout.
			nw, ew := out1.Write(buf[0:nr])
			util.Warn(ew, "Stdout.Write failed")
			if nr != nw {
				util.Warn(io.ErrShortWrite, "ErrShortWrite for Stdout")
			}
			// Then, write to l.
			nw, ew = out2.Write(buf[0:nr])
			util.Warn(ew, "ForwardPipe.Write failed")
			if nr != nw {
				util.Warn(io.ErrShortWrite, "ErrShortWrite for ForwardPipe")
			}
		}
		if err != nil {
			break
		}
	}
	if err != nil && err != io.EOF && err != io.ErrClosedPipe {
		util.Warn(err, "Forwarder finishing with an error")
	}
	return err
}
