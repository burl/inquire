package termui

import (
	"io"
	"os"
	"testing"
	"time"
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
