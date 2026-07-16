package widget

import (
	"context"
	"strings"

	"github.com/burl/inquire/v2/internal/termui"
)

// Note is a non-interactive message printed before the next prompt.
type Note struct {
	Base
	text string
}

// NewNote constructs a Note with the given text.
// Use embedded newlines for multiple lines.
func NewNote(text string) *Note {
	return &Note{text: text}
}

// When registers a skip predicate.
func (w *Note) When(fn func() bool) *Note {
	w.Base.When(fn)
	return w
}

// Run prints the note and continues immediately.
func (w *Note) Run(ctx context.Context, scr *termui.Screen) error {
	lines := noteLines(w.text)
	band, err := scr.OpenBand(ctx, len(lines))
	if err != nil {
		return err
	}

	faint := termui.Style{Faint: true}
	for y, line := range lines {
		if y == 0 {
			band.WriteString(0, y, "› ", stylePrompt)
		}
		band.WriteString(2, y, line, faint)
	}
	if err := band.Flush(); err != nil {
		return err
	}
	return band.FinalizeStatic(len(lines))
}

func noteLines(text string) []string {
	text = strings.TrimRight(text, "\n")
	if text == "" {
		return []string{""}
	}
	return strings.Split(text, "\n")
}