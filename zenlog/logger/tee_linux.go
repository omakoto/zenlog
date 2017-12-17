// +build linux

package logger

//
//import (
//	"os"
//	"syscall"
//
//	"github.com/omakoto/zenlog-go/zenlog/util"
//	"golang.org/x/sys/unix"
//)
//
//func tee_faster(in, out1, out2 *os.File) error {
//	for {
//		len, err := syscall.Tee(int(in.Fd()), int(out1.Fd()), 1024*128, 0 /*unix.SPLICE_F_NONBLOCK*/)
//		if util.Warn(err, "Tee failed") {
//			return err
//		}
//		if len == 0 {
//			return nil
//		}
//
//		for len > 0 {
//			slen, err := syscall.Splice(int(in.Fd()), nil, int(out2.Fd()), nil, int(len), unix.SPLICE_F_MOVE)
//			if util.Warn(err, "Splice failed") {
//				return err
//			}
//			len -= slen
//		}
//	}
//}
