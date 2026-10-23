package termui

import (
	"errors"
	"io"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func testPipes(t *testing.T) (inR, inW, outW *os.File) {
	t.Helper()
	inR, inW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	go func() { _, _ = io.Copy(io.Discard, outR) }()
	t.Cleanup(func() {
		_ = inR.Close()
		_ = inW.Close()
		_ = outW.Close()
		_ = outR.Close()
	})
	return inR, inW, outW
}

// A closed Screen's read loop must be gone: input that arrives afterwards
// belongs to whoever reads next.
func TestCloseStopsReader(t *testing.T) {
	inR, inW, outW := testPipes(t)
	s := newTestScreen(inR, outW, 80, 24)
	closed := make(chan struct{})
	go func() { _ = s.Close(); close(closed) }()
	select {
	case <-closed:
	case <-time.After(2 * time.Second):
		t.Fatal("Close did not return")
	}

	if _, err := inW.Write([]byte("x")); err != nil {
		t.Fatal(err)
	}
	got := make(chan byte, 1)
	go func() {
		b := make([]byte, 1)
		if n, _ := inR.Read(b); n == 1 {
			got <- b[0]
		}
	}()
	select {
	case b := <-got:
		if b != 'x' {
			t.Fatalf("read %q, want x", b)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("input after Close was consumed by the closed Screen's reader")
	}
}

// The second of two Screens on one input gets its cursor position report.
func TestSecondScreenGetsCursorReport(t *testing.T) {
	inR, inW, outW := testPipes(t)
	first := newTestScreen(inR, outW, 80, 24)
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second := newTestScreen(inR, outW, 80, 24)
	defer func() { _ = second.Close() }()
	go func() {
		time.Sleep(50 * time.Millisecond)
		_, _ = inW.Write(CursorReport(17, 1))
	}()
	row, _, err := second.CursorPos(t.Context())
	if err != nil {
		t.Fatalf("CursorPos: %v", err)
	}
	if row != 17 {
		t.Fatalf("row %d, want 17", row)
	}
}

// withSelect replaces selectFn for one test.
func withSelect(t *testing.T, fn func(int, *unix.FdSet, *unix.FdSet, *unix.FdSet, *unix.Timeval) (int, error)) {
	t.Helper()
	orig := selectFn
	selectFn = fn
	t.Cleanup(func() { selectFn = orig })
}

// pollKeys reads n key events from s, failing on error or timeout.
func pollKeys(t *testing.T, s *Screen, n int) []rune {
	t.Helper()
	type res struct {
		r   []rune
		err error
	}
	done := make(chan res, 1)
	go func() {
		var out []rune
		for len(out) < n {
			ev, err := s.Poll(t.Context())
			if err != nil {
				done <- res{out, err}
				return
			}
			if ev.Type == EventKey {
				out = append(out, ev.Rune)
			}
		}
		done <- res{out, nil}
	}()
	select {
	case r := <-done:
		if r.err != nil {
			t.Fatalf("Poll after %d key(s): %v", len(r.r), r.err)
		}
		return r.r
	case <-time.After(5 * time.Second):
		t.Fatal("no key events: the reader stopped reading")
		return nil
	}
}

// EINTR leaves the fd sets as passed in (wake bit included); the reader must
// retry, not take it for Close.
func TestReaderSurvivesEINTR(t *testing.T) {
	var calls atomic.Int32
	withSelect(t, func(nfd int, r, w, e *unix.FdSet, tv *unix.Timeval) (int, error) {
		if calls.Add(1) == 1 {
			return -1, unix.EINTR
		}
		return unix.Select(nfd, r, w, e, tv)
	})
	inR, inW, outW := testPipes(t)
	s := newTestScreen(inR, outW, 80, 24)
	defer func() { _ = s.Close() }()
	if _, err := inW.Write([]byte("ok")); err != nil {
		t.Fatal(err)
	}
	if got := string(pollKeys(t, s, 2)); got != "ok" {
		t.Fatalf("keys %q, want ok", got)
	}
	if calls.Load() < 2 {
		t.Fatal("select was not retried after EINTR")
	}
}

// Real signals arriving while typing must not stop the reader.
func TestReaderSurvivesSignals(t *testing.T) {
	ch := make(chan os.Signal, 64)
	signal.Notify(ch, syscall.SIGWINCH)
	defer signal.Stop(ch)
	go func() {
		for range ch {
		}
	}()

	inR, inW, outW := testPipes(t)
	s := newTestScreen(inR, outW, 80, 24)
	defer func() { _ = s.Close() }()

	const n = 100
	stop := make(chan struct{})
	defer close(stop)
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				_ = syscall.Kill(os.Getpid(), syscall.SIGWINCH)
				time.Sleep(100 * time.Microsecond)
			}
		}
	}()
	go func() {
		for range n {
			_, _ = inW.Write([]byte("x"))
			time.Sleep(time.Millisecond)
		}
	}()
	if got := len(pollKeys(t, s, n)); got != n {
		t.Fatalf("got %d keys, want %d", got, n)
	}
}

// A reader that cannot go on reports it; Poll must not wait forever.
func TestReaderDeathIsAnError(t *testing.T) {
	withSelect(t, func(int, *unix.FdSet, *unix.FdSet, *unix.FdSet, *unix.Timeval) (int, error) {
		return -1, unix.EBADF
	})
	inR, _, outW := testPipes(t)
	s := newTestScreen(inR, outW, 80, 24)
	defer func() { _ = s.Close() }()
	done := make(chan error, 1)
	go func() {
		_, err := s.Poll(t.Context())
		done <- err
	}()
	select {
	case err := <-done:
		if !errors.Is(err, unix.EBADF) {
			t.Fatalf("Poll error %v, want EBADF", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Poll hung after the reader failed")
	}
}
