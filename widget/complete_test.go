package widget

import (
	"strings"
	"testing"
)

func TestLongestCommonPrefix(t *testing.T) {
	got := longestCommonPrefix([]string{"us-east-1", "us-east-2", "us-west-1"})
	if got != "us-" {
		t.Fatalf("got %q want us-", got)
	}
}

func TestCompleteFrom(t *testing.T) {
	fn := CompleteFrom([]string{"African", "European", "Sir Lancelot"})
	m := fn("eu")
	if len(m) != 1 || m[0] != "European" {
		t.Fatalf("eu => %v", m)
	}
	m = fn("sir")
	if len(m) != 1 || m[0] != "Sir Lancelot" {
		t.Fatalf("sir => %v", m)
	}
}

func TestApplyTabCompletion(t *testing.T) {
	ed := NewEditor()
	fn := CompleteFrom([]string{"African", "European"})
	var st tabCompleteState

	ed.Insert('A')
	hint := applyTabCompletion(ed, fn, &st)
	if ed.String() != "African" {
		t.Fatalf("single match: %q", ed.String())
	}
	if hint != "" {
		t.Fatalf("hint %q want empty", hint)
	}

	ed.SetString("e")
	st.reset()
	hint = applyTabCompletion(ed, fn, &st)
	if ed.String() != "European" {
		t.Fatalf("e => %q", ed.String())
	}
	if hint != "" {
		t.Fatalf("hint %q want empty", hint)
	}

	ed.SetString("")
	st.reset()
	hint = applyTabCompletion(ed, fn, &st)
	if hint == "" || ed.String() == "" {
		t.Fatalf("empty prefix: hint=%q val=%q", hint, ed.String())
	}

	ed.SetString("zzz")
	st.reset()
	hint = applyTabCompletion(ed, fn, &st)
	if hint != "no matches" {
		t.Fatalf("got hint %q", hint)
	}
}

func TestFormatCompleteHintTruncates(t *testing.T) {
	matches := make([]string, 50)
	for i := range matches {
		matches[i] = "Sir Knight " + itoa(i+1)
	}
	hint := formatCompleteHint(matches, 0)
	if !strings.Contains(hint, "matches (50):") {
		t.Fatalf("got %q", hint)
	}
	if !strings.Contains(hint, "… (tab cycles)") {
		t.Fatalf("want truncation marker: %q", hint)
	}
}

func TestApplyTabCompletionCycle(t *testing.T) {
	ed := NewEditor()
	fn := CompleteFrom([]string{"Sir Galahad", "Sir Lancelot"})
	var st tabCompleteState

	ed.SetString("Sir ")
	hint := applyTabCompletion(ed, fn, &st)
	if ed.String() != "Sir Galahad" {
		t.Fatalf("first: %q", ed.String())
	}
	if !strings.Contains(hint, "[Sir Galahad]") {
		t.Fatalf("hint %q", hint)
	}

	hint = applyTabCompletion(ed, fn, &st)
	if ed.String() != "Sir Lancelot" {
		t.Fatalf("cycle: %q", ed.String())
	}
	if !strings.Contains(hint, "[Sir Lancelot]") {
		t.Fatalf("hint %q", hint)
	}
}