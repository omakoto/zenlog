// +build !linux

package logger

import "os"

func tee_faster(in, out1, out2 *os.File) error {
	return tee(in, out1, out2)
}
