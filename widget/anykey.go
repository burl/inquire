package widget

import (
	"context"

	"github.com/burl/inquire/v2/internal/termui"
)

// AnyKey waits for any key before continuing.
type AnyKey struct {
	Base
	message string
	hint    string
}

// NewAnyKey constructs an AnyKey prompt.
func NewAnyKey(message string) *AnyKey {
	return &AnyKey{
		message: message,
		hint:    "press any key",
	}
}

// Hint sets hint text shown beside the message.
func (w *AnyKey) Hint(h string) *AnyKey {
	w.hint = h
	return w
}

// When registers a skip predicate.
func (w *AnyKey) When(fn func() bool) *AnyKey {
	w.Base.When(fn)
	return w
}

// Run displays the message and continues on any key.
func (w *AnyKey) Run(ctx context.Context, scr *termui.Screen) error {
	band, err := scr.OpenBand(ctx, 1)
	if err != nil {
		return err
	}

	draw := func() {
		band.Clear()
		band.WriteString(0, 0, "? ", stylePrompt)
		x := 2
		x += writeStyled(band, x, 0, w.message, styleQuestion)
		if w.hint != "" {
			_ = writeStyled(band, x, 0, " ("+w.hint+")", styleHint)
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

		switch ev.Key {
		case termui.KeyCtrlC:
			return termui.ErrInterrupted
		default:
			band.Clear()
			drawSettledRow(band, 0, w.message, "", false, 0)
			return band.FinalizeStatic(1)
		}
	}
}