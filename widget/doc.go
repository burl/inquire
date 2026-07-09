// Package widget provides configurable prompt types used with [inquire.Query].
//
// Configure widgets in the more callback passed to inquire.Session methods:
//
//	inquire.Query().Input(&name, "your name", func(w *widget.Input) {
//	    w.Valid(func(s string) string { ... })
//	})
//
// Conditional prompts use [Base.When] with [WhenEqual] predicates.
package widget