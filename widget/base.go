package widget

import "github.com/burl/termbox-go"

// chars to use for display
const (
	CharQuestion          = '?'
	CharCheckMark         = '✔'
	CharXMark             = '✖'
	CharSpace             = '\x20'
	CharSmallSquareFilled = '◼'
	CharSmallSquare       = '◻'
	CharChevronRight      = '❯'
	CharCircle            = '◯'
	CharCircleFilled      = '◉'
	CharBullseye          = '◎'
	CharBallotX           = '✗'
	CharBallotXHeavy      = '✘'
)

// WBase base widget
type WBase struct {
	name      string
	prompt    string
	hint      string
	lines     int
	row       int
	hintCol   int
	rightCol  int
	when      func() bool
	isMasked  bool
	inputMask rune
}

// Renderable interface
type Renderable interface {
	Render(flush func())
	Lines() int
	SetRow(int)
	DoValidate() string
	DoWhen() bool
}

// Lines - how tall is this widget?
func (w *WBase) Lines() int {
	return w.lines
}

// Hint - set the hint for the display
func (w *WBase) Hint(hint string) *WBase {
	w.hint = hint
	return w
}

// ClearHint - remove hint from widget display
func (w *WBase) ClearHint() *WBase {
	tbClearToEOL(w.hintCol, 0)
	return w
}

// Init - initialize any widget
func (w *WBase) Init() {
}

// When - set when
func (w *WBase) When(when func() bool) *WBase {
	w.when = when
	return w
}

// WhenEqual - shortcut for oft-used when logic
func (w *WBase) WhenEqual(a, b interface{}) *WBase {
	w.When(func() bool {
		switch v := a.(type) {
		case *string:
			a = *v
		case *bool:
			a = *v
		}
		switch v := b.(type) {
		case *string:
			b = *v
		case *bool:
			b = *v
		}
		if a == b {
			return true
		}
		return false
	})
	return w
}

// DoWhen - return result of when callback, or true
func (w *WBase) DoWhen() bool {
	if w.when != nil {
		return w.when()
	}
	return true
}

// SetRow - set row
func (w *WBase) SetRow(row int) {
	w.row = row
}

// DoValidate validation routine
func (w *WBase) DoValidate() (msg string) {
	return ""
}

func (w *WBase) drawPrompt() {
	tbPrint(0, 0, termbox.ColorGreen|termbox.AttrBold, coldef, string(CharQuestion))
	w.hintCol = tbPrint(2, 0, termbox.AttrBold|coldef, coldef, w.prompt+"?")
	w.rightCol = w.hintCol
	if w.hint != "" {
		w.rightCol = tbPrint(w.rightCol+1, 0, coldef, coldef, "("+w.hint+")")
	}
	w.rightCol = tbPrint(w.rightCol, 0, coldef, coldef, string(CharSpace))
}

func (w *WBase) drawResult(str string) {
	cols, _ := termbox.Size()
	stri := 0
	strl := len(str)
	w.hintCol = tbPrint(2, 0, coldef, coldef, w.prompt+"?")
	termbox.SetCell(0, 0, CharCheckMark, termbox.ColorGreen, coldef)
	for c := w.hintCol + 1; c < cols; c++ {
		var ch rune
		if stri < strl {
			if w.isMasked {
				ch = w.inputMask
			} else {
				ch = rune(str[stri])
			}
		} else {
			ch = '\x20'
		}
		stri = stri + 1
		// fmt.Printf("%d,%d : %c\n", c, w.row-1, ch)
		termbox.SetCell(c, 0, ch, termbox.ColorCyan, coldef)
	}
}

// ErrorMessage - show an error message WRT validation failure
// TODO: this needs to be attached to the WBaseT type
//       and it should be integrated into the "bottom bar" goroutine
//
func (w *WBase) ErrorMessage(msg string) {
	tbPrint(0, 0, termbox.ColorRed, coldef, string(CharXMark))
	tbPrint(2, w.lines-1, coldef, coldef, "error: "+msg)
}

// ErrorClear - clear error message
func (w *WBase) ErrorClear() {
	tbPrint(0, 0, termbox.ColorGreen, coldef, string(CharQuestion))
	tbClearToEOL(0, w.lines-1)
}
