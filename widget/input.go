package widget

import "github.com/burl/termbox-go"

// TODO:
//  * support "delete" (right now, its just backspace)
//  * support CTRL+K, kill to end of line
//  * forward/backward one word
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
	w := WInput(prompt, dfl)
	w.bound = value
	return w
}

// Default - set default value
func (w *Input) Default(val string) *Input {
	w.defaultValue = val
	w.hint = val
	return w
}

// Valid - validation for input
func (w *Input) Valid(validator func(str string) (msg string)) *Input {
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

	// call this before NewStrBuf() because we need the column offset
	w.drawPrompt()

	buf := NewStrBuf(w.rightCol-1, 0)
	buf.Draw()
	flush()

	hasError := false
EventLoop:
	for {
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
					buf.Insert('\x20')
				case termbox.KeyTab:
					buf.Insert('\x20')
					buf.Insert('\x20')
					buf.Insert('\x20')
					buf.Insert('\x20')
				case 5: // ^E
					buf.End()
				case 1: // ^A
					buf.Beginning()
				case termbox.KeyArrowUp:
				case termbox.KeyArrowDown:
				case termbox.KeyArrowLeft, 02: // 02 == Ctrl+B
					buf.Left()
				case termbox.KeyArrowRight, 06: // 06 == Ctrl+F
					buf.Right()
				//TODO: backspace and delete / Ctrl+D, etc. are different.. make it so...
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
						break EventLoop
					} else {
						w.Value = ""
						w.ErrorMessage(msg)
						buf.End()
						hasError = true
						flush()
					}
				}
			}
		}
		flush()
	}

	w.drawResult(buf.Buf)
	flush()
}
