// +build linux

package logger

import (
	"os"
	"syscall"

	"github.com/omakoto/zenlog-go/zenlog/util"
	"golang.org/x/sys/unix"
)

func forward(in, out *os.File) error {
	// TODO Is it actuall faster?
	r, w, e := os.Pipe()
	if e != nil {
		return nil
	}
	for {
		len, err := syscall.Splice(int(in.Fd()), nil, int(w.Fd()), nil, 1024*1024*32, unix.SPLICE_F_MOVE)
		if len == 0 {
			break
		}
		if util.Warn(err, "Splice failed") {
			return err
		}
		for len > 0 {
			len2, err := syscall.Splice(int(r.Fd()), nil, int(out.Fd()), nil, 1024*1024*32, unix.SPLICE_F_MOVE)
			if len2 == 0 {
				break
			}
			if util.Warn(err, "Splice failed") {
				return err
			}
			len -= len2
		}
	}
	return nil
}

func tee(in, out1, out2 *os.File) error {
	return tee_simple(in, out1, out2)
	//for {
	//	len, err := syscall.Tee(int(in.Fd()), int(out1.Fd()), 1024*128, 0 /*unix.SPLICE_F_NONBLOCK*/)
	//	if util.Warn(err, "Tee failed") {
	//		return err
	//	}
	//	if len == 0 {
	//		return nil
	//	}
	//
	//	for len > 0 {
	//		slen, err := syscall.Splice(int(in.Fd()), nil, int(out2.Fd()), nil, int(len), unix.SPLICE_F_MOVE)
	//		if util.Warn(err, "Splice failed") {
	//			return err
	//		}
	//		len -= slen
	//	}
	//}
}
