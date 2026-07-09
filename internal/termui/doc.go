// Package termui implements the inline terminal band layer for inquire.
//
// It provides raw TTY mode, ANSI cursor and line control, key decoding, and
// fixed-height bands anchored at the current cursor — without entering the
// alternate screen. Widgets paint into bands; answered prompts settle to
// static scrollback via FinalizeStatic.
package termui
