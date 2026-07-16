package termui

import (
	"context"

	"github.com/mattn/go-runewidth"
)

type cell struct {
	Ch rune
	St Style
}

// Band is an N-line paint region anchored at a screen row (no alt-screen).
type Band struct {
	screen *Screen
	lines  int
	cols   int
	origin int // 1-based row of top line
	cells  [][]cell
	closed bool
}

// Lines returns the band height.
func (b *Band) Lines() int { return b.lines }

// Cols returns the band width (terminal columns at last layout).
func (b *Band) Cols() int { return b.cols }

// OriginRow is the 1-based absolute row of the top of the band (if known).
func (b *Band) OriginRow() int { return b.origin }

// Clear fills the band with spaces.
func (b *Band) Clear() {
	for y := 0; y < b.lines; y++ {
		for x := 0; x < b.cols; x++ {
			b.cells[y][x] = cell{Ch: ' '}
		}
	}
}

// SetCell writes a rune at band-relative (x,y). Wide runes consume extra cells.
func (b *Band) SetCell(x, y int, r rune, st Style) {
	if y < 0 || y >= b.lines || x < 0 || x >= b.cols {
		return
	}
	w := runewidth.RuneWidth(r)
	if w <= 0 {
		w = 1
	}
	b.cells[y][x] = cell{Ch: r, St: st}
	// mark continuation cells so Flush doesn't redraw junk mid-rune
	for i := 1; i < w && x+i < b.cols; i++ {
		b.cells[y][x+i] = cell{Ch: 0} // 0 = wide-tail
	}
}

// WriteString paints s at (x,y) with style; stops at band edge.
func (b *Band) WriteString(x, y int, s string, st Style) {
	for _, r := range s {
		if x >= b.cols {
			return
		}
		w := runewidth.RuneWidth(r)
		if w <= 0 {
			w = 1
		}
		b.SetCell(x, y, r, st)
		x += w
	}
}

// ResizeCols adjusts the band width after a terminal resize (keeps height).
func (b *Band) ResizeCols(cols int) {
	if cols < 1 {
		cols = 1
	}
	if cols == b.cols {
		return
	}
	old := b.cells
	b.cols = cols
	b.cells = make([][]cell, b.lines)
	for y := 0; y < b.lines; y++ {
		b.cells[y] = make([]cell, cols)
		for x := 0; x < cols; x++ {
			if y < len(old) && x < len(old[y]) {
				b.cells[y][x] = old[y][x]
			} else {
				b.cells[y][x] = cell{Ch: ' '}
			}
		}
	}
}

// Reanchor updates the band origin from the current cursor position (CSI 6n).
// Call after resize; assumes the cursor is parked on the bottom row of the band.
func (b *Band) Reanchor(ctx context.Context) error {
	if b.closed {
		return ErrClosed
	}
	row, _, err := b.screen.CursorPos(ctx)
	if err != nil {
		return err
	}
	origin := row - b.lines + 1
	if origin < 1 {
		origin = 1
	}
	b.origin = origin
	return nil
}

// OnResize re-anchors the band and adjusts width after SIGWINCH.
func (b *Band) OnResize(ctx context.Context, cols int) error {
	oldOrigin := b.origin
	if err := b.Reanchor(ctx); err != nil {
		// If DSR fails mid-session, keep the old origin and still reflow width.
		b.ResizeCols(cols)
		return nil
	}
	if oldOrigin != b.origin {
		for y := 0; y < b.lines; y++ {
			row := oldOrigin + y
			if err := cup(b.screen.out, row, 1); err != nil {
				return err
			}
			if err := clearLine(b.screen.out); err != nil {
				return err
			}
		}
	}
	b.ResizeCols(cols)
	return nil
}

func cellHasPaint(c cell) bool {
	if c.Ch == 0 {
		return false
	}
	if c.Ch != ' ' {
		return true
	}
	return c.St.Rev || c.St.Bold || c.St.Faint || c.St.Fg != ColorDefault
}

func (b *Band) lineLastContent(xCells []cell) int {
	last := -1
	for x := 0; x < len(xCells); x++ {
		c := xCells[x]
		if !cellHasPaint(c) {
			continue
		}
		w := runewidth.RuneWidth(c.Ch)
		if w <= 0 {
			w = 1
		}
		last = x + w - 1
	}
	return last
}

func (b *Band) flushLine(y int) error {
	if b.closed {
		return ErrClosed
	}
	s := b.screen
	row := b.origin + y
	if err := cup(s.out, row, 1); err != nil {
		return err
	}
	if err := clearLine(s.out); err != nil {
		return err
	}

	last := b.lineLastContent(b.cells[y])
	if last < 0 {
		return nil
	}

	if err := cup(s.out, row, 1); err != nil {
		return err
	}
	x := 0
	for x < b.cols && x <= last {
		c := b.cells[y][x]
		if c.Ch == 0 {
			x++
			continue
		}
		ch := c.Ch
		if err := sgr0(s.out); err != nil {
			return err
		}
		if err := c.St.apply(s.out, s.colorEnabled); err != nil {
			return err
		}
		if _, err := s.out.WriteString(string(ch)); err != nil {
			return err
		}
		w := runewidth.RuneWidth(ch)
		if w <= 0 {
			w = 1
		}
		x += w
	}
	return sgr0(s.out)
}

func (b *Band) eraseTerminalLine(y int) error {
	row := b.origin + y
	if err := cup(b.screen.out, row, 1); err != nil {
		return err
	}
	return clearLine(b.screen.out)
}

// Flush paints the entire band to the terminal at its origin.
func (b *Band) Flush() error {
	if b.closed {
		return ErrClosed
	}
	for y := 0; y < b.lines; y++ {
		if err := b.flushLine(y); err != nil {
			return err
		}
	}
	return cup(b.screen.out, b.origin+b.lines-1, 1)
}

// FlushLines paints only the first n band rows.
func (b *Band) FlushLines(n int) error {
	if n < 0 {
		n = 0
	}
	if n > b.lines {
		n = b.lines
	}
	for y := 0; y < n; y++ {
		if err := b.flushLine(y); err != nil {
			return err
		}
	}
	if n == 0 {
		return nil
	}
	return cup(b.screen.out, b.origin+n-1, 1)
}

// Close parks the cursor on the line below the band. Does not clear painted cells.
func (b *Band) Close() error {
	return b.closeAfter(b.lines)
}

func (b *Band) closeAfter(visibleLines int) error {
	if b.closed {
		return nil
	}
	if visibleLines < 1 {
		visibleLines = 1
	}
	b.closed = true
	// Newline from the last settled row commits it to terminal scrollback.
	// Parking on the cleared line below (old behavior) left CUP-painted rows
	// uncommitted; a later OpenBand scroll could erase them from history.
	lastRow := b.origin + visibleLines - 1
	if err := cup(b.screen.out, lastRow, 1); err != nil {
		return err
	}
	return newline(b.screen.out)
}

// FinalizeStatic settles the band to keepLines of content, erases any extra
// interactive rows on the terminal, commits the settled rows to scrollback,
// and parks the cursor on the line below for the next OpenBand.
func (b *Band) FinalizeStatic(keepLines int) error {
	if keepLines < 1 {
		keepLines = 1
	}
	if keepLines > b.lines {
		keepLines = b.lines
	}
	if err := b.FlushLines(keepLines); err != nil {
		return err
	}
	for y := keepLines; y < b.lines; y++ {
		if err := b.eraseTerminalLine(y); err != nil {
			return err
		}
	}
	return b.closeAfter(keepLines)
}
