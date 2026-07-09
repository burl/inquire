package termui

import (
	"testing"
)

func TestFakeCapturesOutput(t *testing.T) {
	ctx := t.Context()
	script := Script(CursorReport(4, 1), KeyEnter)

	fake, err := NewFake(40, 10, script)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = fake.Close() }()

	band, err := fake.Screen.OpenBand(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	band.WriteString(0, 0, "hello", Style{})
	if err := band.FinalizeStatic(1); err != nil {
		t.Fatal(err)
	}
	if err := fake.Close(); err != nil {
		t.Fatal(err)
	}
	if fake.Output() == "" {
		t.Fatal("expected captured output")
	}
}
