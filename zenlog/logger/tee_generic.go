// +build !linux amd64

package logger

import "os"

func forward(in, out *os.File) error {
	return forwardSimple(in, out)
}

func tee(in, out1, out2 *os.File) error {
	return teeSimple(in, out1, out2)
}
