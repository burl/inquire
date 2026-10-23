package termui

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// The read loop waits in select(2) on the input fd and a wake pipe, and
// reads only when input is ready. A terminal's stdin is a blocking fd outside
// Go's netpoller, so read deadlines don't apply to it; a loop parked in Read
// would outlive Close and steal the next Screen's input (cursor position
// reports, keys). select, unlike poll on macOS, works with tty devices.
//
// A signal interrupts select with EINTR and leaves the fd sets as they were
// passed in, so the sets are read only after a successful select. The loop
// returns quietly only on Close; any other exit delivers an error, so Poll
// fails instead of waiting forever with nobody reading the terminal.

// selectFn is unix.Select; tests replace it.
var selectFn = unix.Select

func (s *Screen) startReader() error {
	r, w, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("termui: wake pipe: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.wakeR, s.wakeW = r, w
	s.stopRead = cancel
	s.readDone = make(chan struct{})
	wakeFd := int(r.Fd())
	go s.readLoop(ctx, wakeFd)
	return nil
}

// stopReader stops the read loop and waits until it has returned, so no
// input is consumed after Close.
func (s *Screen) stopReader() {
	if s.stopRead == nil {
		return
	}
	s.stopRead()
	_, _ = s.wakeW.Write([]byte{0})
	<-s.readDone
	_ = s.wakeW.Close()
	_ = s.wakeR.Close()
	s.stopRead = nil
}

func (s *Screen) readLoop(ctx context.Context, wakeFd int) {
	defer close(s.readDone)
	buf := make([]byte, 256)
	nfd := max(s.fdIn, wakeFd) + 1
	for {
		var fds unix.FdSet
		fds.Set(s.fdIn)
		fds.Set(wakeFd)
		_, err := selectFn(nfd, &fds, nil, nil, nil)
		if errors.Is(err, unix.EINTR) {
			if ctx.Err() != nil {
				return
			}
			continue
		}
		if err != nil {
			s.deliver(ctx, readResult{err: fmt.Errorf("termui: select: %w", err)})
			return
		}
		if ctx.Err() != nil {
			return
		}
		if fds.IsSet(wakeFd) {
			// Only Close writes the wake byte, after cancelling ctx.
			s.deliver(ctx, readResult{err: ErrClosed})
			return
		}
		if !fds.IsSet(s.fdIn) {
			continue
		}
		n, err := s.in.Read(buf)
		if n > 0 {
			cp := make([]byte, n)
			copy(cp, buf[:n])
			if !s.deliver(ctx, readResult{data: cp}) {
				return
			}
		}
		if err != nil {
			s.deliver(ctx, readResult{err: err})
			return
		}
	}
}

// deliver hands r to the Screen; false once the Screen is closing.
func (s *Screen) deliver(ctx context.Context, r readResult) bool {
	select {
	case s.readCh <- r:
		return true
	case <-ctx.Done():
		return false
	}
}
