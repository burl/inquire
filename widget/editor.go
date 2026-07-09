package widget

import (
	"unicode"

	"github.com/burl/inquire/v2/internal/termui"
	"github.com/mattn/go-runewidth"
)

// Editor is a rune-safe single-line text editor.
type Editor struct {
	runes []rune
	pos   int
	mask  rune
}

// NewEditor returns an empty editor.
func NewEditor() *Editor {
	return &Editor{}
}

// String returns the edited text.
func (e *Editor) String() string {
	return string(e.runes)
}

// SetString replaces the buffer and parks the cursor at the end.
func (e *Editor) SetString(s string) {
	e.runes = []rune(s)
	e.pos = len(e.runes)
}

// Insert adds a rune at the cursor.
func (e *Editor) Insert(r rune) {
	e.runes = append(e.runes[:e.pos], append([]rune{r}, e.runes[e.pos:]...)...)
	e.pos++
}

// Backspace removes the rune before the cursor.
func (e *Editor) Backspace() {
	if e.pos == 0 {
		return
	}
	e.runes = append(e.runes[:e.pos-1], e.runes[e.pos:]...)
	e.pos--
}

// DeleteForward removes the rune at the cursor.
func (e *Editor) DeleteForward() {
	if e.pos >= len(e.runes) {
		return
	}
	e.runes = append(e.runes[:e.pos], e.runes[e.pos+1:]...)
}

// Left moves the cursor one rune left.
func (e *Editor) Left() {
	if e.pos > 0 {
		e.pos--
	}
}

// Right moves the cursor one rune right.
func (e *Editor) Right() {
	if e.pos < len(e.runes) {
		e.pos++
	}
}

// Home moves the cursor to the start.
func (e *Editor) Home() {
	e.pos = 0
}

// End moves the cursor to the end.
func (e *Editor) End() {
	e.pos = len(e.runes)
}

// KillToEnd removes text from the cursor through the end of the line.
func (e *Editor) KillToEnd() {
	e.runes = e.runes[:e.pos]
}

// KillWordBackward removes the word (readline-style) before the cursor.
func (e *Editor) KillWordBackward() {
	if e.pos == 0 {
		return
	}
	i := e.pos
	for i > 0 && unicode.IsSpace(e.runes[i-1]) {
		i--
	}
	if i == e.pos {
		for i > 0 && isWordRune(e.runes[i-1]) {
			i--
		}
	}
	if i == e.pos {
		i--
	}
	e.runes = append(e.runes[:i], e.runes[e.pos:]...)
	e.pos = i
}

func isWordRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// SetMask sets the echo character; zero disables masking.
func (e *Editor) SetMask(ch rune) {
	e.mask = ch
}

// CursorCol returns the display column of the cursor relative to startX.
func (e *Editor) CursorCol(startX int) int {
	w := 0
	for i := 0; i < e.pos; i++ {
		rw := runewidth.RuneWidth(e.runes[i])
		if rw <= 0 {
			rw = 1
		}
		w += rw
	}
	return startX + w
}

// Draw paints the buffer and a reverse-video cursor at (y, startX).
func (e *Editor) Draw(band *termui.Band, y, startX int) {
	x := startX
	for i := 0; i < len(e.runes); i++ {
		ch := e.runes[i]
		if e.mask != 0 {
			ch = e.mask
		}
		st := termui.Style{}
		if i == e.pos {
			st = styleCursor
		}
		rw := runewidth.RuneWidth(e.runes[i])
		if rw <= 0 {
			rw = 1
		}
		band.SetCell(x, y, ch, st)
		x += rw
	}
	if e.pos == len(e.runes) {
		band.SetCell(x, y, ' ', styleCursor)
	}
}
