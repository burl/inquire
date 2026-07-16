package termui

import (
	"strings"
	"testing"
)

func TestFinalizeStaticParksCursorOnNextLine(t *testing.T) {
	ctx := t.Context()
	fake, err := NewFake(40, 24, Script(CursorReport(3, 1), CursorReport(4, 1)))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = fake.Close() }()

	band, err := fake.Screen.OpenBand(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	band.WriteString(0, 0, "✔ first? a", Style{})
	if err := band.FinalizeStatic(1); err != nil {
		t.Fatal(err)
	}

	next, err := fake.Screen.OpenBand(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if got := next.OriginRow(); got != 4 {
		t.Fatalf("next origin row = %d want 4 (no blank line between prompts)", got)
	}
}

func TestFinalizeStaticCommitsLineToScrollback(t *testing.T) {
	ctx := t.Context()
	fake, err := NewFake(40, 24, Script(CursorReport(3, 1)))
	if err != nil {
		t.Fatal(err)
	}

	band, err := fake.Screen.OpenBand(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	band.WriteString(0, 0, "› note text", Style{})
	if err := band.FinalizeStatic(1); err != nil {
		t.Fatal(err)
	}
	if err := fake.Close(); err != nil {
		t.Fatal(err)
	}

	out := fake.Output()
	if !strings.Contains(out, "\n") {
		t.Fatal("expected newline after FinalizeStatic to commit settled line")
	}
}