package widget

import (
	"strings"

	"github.com/burl/termbox-go"
)

// StrBuf is the basis for an input field widget
type StrBuf struct {
	Buf       string
	Pos       int
	Row       int
	Col       int
	isMasked  bool
	inputMask rune
}

// NewStrBuf - create a new StrBuf at x,y in the screen
func NewStrBuf(x, y int) *StrBuf {
	return &StrBuf{
		Buf: "",
		Pos: 0,
		Row: y,
		Col: x,
	}
}

// MaskInput - mask input (echo a mask char instead of actual char)
func (b *StrBuf) MaskInput(ch rune) {
	b.isMasked = true
	b.inputMask = ch
}

// Draw the buffer (not so performant.. could be optimized) - and draw
// a "cursor" by showing inverse...
//
// I previously turned on the real terminal cursor and positioned it, but
// there was a lot of extra accounting that had to take place.  This just
// draws an inverse character (or space) at the insertion point, it looks
// fine on most terminals and works with white backgrounds, etc.
//
func (b *StrBuf) Draw() {
	buf := b.Buf
	if b.isMasked {
		buf = strings.Repeat(string(b.inputMask), len(buf))
	}
	tbPrint(b.Col+1, b.Row, coldef, coldef, buf+"\x20\x20")
	ch := '\x20'
	if b.Pos < len(b.Buf) {
		if b.isMasked {
			ch = b.inputMask
		} else {
			ch = rune(b.Buf[b.Pos])
		}
	}
	termbox.SetCell(b.Pos+b.Col+1, 0, ch, termbox.AttrReverse, coldef)
}

// SetValue set the value of the StrBuf internal string
func (b *StrBuf) SetValue(value string) {
	for i := 0; i < len(value); i++ {
		b.Insert(rune(value[i]))
	}
}

// String - return value of strbuf
func (b *StrBuf) String() string {
	return b.Buf
}

// Append - add a character to the end of the StrBuf
func (b *StrBuf) Append(ch rune) {
	b.Buf = b.Buf + string(ch)
	b.Right()
}

// Insert a character at the current edit point in the StrBuf
func (b *StrBuf) Insert(ch rune) {
	b.Buf = b.Buf[0:b.Pos] + string(ch) + b.Buf[b.Pos:]
	b.Right()
}

// Delete a character at the current position in the StrBuf
func (b *StrBuf) Delete() {
	if len(b.Buf) < 1 {
		return
	}
	if b.Pos == 0 {
		return
	}
	b.Buf = b.Buf[0:b.Pos-1] + b.Buf[b.Pos:]
	b.Left()
}

// Left - move cursor left in the StrBuf
func (b *StrBuf) Left() {
	if b.Pos > 0 {
		b.Pos = b.Pos - 1
		b.Draw()
	}
}

// Beginning - move cursor/insertion point to beginning of StrBuf
func (b *StrBuf) Beginning() {
	if b.Pos > 0 {
		b.Pos = 0
		b.Draw()
	}
}

// Right - move cursor/insertion-point to the right in the StrBuf
func (b *StrBuf) Right() {
	if b.Pos < len(b.Buf) {
		b.Pos = b.Pos + 1
		b.Draw()
	}
}

// End - move cursor/insertion point to the end of the StrBuf
func (b *StrBuf) End() {
	dx := len(b.Buf) - b.Pos
	if dx > 1 {
		b.Pos += dx
		b.Draw()
	}
}
