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
	name     string
	prompt   string
	hint     string
	lines    int
	row      int
	hintCol  int
	rightCol int
}

// Init - initialize any widget
func (w *WBase) Init() {
}

// Lines - how tall is this widget?
func (w *WBase) Lines() int {
	return w.lines
}

// SetRow - set row
func (w *WBase) SetRow(row int) {
	w.row = row
}

// DoValidate validation routine
func (w *WBase) DoValidate() (msg string) {
	return ""
}

// Renderable interface
type Renderable interface {
	Render(flush func())
	Lines() int
	SetRow(int)
	DoValidate() string
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
			ch = rune(str[stri])
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
