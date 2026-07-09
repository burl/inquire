package inquire

import "os"

// Option configures a query session.
type Option func(*Session)

// WithInput sets the terminal input stream (must be a TTY for Run).
func WithInput(in *os.File) Option {
	return func(s *Session) {
		s.in = in
	}
}

// WithOutput sets the terminal output stream (must be a TTY for Run).
func WithOutput(out *os.File) Option {
	return func(s *Session) {
		s.out = out
	}
}

// WithColor forces ANSI color on (true) or off (false).
// When unset, color follows NO_COLOR, TERM, and COLORTERM.
func WithColor(enabled bool) Option {
	return func(s *Session) {
		s.colorEnabled = &enabled
	}
}