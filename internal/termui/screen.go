package termui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/term"
)

// Sentinel errors for the spike / future library.
var (
	ErrNotTerminal = errors.New("inquire: stdin/stdout is not a terminal")
	ErrInterrupted = errors.New("inquire: interrupted")
	ErrClosed      = errors.New("inquire: screen closed")
)

// Screen owns raw mode and input demux (keys, cursor reports, resize).
type Screen struct {
	in  *os.File
	out *os.File

	fdIn  int
	fdOut int
	state *term.State

	mu     sync.Mutex
	cols   int
	rows   int
	closed bool

	colorEnabled bool

	// input
	dec     decoder
	pending []byte // raw bytes waiting for DSR parse / decode

	// read loop
	readCh   chan readResult
	winchCh  chan os.Signal
	stopRead context.CancelFunc
}

type readResult struct {
	data []byte
	err  error
}

// ScreenOption configures [OpenScreen].
type ScreenOption func(*screenConfig)

type screenConfig struct {
	color *bool
}

// WithColor forces ANSI color on (true) or off (false). Default follows env.
func WithColor(enabled bool) ScreenOption {
	return func(c *screenConfig) {
		c.color = &enabled
	}
}

// OpenScreen puts the terminal into raw mode for inline (no alt-screen) use.
func OpenScreen(in, out *os.File, opts ...ScreenOption) (*Screen, error) {
	if in == nil {
		in = os.Stdin
	}
	if out == nil {
		out = os.Stdout
	}
	fdIn := int(in.Fd())
	fdOut := int(out.Fd())
	if !term.IsTerminal(fdIn) || !term.IsTerminal(fdOut) {
		return nil, ErrNotTerminal
	}
	st, err := term.MakeRaw(fdIn)
	if err != nil {
		return nil, fmt.Errorf("termui: make raw: %w", err)
	}
	cols, rows, err := term.GetSize(fdOut)
	if err != nil {
		_ = term.Restore(fdIn, st)
		return nil, fmt.Errorf("termui: get size: %w", err)
	}

	var cfg screenConfig
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	colorEnabled := detectColor()
	if cfg.color != nil {
		colorEnabled = *cfg.color
	}

	s := &Screen{
		in:           in,
		out:          out,
		fdIn:         fdIn,
		fdOut:        fdOut,
		state:        st,
		cols:         cols,
		rows:         rows,
		colorEnabled: colorEnabled,
		readCh:       make(chan readResult, 8),
		winchCh:      make(chan os.Signal, 1),
	}

	if err := hideCursor(s.out); err != nil {
		_ = s.restore()
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.stopRead = cancel
	signal.Notify(s.winchCh, syscall.SIGWINCH)
	go s.readLoop(ctx)

	return s, nil
}

func detectColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	// Prefer COLORTERM / TERM hints; default on for real TTYs.
	if os.Getenv("COLORTERM") != "" {
		return true
	}
	termName := os.Getenv("TERM")
	if termName == "" || termName == "dumb" {
		return false
	}
	return true
}

func (s *Screen) restore() error {
	_ = showCursor(s.out)
	_ = sgr0(s.out)
	if s.state != nil {
		return term.Restore(s.fdIn, s.state)
	}
	return nil
}

// Close restores the terminal. Safe to call once.
func (s *Screen) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	s.mu.Unlock()

	if s.stopRead != nil {
		s.stopRead()
	}
	signal.Stop(s.winchCh)
	return s.restore()
}

// Size returns the last known terminal dimensions.
func (s *Screen) Size() (cols, rows int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cols, s.rows
}

// ColorEnabled reports whether ANSI colors will be emitted.
func (s *Screen) ColorEnabled() bool {
	return s.colorEnabled
}

func (s *Screen) readLoop(ctx context.Context) {
	buf := make([]byte, 256)
	for {
		// Set a short deadline so we can notice cancel.
		_ = s.in.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		n, err := s.in.Read(buf)
		if n > 0 {
			cp := make([]byte, n)
			copy(cp, buf[:n])
			select {
			case s.readCh <- readResult{data: cp}:
			case <-ctx.Done():
				return
			}
		}
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				select {
				case <-ctx.Done():
					return
				default:
					continue
				}
			}
			select {
			case s.readCh <- readResult{err: err}:
			case <-ctx.Done():
			}
			return
		}
	}
}

