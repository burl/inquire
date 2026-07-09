package termui

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
)

// Fake is a headless Screen for tests: scripted stdin bytes and captured stdout.
type Fake struct {
	Screen *Screen

	mu   sync.Mutex
	out  bytes.Buffer
	done chan struct{}
}

// NewFake builds a Screen backed by a byte script and output capture.
func NewFake(cols, rows int, script []byte) (*Fake, error) {
	inR, inW, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	outR, outW, err := os.Pipe()
	if err != nil {
		_ = inR.Close()
		_ = inW.Close()
		return nil, err
	}

	f := &Fake{Screen: newTestScreen(inR, outW, cols, rows)}

	done := make(chan struct{})
	go func() {
		defer close(done)
		buf := make([]byte, 4096)
		for {
			n, err := outR.Read(buf)
			if n > 0 {
				f.mu.Lock()
				_, _ = f.out.Write(buf[:n])
				f.mu.Unlock()
			}
			if err != nil {
				return
			}
		}
	}()

	if len(script) > 0 {
		if _, err := inW.Write(script); err != nil {
			_ = f.Close()
			return nil, err
		}
	}
	_ = inW.Close()

	f.done = done
	return f, nil
}

// Close shuts down the fake screen and finalizes output capture.
func (f *Fake) Close() error {
	err := f.Screen.Close()
	_ = f.Screen.out.Close()
	f.Output()
	return err
}

// Output returns captured terminal writes.
func (f *Fake) Output() string {
	if f.done != nil {
		<-f.done
		f.done = nil
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.out.String()
}

func newTestScreen(in, out *os.File, cols, rows int) *Screen {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Screen{
		in:           in,
		out:          out,
		fdIn:         int(in.Fd()),
		fdOut:        int(out.Fd()),
		cols:         cols,
		rows:         rows,
		colorEnabled: false,
		readCh:       make(chan readResult, 8),
		winchCh:      make(chan os.Signal, 1),
		stopRead:     cancel,
	}
	go s.readLoop(ctx)
	return s
}

// CursorReport encodes a DSR cursor-position response for scripted input.
func CursorReport(row, col int) []byte {
	return fmt.Appendf(nil, "\x1b[%d;%dR", row, col)
}

// EncodeKey returns terminal bytes for a logical key event.
func EncodeKey(key Key, r rune) []byte {
	switch key {
	case KeyEnter:
		return []byte{'\r'}
	case KeyBackspace:
		return []byte{0x7f}
	case KeyDelete:
		return []byte("\x1b[3~")
	case KeyTab:
		return []byte{'\t'}
	case KeySpace:
		return []byte{' '}
	case KeyUp:
		return []byte("\x1b[A")
	case KeyDown:
		return []byte("\x1b[B")
	case KeyLeft:
		return []byte("\x1b[D")
	case KeyRight:
		return []byte("\x1b[C")
	case KeyCtrlC:
		return []byte{0x03}
	case KeyCtrlA:
		return []byte{0x01}
	case KeyCtrlE:
		return []byte{0x05}
	case KeyCtrlK:
		return []byte{0x0b}
	case KeyCtrlD:
		return []byte{0x04}
	case KeyCtrlW:
		return []byte{0x17}
	case KeyRune:
		return []byte(string(r))
	default:
		return nil
	}
}

// Script builds a byte script from cursor reports, runes, and encoded keys.
func Script(parts ...any) []byte {
	var b []byte
	for _, p := range parts {
		switch v := p.(type) {
		case byte:
			b = append(b, v)
		case rune:
			b = append(b, string(v)...)
		case string:
			b = append(b, v...)
		case []byte:
			b = append(b, v...)
		case Key:
			b = append(b, EncodeKey(v, 0)...)
		default:
			panic(fmt.Sprintf("termui.Script: unsupported %T", p))
		}
	}
	return b
}
