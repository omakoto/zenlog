//go:build linux && (amd64 || arm64)

package logger

import (
	"os"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/omakoto/zenlog/zenlog/util"
	"golang.org/x/sys/unix"
)

func forward(in, out *os.File) error {
	return forwardSimple(in, out)
	//// TODO Is it actually faster?
	//r, w, e := os.Pipe()
	//if e != nil {
	//	return nil
	//}
	//for {
	//	len, err := syscall.Splice(int(in.Fd()), nil, int(w.Fd()), nil, 1024*1024*32, unix.SPLICE_F_MOVE)
	//	if len == 0 {
	//		return nil
	//	}
	//	if len == -1 && util.Warn(err, "Splice failed (a1/2)") {
	//		return err
	//	}
	//	for len > 0 {
	//		len2, err := syscall.Splice(int(r.Fd()), nil, int(out.Fd()), nil, int(len), unix.SPLICE_F_MOVE)
	//		if len2 == 0 {
	//			return nil
	//		}
	//		if len2 == -1 && util.Warn(err, "Splice failed (a2/2)") {
	//			return err
	//		}
	//		len -= len2
	//	}
	//}
}

func tee(in, out1, out2 *os.File) error {
	if os.Getenv("ZENLOG_USE_SPLICE") != "1" {
		return teeSimple(in, out1, out2)
	}
	// TODO Is it actually faster?
	r, w, e := os.Pipe()
	if e != nil {
		return nil
	}
	for {
		tlen, err := syscall.Splice(int(in.Fd()), nil, int(w.Fd()), nil, 1024*1024*32, unix.SPLICE_F_MOVE)
		if tlen == 0 {
			return nil
		}
		if tlen == -1 {
			if err == syscall.Errno(5) {
				// EIO.
				// man 2 splice doesn't mention it, but we get it when zenlog closes.
				return nil
			}
			if util.Warn(err, "Splice failed (b1/2)") {
				util.Say("err=%s %d", spew.Sdump(err), err)
				return err
			}
		}
		for tlen > 0 {
			len, err := syscall.Tee(int(r.Fd()), int(out1.Fd()), int(tlen), 0 /*unix.SPLICE_F_NONBLOCK*/)
			if len == 0 {
				return nil
			}
			if len == -1 && util.Warn(err, "Tee failed") {
				return err
			}
			for len > 0 {
				len2, err := syscall.Splice(int(r.Fd()), nil, int(out2.Fd()), nil, int(len), unix.SPLICE_F_MOVE)
				if len2 == 0 {
					return nil
				}
				if len2 == -1 && util.Warn(err, "Splice failed (b2/2)") {
					return err
				}
				len -= len2
				tlen -= len2
			}
		}
	}
}
