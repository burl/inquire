package widget

import "github.com/burl/termbox-go"

// YesNo an input widget, single line/string of text
type YesNo struct {
	WBase
	boundVar *bool
}

// YesNoBoolVar - create new yes/no widget
func YesNoBoolVar(value *bool, prompt string) *YesNo {
	w := &YesNo{
		WBase:    WBase{prompt: prompt, lines: 2, hint: "Yes/No"},
		boundVar: value,
	}
	w.Init()
	return w
}

// Render YesNo widget
func (w *YesNo) Render(flush func()) {
	drawn := false
	answer := false
	result := "No"
	if w.boundVar != nil {
		answer = *w.boundVar
	}

	update := func() {
		if w.boundVar != nil {
			*w.boundVar = answer
		}
		if answer {
			result = "Yes"
		} else {
			result = "No"
		}
	}

	// draw a "phony cursor"
	drawCursor := func(x, y int) {
		tbPrint(x, y, termbox.ColorDefault|termbox.AttrReverse, termbox.ColorDefault, "\x20")
	}

	draw := func() {
		drawn = true
		update()
		tbPrint(w.rightCol, 0, termbox.ColorCyan, coldef, result+"\x20\x20")
		drawCursor(w.rightCol+len(result), 0)
	}

	w.drawPrompt()
	drawCursor(w.rightCol, 0)
	flush()

EventLoop:
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			if ev.Key == 0 {
				switch ev.Ch {
				case 'y', 'Y':
					answer = true
				case 'n', 'N':
					answer = false
				}
			} else {
				switch ev.Key {
				case termbox.KeyBackspace2,
					termbox.KeyBackspace,
					termbox.KeyDelete,
					termbox.KeySpace,
					termbox.KeyTab:
					answer = !answer
				case termbox.KeyArrowLeft:
					answer = false
				case termbox.KeyArrowRight:
					answer = true
				case 3:
					dieFromCtlc()
				case 13:
					if drawn {
						break EventLoop
					}
				}
			}
		}
		draw()
		flush()
	}

	update()
	w.drawResult(result)
	flush()
}
