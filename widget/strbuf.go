package widget

import "github.com/burl/termbox-go"

// StrBuf is the basis for an input field widget
type StrBuf struct {
	Buf  string
	Pos  int
	Row  int
	Col  int
	Xoff int
	Yoff int
}

// NewStrBuf - create a new StrBuf at x,y in the screen
func NewStrBuf(x, y, xoff, yoff int) *StrBuf {
	return &StrBuf{
		Buf:  "",
		Pos:  0,
		Row:  y,
		Col:  x,
		Xoff: xoff,
		Yoff: yoff,
	}
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
	for i := b.Pos; i < len(b.Buf); i++ {
		termbox.SetCell(b.Col+i+2, b.Row, rune(b.Buf[i]), termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.SetCell(b.Col+b.Pos+1, b.Row, ch, termbox.ColorDefault, termbox.ColorDefault)
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
	for i := b.Pos; i < len(b.Buf); i++ {
		termbox.SetCell(b.Col+i, b.Row, rune(b.Buf[i]), termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.SetCell(b.Col+len(b.Buf), b.Row, '\x20', termbox.ColorDefault, termbox.ColorDefault)
	b.Buf = b.Buf[0:b.Pos-1] + b.Buf[b.Pos:]
	b.Left()
}

// Left - move cursor left in the StrBuf
func (b *StrBuf) Left() {
	if b.Pos > 0 {
		b.Pos = b.Pos - 1
		termbox.SetCursor(b.Col+b.Pos+b.Xoff, b.Row+b.Yoff)
	}
}

// Beginning - move cursor/insertion point to beginning of StrBuf
func (b *StrBuf) Beginning() {
	if b.Pos > 0 {
		b.Pos = 0
		termbox.SetCursor(b.Col+b.Pos+b.Xoff, b.Row+b.Yoff)
	}
}

// Right - move cursor/insertion-point to the right in the StrBuf
func (b *StrBuf) Right() {
	if b.Pos < len(b.Buf) {
		b.Pos = b.Pos + 1
		termbox.SetCursor(b.Col+b.Pos+b.Xoff, b.Row+b.Yoff)
	}
}

// End - move cursor/insertion point to the end of the StrBuf
func (b *StrBuf) End() {
	dx := len(b.Buf) - b.Pos
	if dx > 1 {
		b.Pos += dx
		termbox.SetCursor(b.Col+b.Pos+b.Xoff, b.Row+b.Yoff)
	}
}
