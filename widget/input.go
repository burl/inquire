package widget

import (
	"context"

	"github.com/burl/inquire/v2/internal/termui"
)

// Input is a single-line text prompt.
type Input struct {
	Base
	value    *string
	prompt   string
	hint     string
	validate func(string) string
	mask     rune
}

// NewInput constructs an Input bound to value.
func NewInput(value *string, prompt string) *Input {
	w := &Input{value: value, prompt: prompt}
	if value != nil && *value != "" {
		w.hint = *value
	}
	return w
}

// Hint sets hint text shown beside the prompt.
func (w *Input) Hint(h string) *Input {
	w.hint = h
	return w
}

// Valid registers a validation callback; non-empty return is shown as an error.
func (w *Input) Valid(fn func(string) string) *Input {
	w.validate = fn
	return w
}

// MaskInput masks typed characters (password entry). Default mask is •.
func (w *Input) MaskInput(ch ...rune) *Input {
	if len(ch) == 0 {
		w.mask = '•'
	} else {
		w.mask = ch[0]
	}
	return w
}

// When registers a skip predicate.
func (w *Input) When(fn func() bool) *Input {
	w.Base.When(fn)
	return w
}

// Run interactively collects text input.
func (w *Input) Run(ctx context.Context, scr *termui.Screen) error {
	const lines = 2
	band, err := scr.OpenBand(ctx, lines)
	if err != nil {
		return err
	}

	defaultVal := ""
	if w.value != nil {
		defaultVal = *w.value
	}
	hint := w.hint
	if hint == "" {
		hint = defaultVal
	}

	ed := NewEditor()
	if w.mask != 0 {
		ed.SetMask(w.mask)
	}

	var valueCol int
	var errMsg string
	showHint := hint != ""

	draw := func() {
		band.Clear()
		h := hint
		if !showHint {
			h = ""
		}
		valueCol = drawPromptRow(band, 0, w.prompt, h)
		ed.Draw(band, 0, valueCol)
		if errMsg != "" {
			drawErrorRow(band, 1, errMsg)
		}
		_ = band.Flush()
	}
	draw()

	for {
		ev, err := PollKey(ctx, scr, band, draw)
		if err != nil {
			return err
		}
		if ev.Type != termui.EventKey {
			continue
		}

		if errMsg != "" {
			errMsg = ""
		}
		showHint = false

		switch ev.Key {
		case termui.KeyCtrlC:
			return termui.ErrInterrupted
		case termui.KeyEnter:
			val := ed.String()
			if val == "" && defaultVal != "" {
				val = defaultVal
				ed.SetString(val)
				draw()
			}
			if w.validate != nil {
				if msg := w.validate(val); msg != "" {
					errMsg = msg
					draw()
					continue
				}
			}
			if w.value != nil {
				*w.value = val
			}
			band.Clear()
			drawSettledRow(band, 0, w.prompt, val, w.mask != 0, w.mask)
			return band.FinalizeStatic(1)
		case termui.KeyBackspace:
			ed.Backspace()
			draw()
		case termui.KeyDelete:
			ed.DeleteForward()
			draw()
		case termui.KeyLeft:
			ed.Left()
			draw()
		case termui.KeyRight:
			ed.Right()
			draw()
		case termui.KeyHome, termui.KeyCtrlA:
			ed.Home()
			draw()
		case termui.KeyEnd, termui.KeyCtrlE:
			ed.End()
			draw()
		case termui.KeyCtrlK:
			ed.KillToEnd()
			draw()
		case termui.KeyCtrlD:
			ed.DeleteForward()
			draw()
		case termui.KeyCtrlW:
			ed.KillWordBackward()
			draw()
		case termui.KeyTab:
			for range 4 {
				ed.Insert(' ')
			}
			draw()
		case termui.KeySpace:
			ed.Insert(' ')
			draw()
		case termui.KeyRune:
			ed.Insert(ev.Rune)
			draw()
		}
	}
}
