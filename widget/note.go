package widget

import (
	"context"

	"github.com/burl/inquire/v2/internal/termui"
)

// Note is a non-interactive message shown before the next prompt.
type Note struct {
	Base
	text string
}

// NewNote constructs a Note with the given text.
func NewNote(text string) *Note {
	return &Note{text: text}
}

// When registers a skip predicate.
func (w *Note) When(fn func() bool) *Note {
	w.Base.When(fn)
	return w
}

// Run displays the note and waits for Enter.
func (w *Note) Run(ctx context.Context, scr *termui.Screen) error {
	band, err := scr.OpenBand(ctx, 1)
	if err != nil {
		return err
	}

	draw := func() {
		band.Clear()
		band.WriteString(0, 0, "› ", stylePrompt)
		band.WriteString(2, 0, w.text, termui.Style{Faint: true})
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

		switch ev.Key {
		case termui.KeyCtrlC:
			return termui.ErrInterrupted
		case termui.KeyEnter:
			band.Clear()
			band.WriteString(0, 0, "› ", stylePrompt)
			band.WriteString(2, 0, w.text, termui.Style{Faint: true})
			return band.FinalizeStatic(1)
		}
	}
}