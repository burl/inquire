package termui

import (
	"bytes"
	"unicode/utf8"
)

// EventType classifies input events.
type EventType int

const (
	EventNone EventType = iota
	EventKey
	EventResize
	EventError
)

// Key is a logical key.
type Key int

const (
	KeyUnknown Key = iota
	KeyRune
	KeyEnter
	KeyBackspace
	KeyDelete
	KeyTab
	KeyEscape
	KeySpace
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyCtrlC
	KeyCtrlA
	KeyCtrlE
	KeyCtrlK
	KeyCtrlD
	KeyCtrlW
)

// Event is a single input or signal event.
type Event struct {
	Type EventType
	Key  Key
	Rune rune
	// Resize
	Cols int
	Rows int
	Err  error
}

// decoder turns a byte stream into Events. Escape sequences may be incomplete
// across reads; leftover bytes stay in buf.
type decoder struct {
	buf []byte
}

func (d *decoder) feed(p []byte) { d.buf = append(d.buf, p...) }

func (d *decoder) next() (Event, bool) {
	if len(d.buf) == 0 {
		return Event{}, false
	}

	// CSI / SS3 escape sequences
	if d.buf[0] == 0x1b {
		ev, n, ok, needMore := parseEsc(d.buf)
		if needMore {
			return Event{}, false
		}
		if ok {
			d.buf = d.buf[n:]
			return ev, true
		}
		// lone ESC
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyEscape}, true
	}

	b := d.buf[0]

	// Ctrl keys (ASCII)
	switch b {
	case 0x03: // Ctrl+C
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyCtrlC}, true
	case 0x01: // Ctrl+A
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyCtrlA}, true
	case 0x02: // Ctrl+B
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyLeft}, true
	case 0x05: // Ctrl+E
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyCtrlE}, true
	case 0x06: // Ctrl+F
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyRight}, true
	case 0x0b: // Ctrl+K
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyCtrlK}, true
	case 0x0e: // Ctrl+N
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyDown}, true
	case 0x10: // Ctrl+P
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyUp}, true
	case 0x17: // Ctrl+W
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyCtrlW}, true
	case 0x04: // Ctrl+D
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyCtrlD}, true
	case 0x0d, 0x0a: // CR / LF
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyEnter}, true
	case 0x09:
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyTab}, true
	case 0x7f, 0x08:
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyBackspace}, true
	case 0x20:
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeySpace, Rune: ' '}, true
	}

	// UTF-8 rune
	r, size := utf8.DecodeRune(d.buf)
	if r == utf8.RuneError && size == 1 {
		// incomplete multi-byte sequence?
		if !utf8.FullRune(d.buf) {
			return Event{}, false
		}
		d.buf = d.buf[1:]
		return Event{Type: EventKey, Key: KeyUnknown}, true
	}
	d.buf = d.buf[size:]
	return Event{Type: EventKey, Key: KeyRune, Rune: r}, true
}

// parseEsc parses ESC-led sequences. needMore means wait for more bytes.
func parseEsc(b []byte) (ev Event, n int, ok bool, needMore bool) {
	if len(b) < 2 {
		return Event{}, 0, false, true
	}

	// SS3: ESC O A/B/C/D (application cursor keys)
	if b[1] == 'O' {
		if len(b) < 3 {
			return Event{}, 0, false, true
		}
		switch b[2] {
		case 'A':
			return Event{Type: EventKey, Key: KeyUp}, 3, true, false
		case 'B':
			return Event{Type: EventKey, Key: KeyDown}, 3, true, false
		case 'C':
			return Event{Type: EventKey, Key: KeyRight}, 3, true, false
		case 'D':
			return Event{Type: EventKey, Key: KeyLeft}, 3, true, false
		case 'H':
			return Event{Type: EventKey, Key: KeyHome}, 3, true, false
		case 'F':
			return Event{Type: EventKey, Key: KeyEnd}, 3, true, false
		}
		return Event{}, 3, false, false
	}

	// CSI: ESC [
	if b[1] != '[' {
		return Event{}, 0, false, false
	}

	// scan to final byte (0x40–0x7E)
	i := 2
	for i < len(b) {
		c := b[i]
		if c >= 0x40 && c <= 0x7e {
			seq := b[2:i]
			final := c
			total := i + 1
			ev, ok := mapCSI(seq, final)
			return ev, total, ok, false
		}
		i++
	}
	// incomplete CSI — but protect against runaway garbage
	if len(b) > 64 {
		return Event{}, 2, false, false
	}
	return Event{}, 0, false, true
}

func mapCSI(params []byte, final byte) (Event, bool) {
	// Cursor position report: ESC [ row ; col R  — consumed by reader, not keys
	if final == 'R' {
		return Event{}, false
	}

	p := string(params)
	switch final {
	case 'A':
		return Event{Type: EventKey, Key: KeyUp}, true
	case 'B':
		return Event{Type: EventKey, Key: KeyDown}, true
	case 'C':
		return Event{Type: EventKey, Key: KeyRight}, true
	case 'D':
		return Event{Type: EventKey, Key: KeyLeft}, true
	case 'H':
		return Event{Type: EventKey, Key: KeyHome}, true
	case 'F':
		return Event{Type: EventKey, Key: KeyEnd}, true
	case '~':
		// numbered keys: 3~ = Delete, 1~ = Home, 4~ = End
		switch p {
		case "3":
			return Event{Type: EventKey, Key: KeyDelete}, true
		case "1", "7":
			return Event{Type: EventKey, Key: KeyHome}, true
		case "4", "8":
			return Event{Type: EventKey, Key: KeyEnd}, true
		}
	}
	return Event{}, false
}

// parseCursorReport looks for ESC [ row ; col R in b, returns remaining bytes.
func parseCursorReport(b []byte) (row, col int, rest []byte, ok bool) {
	// find ESC [
	idx := bytes.Index(b, []byte{0x1b, '['})
	if idx < 0 {
		return 0, 0, b, false
	}
	i := idx + 2
	row = 0
	for i < len(b) && b[i] >= '0' && b[i] <= '9' {
		row = row*10 + int(b[i]-'0')
		i++
	}
	if i >= len(b) || b[i] != ';' {
		return 0, 0, b, false
	}
	i++
	col = 0
	for i < len(b) && b[i] >= '0' && b[i] <= '9' {
		col = col*10 + int(b[i]-'0')
		i++
	}
	if i >= len(b) || b[i] != 'R' {
		// incomplete
		if len(b)-idx > 32 {
			// skip garbage ESC
			return 0, 0, b[idx+1:], false
		}
		return 0, 0, b, false
	}
	i++ // consume R
	rest = append([]byte{}, b[:idx]...)
	rest = append(rest, b[i:]...)
	return row, col, rest, true
}


