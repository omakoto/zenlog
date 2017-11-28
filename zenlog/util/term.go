package util

// #include <unistd.h>
import "C"

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type window struct {
	row    uint16
	col    uint16
	xpixel uint16
	ypixel uint16
}

func Ttyname(fd uintptr) string {
	name, _ := C.ttyname(C.int(fd))

	return C.GoString(name)
}

func Tty() string {
	tty := Ttyname(os.Stdin.Fd())
	if tty == "" {
		tty = Ttyname(os.Stdout.Fd())
	}
	if tty == "" {
		tty = Ttyname(os.Stderr.Fd())
	}
	if tty == "" {
		Fatalf("Failed to infer TTY name.")
	}
	return tty
}

//func getTerminalSize(fd uintptr) (rows uint16, cols uint16, e error) {
//	w := new(window)
//	_, _, err := syscall.Syscall(syscall.SYS_IOCTL,
//		fd,
//		syscall.TIOCGWINSZ,
//		uintptr(unsafe.Pointer(w)),
//	)
//	if err != 0 {
//		rows = 0
//		cols = 0
//		e = fmt.Errorf("Ioctl failed. errno=%d", err)
//		return
//	}
//	rows = w.row;
//	cols = w.col;
//	e = nil
//	return
//}
//
//func setTerminalSize(fd uintptr, rows uint16, cols uint16) error {
//	w := new(window)
//	w.row = rows
//	w.col = cols
//	_, _, err := syscall.Syscall(syscall.SYS_IOCTL,
//		fd,
//		syscall.TIOCSWINSZ,
//		uintptr(unsafe.Pointer(w)),
//	)
//	if err != 0 {
//		return fmt.Errorf("Ioctl failed. errno=%d", err)
//	}
//	return nil
//}

func PropagateTerminalSize(from *os.File, to *os.File) error {
	w := new(window)
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		from.Fd(),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(w)),
	)
	if err != 0 {
		return fmt.Errorf("Ioctl(TIOCGWINSZ) failed. errno=%d", err)
	}
	_, _, err = syscall.Syscall(syscall.SYS_IOCTL,
		to.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(w)),
	)
	if err != 0 {
		return fmt.Errorf("Ioctl(TIOCSWINSZ) failed. errno=%d", err)
	}
	return nil
}