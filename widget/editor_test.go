package widget_test

import (
	"testing"

	"github.com/burl/inquire/v2/widget"
)

func TestEditorInsertAndBackspace(t *testing.T) {
	ed := widget.NewEditor()
	ed.Insert('h')
	ed.Insert('i')
	if ed.String() != "hi" {
		t.Fatalf("got %q", ed.String())
	}
	ed.Backspace()
	if ed.String() != "h" {
		t.Fatalf("got %q", ed.String())
	}
}

func TestEditorWideRuneCursorCol(t *testing.T) {
	ed := widget.NewEditor()
	ed.SetString("a界")
	ed.Home()
	ed.Right()
	if got := ed.CursorCol(0); got != 1 {
		t.Fatalf("col after 'a' = %d want 1", got)
	}
	ed.Right()
	if got := ed.CursorCol(0); got < 2 {
		t.Fatalf("col after wide rune = %d want >= 2", got)
	}
}

func TestEditorDeleteForwardAndKill(t *testing.T) {
	ed := widget.NewEditor()
	ed.SetString("hello")
	ed.Home()
	ed.DeleteForward()
	if ed.String() != "ello" {
		t.Fatalf("delete forward: %q", ed.String())
	}
	ed.Right()
	ed.KillToEnd()
	if ed.String() != "e" {
		t.Fatalf("kill to end: %q", ed.String())
	}
}

func TestEditorKillWordBackward(t *testing.T) {
	ed := widget.NewEditor()
	ed.SetString("hello world")
	ed.KillWordBackward()
	if ed.String() != "hello " {
		t.Fatalf("end of line: %q", ed.String())
	}

	ed.SetString("one two")
	ed.Home()
	ed.Right()
	ed.Right()
	ed.Right()
	ed.KillWordBackward()
	if ed.String() != " two" {
		t.Fatalf("mid word: %q", ed.String())
	}

	ed.SetString("a  bc")
	ed.End()
	ed.KillWordBackward()
	if ed.String() != "a  " {
		t.Fatalf("trailing spaces: %q", ed.String())
	}
}
