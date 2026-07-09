// Package inquire provides line-oriented interactive CLI prompts that render
// inline at the current cursor, leaving answered questions as static scrollback.
//
// # TTY requirements
//
// Run requires both stdin and stdout to be terminals. Piping either stream
// returns [ErrNotTerminal]. Stderr may be redirected.
//
// # Interrupts
//
// Ctrl+C aborts the entire session and returns [ErrInterrupted]. Answers bound
// before the interrupt are kept; later prompts are not run.
//
// See https://pkg.go.dev/github.com/burl/inquire/v2 for API reference.
package inquire

import (
	"context"
	"errors"
	"os"

	"github.com/burl/inquire/v2/internal/termui"
	"github.com/burl/inquire/v2/widget"
)

var (
	// ErrNotTerminal is returned when stdin or stdout is not a terminal.
	ErrNotTerminal = errors.New("inquire: stdin/stdout is not a terminal")
	// ErrInterrupted is returned when the user presses Ctrl+C.
	ErrInterrupted = errors.New("inquire: interrupted")
)

// Query builds an interactive prompt session.
func Query(opts ...Option) *Session {
	s := &Session{
		in:  os.Stdin,
		out: os.Stdout,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(s)
		}
	}
	return s
}

// Session is a fluent builder for interactive prompts.
type Session struct {
	in           *os.File
	out          *os.File
	err          error
	widgets      []widget.Runner
	colorEnabled *bool
}

// Input adds a text prompt bound to value.
func (s *Session) Input(value *string, prompt string, more func(*widget.Input)) *Session {
	w := widget.NewInput(value, prompt)
	if more != nil {
		more(w)
	}
	s.widgets = append(s.widgets, w)
	return s
}

// YesNo adds a yes/no prompt bound to value.
func (s *Session) YesNo(value *bool, prompt string, more func(*widget.YesNo)) *Session {
	w := widget.NewYesNo(value, prompt)
	if more != nil {
		more(w)
	}
	s.widgets = append(s.widgets, w)
	return s
}

// Menu adds a single-select menu bound to value.
func (s *Session) Menu(value *string, prompt string, more func(*widget.Menu)) *Session {
	w := widget.NewMenu(value, prompt)
	if more != nil {
		more(w)
	}
	s.widgets = append(s.widgets, w)
	return s
}

// Select adds a multi-select checkbox group.
func (s *Session) Select(prompt string, more func(*widget.Select)) *Session {
	w := widget.NewSelect(prompt)
	if more != nil {
		more(w)
	}
	s.widgets = append(s.widgets, w)
	return s
}

// Note adds a non-interactive message; the user presses Enter to continue.
func (s *Session) Note(text string, more func(*widget.Note)) *Session {
	w := widget.NewNote(text)
	if more != nil {
		more(w)
	}
	s.widgets = append(s.widgets, w)
	return s
}

// AnyKey adds a prompt that continues on any key (except Ctrl+C).
func (s *Session) AnyKey(message string, more func(*widget.AnyKey)) *Session {
	w := widget.NewAnyKey(message)
	if more != nil {
		more(w)
	}
	s.widgets = append(s.widgets, w)
	return s
}

// Run executes the prompt session.
func (s *Session) Run(ctx context.Context) error {
	if s.err != nil {
		return s.err
	}

	var scrOpts []termui.ScreenOption
	if s.colorEnabled != nil {
		scrOpts = append(scrOpts, termui.WithColor(*s.colorEnabled))
	}

	scr, err := termui.OpenScreen(s.in, s.out, scrOpts...)
	if err != nil {
		if errors.Is(err, termui.ErrNotTerminal) {
			return ErrNotTerminal
		}
		return err
	}
	defer func() { _ = scr.Close() }()

	for _, w := range s.widgets {
		if !w.DoWhen() {
			continue
		}
		if err := w.Run(ctx, scr); err != nil {
			return mapRunError(err)
		}
	}
	return nil
}

func mapRunError(err error) error {
	if errors.Is(err, termui.ErrInterrupted) {
		return ErrInterrupted
	}
	if errors.Is(err, context.Canceled) {
		return ErrInterrupted
	}
	return err
}