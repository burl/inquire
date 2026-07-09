package widget

import (
	"context"

	"github.com/burl/inquire/v2/internal/termui"
)

// YesNo is a yes/no toggle prompt.
type YesNo struct {
	Base
	value  *bool
	prompt string
	hint   string
}

// NewYesNo constructs a YesNo bound to value.
func NewYesNo(value *bool, prompt string) *YesNo {
	return &YesNo{
		value:  value,
		prompt: prompt,
		hint:   "Yes/No",
	}
}

// Hint sets hint text shown beside the prompt.
func (w *YesNo) Hint(h string) *YesNo {
	w.hint = h
	return w
}

// When registers a skip predicate.
func (w *YesNo) When(fn func() bool) *YesNo {
	w.Base.When(fn)
	return w
}

// Run interactively collects a yes/no answer.
func (w *YesNo) Run(ctx context.Context, scr *termui.Screen) error {
	band, err := scr.OpenBand(ctx, 1)
	if err != nil {
		return err
	}

	answer := false
	if w.value != nil {
		answer = *w.value
	}

	showHint := w.hint != ""
	valueCol := 0

	draw := func() {
		band.Clear()
		h := w.hint
		if !showHint {
			h = ""
		}
		valueCol = drawPromptRow(band, 0, w.prompt, h)
		label := "No"
		if answer {
			label = "Yes"
		}
		used := writeStyled(band, valueCol, 0, label, styleAnswer)
		band.SetCell(valueCol+used, 0, ' ', styleCursor)
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
		showHint = false

		switch ev.Key {
		case termui.KeyCtrlC:
			return termui.ErrInterrupted
		case termui.KeyEnter:
			if w.value != nil {
				*w.value = answer
			}
			label := "No"
			if answer {
				label = "Yes"
			}
			band.Clear()
			drawSettledRow(band, 0, w.prompt, label, false, 0)
			return band.FinalizeStatic(1)
		case termui.KeyLeft, termui.KeyUp:
			answer = false
			draw()
		case termui.KeyRight, termui.KeyDown:
			answer = true
			draw()
		case termui.KeySpace, termui.KeyTab, termui.KeyBackspace, termui.KeyDelete:
			answer = !answer
			draw()
		case termui.KeyRune:
			switch ev.Rune {
			case 'y', 'Y':
				answer = true
				draw()
			case 'n', 'N':
				answer = false
				draw()
			}
		}
	}
}
