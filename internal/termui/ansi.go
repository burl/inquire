package termui

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Minimal ANSI helpers. No terminfo — modern terminals speak CSI.

const (
	esc = "\x1b"
	csi = "\x1b["
)

func writeString(w io.Writer, s string) error {
	_, err := io.WriteString(w, s)
	return err
}

func hideCursor(w io.Writer) error { return writeString(w, csi+"?25l") }
func showCursor(w io.Writer) error { return writeString(w, csi+"?25h") }
func clearLine(w io.Writer) error  { return writeString(w, csi+"2K") }
func clearToEOS(w io.Writer) error { return writeString(w, csi+"J") }
func newline(w io.Writer) error    { return writeString(w, "\n") }

// cup moves to 1-based row,col (absolute).
func cup(w io.Writer, row, col int) error {
	return writeString(w, fmt.Sprintf("%s%d;%dH", csi, row, col))
}

// cuu moves cursor up n lines.
func cuu(w io.Writer, n int) error {
	if n <= 0 {
		return nil
	}
	return writeString(w, fmt.Sprintf("%s%dA", csi, n))
}

// sgr0 resets attributes.
func sgr0(w io.Writer) error { return writeString(w, csi+"0m") }

// Style is a minimal cell style.
type Style struct {
	Fg    Color
	Bold  bool
	Rev   bool // reverse video
	Faint bool
}

// Color is a basic 8-color + default set.
type Color int

const (
	ColorDefault Color = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

func (s Style) apply(w io.Writer, colorEnabled bool) error {
	if !colorEnabled {
		if s.Rev {
			return writeString(w, csi+"7m")
		}
		return nil
	}
	var parts []string
	if s.Bold {
		parts = append(parts, "1")
	}
	if s.Faint {
		parts = append(parts, "2")
	}
	if s.Rev {
		parts = append(parts, "7")
	}
	if s.Fg != ColorDefault {
		// 30–37 foreground
		parts = append(parts, strconv.Itoa(29+int(s.Fg)))
	}
	if len(parts) == 0 {
		return nil
	}
	return writeString(w, csi+strings.Join(parts, ";")+"m")
}

func requestCursorPos(w io.Writer) error {
	// DSR: Device Status Report — cursor position
	return writeString(w, csi+"6n")
}
