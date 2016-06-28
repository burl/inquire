package widget

import (
	"fmt"

	"github.com/burl/termbox-go"
)

// TODO:
//  * support "delete" (right now, its just backspace)
//  * support CTRL+K, kill to end of line
//  * forward/backward one word
//  * up/down for history?
//

// ------------------------------------------------------------- Input Widget --

// Input an input widget, single line/string of text
type Input struct {
	WBase
	defaultValue string
	Value        string
	bound        *string
	validator    func(str string) (msg string)
}

// WInput - create new input widget
func WInput(prompt, dfl string) *Input {
	w := &Input{
		WBase:        WBase{prompt: prompt, lines: 2, hint: dfl},
		defaultValue: dfl,
	}
	w.Init()
	return w
}

// InputStringVar - create new input widget
func InputStringVar(value *string, prompt, dfl string) *Input {
	w := &Input{
		WBase:        WBase{prompt: prompt, lines: 2, hint: dfl},
		defaultValue: dfl,
		bound:        value,
	}
	w.Init()
	return w
}

// Default - set default value
func (w *Input) Default(val string) *Input {
	w.defaultValue = val
	w.hint = val
	fmt.Printf("*** set default value to %s\n", w.defaultValue)
	return w
}

// Validate - validation for input
func (w *Input) Validate(validator func(str string) (msg string)) *Input {
	w.validator = validator
	return w
}

// DoValidate validation routine
func (w *Input) DoValidate() (msg string) {
	if w.validator != nil {
		msg = w.validator(w.Value)
	}
	return
}

// Render WProtoT widget
func (w *Input) Render(flush func()) {
	if w.defaultValue == "" && w.bound != nil {
		w.defaultValue = *w.bound
		w.hint = *w.bound
	}
	w.drawPrompt()

	row := w.row
	col := w.rightCol - 1
	fmt.Print("\x1b[?25h")       // ShowCursor
	defer fmt.Print("\x1b[?25l") // HideCursor

	buf := NewStrBuf(col, 0, 1, row-1)
	buf.Append('\x20')
	buf.Delete()
	flush()

	quit := false
	hasError := false
	for {
		doFlush := true
		ev := termbox.PollEvent()
		if hasError {
			w.ErrorClear()
			hasError = false
		}
		switch ev.Type {
		case termbox.EventKey:
			if ev.Key == 0 {
				buf.Insert(ev.Ch)
			} else {
				switch ev.Key {
				case termbox.KeySpace:
					buf.Insert(' ')
				case termbox.KeyTab:
					buf.Insert(' ')
					buf.Insert(' ')
					buf.Insert(' ')
					buf.Insert(' ')
				case 5: // ^E
					buf.End()
				case 1: // ^A
					buf.Beginning()
				case termbox.KeyArrowUp:
				case termbox.KeyArrowDown:
				case termbox.KeyArrowLeft:
					buf.Left()
				case termbox.KeyArrowRight:
					buf.Right()
				case termbox.KeyDelete, termbox.KeyBackspace, termbox.KeyBackspace2:
					buf.Delete()
				case 3:
					dieFromCtlc()
				case 13:
					if buf.Buf == "" && w.defaultValue != "" {
						buf.SetValue(w.defaultValue)
						flush()
					}
					w.Value = buf.Buf
					msg := w.DoValidate()
					if msg == "" {
						if w.bound != nil {
							*w.bound = buf.Buf
						}
						quit = true
					} else {
						w.Value = ""
						w.ErrorMessage(msg)
						buf.End()
						hasError = true
						flush()
					}
				default:
					doFlush = false
				}
			}
		default:
			doFlush = false
		}
		if doFlush {
			flush()
		}
		if quit {
			break
		}
	}

	w.drawResult(buf.Buf)
	flush()
}