func (s *Screen) refreshSize() {
	cols, rows, err := term.GetSize(s.fdOut)
	if err != nil {
		return
	}
	s.mu.Lock()
	s.cols, s.rows = cols, rows
	s.mu.Unlock()
}

// CursorPos queries the terminal for the cursor position (1-based).
func (s *Screen) CursorPos(ctx context.Context) (row, col int, err error) {
	if err := requestCursorPos(s.out); err != nil {
		return 0, 0, err
	}
	deadline := time.Now().Add(500 * time.Millisecond)
	if dl, ok := ctx.Deadline(); ok && dl.Before(deadline) {
		deadline = dl
	}

	for {
		// first check pending
		if r, c, rest, ok := parseCursorReport(s.pending); ok {
			s.pending = rest
			return r, c, nil
		}
		if time.Now().After(deadline) {
			return 0, 0, fmt.Errorf("termui: timeout waiting for cursor position report")
		}
		remaining := time.Until(deadline)
		timer := time.NewTimer(remaining)
		select {
		case <-ctx.Done():
			timer.Stop()
			return 0, 0, ctx.Err()
		case sig := <-s.winchCh:
			timer.Stop()
			_ = sig
			s.refreshSize()
		case rr := <-s.readCh:
			timer.Stop()
			if rr.err != nil {
				return 0, 0, rr.err
			}
			s.pending = append(s.pending, rr.data...)
		case <-timer.C:
			return 0, 0, fmt.Errorf("termui: timeout waiting for cursor position report")
		}
	}
}

// Poll waits for the next key or resize event.
func (s *Screen) Poll(ctx context.Context) (Event, error) {
	for {
		if len(s.pending) > 0 {
			// strip any stray cursor reports
			if _, _, rest, ok := parseCursorReport(s.pending); ok {
				s.pending = rest
				continue
			}
			s.dec.feed(s.pending)
			s.pending = s.pending[:0]
		}
		if ev, ok := s.dec.next(); ok {
			if ev.Type == EventKey && ev.Key == KeyCtrlC {
				return ev, ErrInterrupted
			}
			return ev, nil
		}

		select {
		case <-ctx.Done():
			return Event{}, ctx.Err()
		case <-s.winchCh:
			s.refreshSize()
			cols, rows := s.Size()
			return Event{Type: EventResize, Cols: cols, Rows: rows}, nil
		case rr := <-s.readCh:
			if rr.err != nil {
				return Event{Type: EventError, Err: rr.err}, rr.err
			}
			s.pending = append(s.pending, rr.data...)
		}
	}
}

// OpenBand reserves `lines` rows at the current cursor for interactive painting.
func (s *Screen) OpenBand(ctx context.Context, lines int) (*Band, error) {
	if lines < 1 {
		return nil, fmt.Errorf("termui: band lines must be >= 1")
	}
	s.refreshSize()
	cols, rows := s.Size()

	// Ensure room: emit (lines-1) newlines then move back up and clear below.
	// This matches the original ViewPortSetHeight idea and scrolls if needed.
	for i := 0; i < lines-1; i++ {
		if err := newline(s.out); err != nil {
			return nil, err
		}
	}
	if lines > 1 {
		if err := cuu(s.out, lines-1); err != nil {
			return nil, err
		}
	}
	if err := clearToEOS(s.out); err != nil {
		return nil, err
	}

	row, col, err := s.CursorPos(ctx)
	if err != nil {
		// Fallback: assume top of reserved region is current logical origin.
		// Without DSR we cannot know absolute row; use 1 as last resort.
		row, col = 1, 1
		_ = col
	}

	// If DSR says we don't fit, we've already scrolled via newlines; trust row.
	_ = rows

	b := &Band{
		screen: s,
		lines:  lines,
		cols:   cols,
		origin: row, // 1-based
		cells:  make([][]cell, lines),
	}
	for y := 0; y < lines; y++ {
		b.cells[y] = make([]cell, cols)
		for x := 0; x < cols; x++ {
			b.cells[y][x] = cell{Ch: ' '}
		}
	}
	return b, nil
}
