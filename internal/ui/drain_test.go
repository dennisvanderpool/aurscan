package ui

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

const (
	tiocsptlck = 0x40045431
	tiocgptn   = 0x80045430
	fionread   = 0x541B
)

func ioctlPtr(fd, req uintptr, arg unsafe.Pointer) syscall.Errno {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, req, uintptr(arg))
	return e
}

// openPTY returns a connected (master, slave) pty pair, or skips the test if
// the environment doesn't allow pty allocation.
func openPTY(t *testing.T) (*os.File, *os.File) {
	t.Helper()
	m, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		t.Skipf("no /dev/ptmx: %v", err)
	}
	var unlock int
	if e := ioctlPtr(m.Fd(), tiocsptlck, unsafe.Pointer(&unlock)); e != 0 {
		m.Close()
		t.Skipf("TIOCSPTLCK: %v", e)
	}
	var n uint32
	if e := ioctlPtr(m.Fd(), tiocgptn, unsafe.Pointer(&n)); e != 0 {
		m.Close()
		t.Skipf("TIOCGPTN: %v", e)
	}
	s, err := os.OpenFile(fmt.Sprintf("/dev/pts/%d", n), os.O_RDWR, 0)
	if err != nil {
		m.Close()
		t.Skipf("open slave: %v", err)
	}
	return m, s
}

func pending(t *testing.T, f *os.File) int {
	t.Helper()
	var n int
	if e := ioctlPtr(f.Fd(), fionread, unsafe.Pointer(&n)); e != 0 {
		t.Fatalf("FIONREAD: %v", e)
	}
	return n
}

func TestDrainInputClearsBufferedKeystrokes(t *testing.T) {
	master, slave := openPTY(t)
	defer master.Close()
	defer slave.Close()

	// Simulate the user mashing ENTER five times before the prompt appears.
	if _, err := master.WriteString("\n\n\n\n\n"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)

	if got := pending(t, slave); got == 0 {
		t.Skip("pty did not queue input on this system")
	}
	DrainInput(slave)
	if got := pending(t, slave); got != 0 {
		t.Fatalf("after DrainInput, %d bytes still queued — flush failed", got)
	}
}

func TestDrainInputNilAndNonTTYSafe(t *testing.T) {
	DrainInput(nil) // must not panic
	r, w, _ := os.Pipe()
	defer r.Close()
	defer w.Close()
	DrainInput(r) // not a tty: no-op, must not panic or block
}
