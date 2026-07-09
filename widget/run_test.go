package widget_test

import (
	"errors"
	"testing"

	"github.com/burl/inquire/v2/internal/termui"
	"github.com/burl/inquire/v2/widget"
)

func openFake(t *testing.T, parts ...any) *termui.Fake {
	t.Helper()
	fake, err := termui.NewFake(80, 24, termui.Script(parts...))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = fake.Close() })
	return fake
}

func TestInputRun(t *testing.T) {
	ctx := t.Context()
	var got string
	in := widget.NewInput(&got, "name")

	fake := openFake(t, termui.CursorReport(5, 1), 'h', 'i', termui.KeyEnter)
	if err := in.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
	if got != "hi" {
		t.Fatalf("got %q want hi", got)
	}
}

func TestInputValidationRetry(t *testing.T) {
	ctx := t.Context()
	var got string
	in := widget.NewInput(&got, "n")
	in.Valid(func(s string) string {
		if s == "ok" {
			return ""
		}
		return "nope"
	})

	fake := openFake(t,
		termui.CursorReport(5, 1),
		'b', 'a', 'd', termui.KeyEnter,
		termui.KeyBackspace, termui.KeyBackspace, termui.KeyBackspace,
		'o', 'k', termui.KeyEnter,
	)
	if err := in.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
	if got != "ok" {
		t.Fatalf("got %q want ok", got)
	}
}

func TestYesNoToggle(t *testing.T) {
	ctx := t.Context()
	got := true
	yn := widget.NewYesNo(&got, "go")

	fake := openFake(t, termui.CursorReport(5, 1), termui.KeySpace, termui.KeyEnter)
	if err := yn.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
	if got {
		t.Fatal("expected false after toggle")
	}
}

func TestYesNoInterrupted(t *testing.T) {
	ctx := t.Context()
	yn := widget.NewYesNo(new(bool), "go")

	fake := openFake(t, termui.CursorReport(5, 1), termui.KeyCtrlC)
	err := yn.Run(ctx, fake.Screen)
	if !errors.Is(err, termui.ErrInterrupted) {
		t.Fatalf("got %v want ErrInterrupted", err)
	}
}

func TestMenuEmptyItems(t *testing.T) {
	ctx := t.Context()
	m := widget.NewMenu(new(string), "pick")

	fake := openFake(t, termui.CursorReport(5, 1))
	err := m.Run(ctx, fake.Screen)
	if err == nil {
		t.Fatal("expected empty menu error")
	}
}

func TestMenuNavigate(t *testing.T) {
	ctx := t.Context()
	var got string
	m := widget.NewMenu(&got, "pick")
	m.Item("a", "alpha")
	m.Item("b", "beta")

	fake := openFake(t, termui.CursorReport(5, 1), termui.KeyDown, termui.KeyEnter)
	if err := m.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
	if got != "b" {
		t.Fatalf("got %q want b", got)
	}
}

func TestMenuNavigateCtrlN(t *testing.T) {
	ctx := t.Context()
	var got string
	m := widget.NewMenu(&got, "pick")
	m.Item("a", "alpha")
	m.Item("b", "beta")

	fake := openFake(t, termui.CursorReport(5, 1), byte(0x0e), termui.KeyEnter) // Ctrl+N
	if err := m.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
	if got != "b" {
		t.Fatalf("got %q want b", got)
	}
}

func TestInputEmacsEditing(t *testing.T) {
	ctx := t.Context()
	var got string
	in := widget.NewInput(&got, "name")

	fake := openFake(t,
		termui.CursorReport(5, 1),
		'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd',
		termui.KeyCtrlW,
		byte(0x04), // Ctrl+D
		termui.KeyEnter,
	)
	if err := in.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
	if got != "hello " {
		t.Fatalf("got %q want %q", got, "hello ")
	}
}

func TestMenuSelectWithSpace(t *testing.T) {
	ctx := t.Context()
	var got string
	m := widget.NewMenu(&got, "pick")
	m.Item("a", "alpha")
	m.Item("b", "beta")

	fake := openFake(t, termui.CursorReport(5, 1), termui.KeyDown, termui.KeySpace)
	if err := m.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
	if got != "b" {
		t.Fatalf("got %q want b", got)
	}
}

func TestSelectNoneDisplay(t *testing.T) {
	ctx := t.Context()
	red := false
	sel := widget.NewSelect("colors")
	sel.Item(&red, "red")

	fake := openFake(t, termui.CursorReport(5, 1), termui.KeyEnter)
	if err := sel.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
}

func TestSelectInvertAll(t *testing.T) {
	ctx := t.Context()
	red, blue := true, false
	sel := widget.NewSelect("colors")
	sel.Item(&red, "red")
	sel.Item(&blue, "blue")

	fake := openFake(t, termui.CursorReport(5, 1), 'i', termui.KeyEnter)
	if err := sel.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
	if red || !blue {
		t.Fatalf("invert: red=%v blue=%v want false,true", red, blue)
	}

	red, blue = false, false
	sel = widget.NewSelect("colors")
	sel.Item(&red, "red")
	sel.Item(&blue, "blue")
	fake = openFake(t, termui.CursorReport(5, 1), 'a', termui.KeyEnter)
	if err := sel.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
	if !red || !blue {
		t.Fatalf("all: red=%v blue=%v want true,true", red, blue)
	}
}

func TestNoteRun(t *testing.T) {
	ctx := t.Context()
	n := widget.NewNote("hello")

	fake := openFake(t, termui.CursorReport(5, 1), termui.KeyEnter)
	if err := n.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
}

func TestAnyKeyRun(t *testing.T) {
	ctx := t.Context()
	ak := widget.NewAnyKey("continue")

	fake := openFake(t, termui.CursorReport(5, 1), 'x')
	if err := ak.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
}

func TestSelectToggle(t *testing.T) {
	ctx := t.Context()
	red, blue := false, false
	sel := widget.NewSelect("colors")
	sel.Item(&red, "red")
	sel.Item(&blue, "blue")

	fake := openFake(t, termui.CursorReport(5, 1), termui.KeySpace, termui.KeyEnter)
	if err := sel.Run(ctx, fake.Screen); err != nil {
		t.Fatal(err)
	}
	if !red || blue {
		t.Fatalf("red=%v blue=%v want true,false", red, blue)
	}
}
