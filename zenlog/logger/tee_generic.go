// +build !linux

package logger

import "os"

func forward(in, out *os.File) error {
	return forward_simple(in, out)
}

func tee(in, out1, out2 *os.File) error {
	return tee_simple(in, out1, out2)
}
