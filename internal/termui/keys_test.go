package termui

import (
	"bytes"
	"testing"
)

func collect(t *testing.T, raw string) []Event {
	t.Helper()
	var d decoder
	d.feed([]byte(raw))
	var out []Event
	for {
		ev, ok := d.next()
		if !ok {
			break
		}
		out = append(out, ev)
	}
	return out
}

func TestDecodeArrowsCSI(t *testing.T) {
	cases := map[string]Key{
		"\x1b[A": KeyUp,
		"\x1b[B": KeyDown,
		"\x1b[C": KeyRight,
		"\x1b[D": KeyLeft,
	}
	for in, want := range cases {
		evs := collect(t, in)
		if len(evs) != 1 || evs[0].Key != want {
			t.Fatalf("%q => %+v, want key %v", in, evs, want)
		}
	}
}

func TestDecodeArrowsSS3(t *testing.T) {
	cases := map[string]Key{
		"\x1bOA": KeyUp,
		"\x1bOB": KeyDown,
		"\x1bOC": KeyRight,
		"\x1bOD": KeyLeft,
	}
	for in, want := range cases {
		evs := collect(t, in)
		if len(evs) != 1 || evs[0].Key != want {
			t.Fatalf("%q => %+v, want key %v", in, evs, want)
		}
	}
}

func TestDecodeCtrlAndEnter(t *testing.T) {
	evs := collect(t, "\x03\r\x7f")
	want := []Key{KeyCtrlC, KeyEnter, KeyBackspace}
	if len(evs) != len(want) {
		t.Fatalf("got %d events: %+v", len(evs), evs)
	}
	for i := range want {
		if evs[i].Key != want[i] {
			t.Fatalf("ev[%d]=%v want %v", i, evs[i].Key, want[i])
		}
	}
}

func TestDecodeUTF8(t *testing.T) {
	evs := collect(t, "hi界")
	if len(evs) != 3 {
		t.Fatalf("got %d events: %+v", len(evs), evs)
	}
	if evs[0].Rune != 'h' || evs[1].Rune != 'i' || evs[2].Rune != '界' {
		t.Fatalf("runes: %+v", evs)
	}
}

func TestDecodeSplitEscape(t *testing.T) {
	var d decoder
	d.feed([]byte{0x1b, '['})
	if _, ok := d.next(); ok {
		t.Fatal("expected need more")
	}
	d.feed([]byte{'A'})
	ev, ok := d.next()
	if !ok || ev.Key != KeyUp {
		t.Fatalf("got %+v ok=%v", ev, ok)
	}
}

func TestParseCursorReport(t *testing.T) {
	in := []byte("noise\x1b[12;40Rmore")
	row, col, rest, ok := parseCursorReport(in)
	if !ok || row != 12 || col != 40 {
		t.Fatalf("row=%d col=%d ok=%v", row, col, ok)
	}
	if !bytes.Equal(rest, []byte("noisemore")) {
		t.Fatalf("rest=%q", rest)
	}
}

func TestDeleteKey(t *testing.T) {
	evs := collect(t, "\x1b[3~")
	if len(evs) != 1 || evs[0].Key != KeyDelete {
		t.Fatalf("got %+v", evs)
	}
}

func TestDecodeEmacsCtrl(t *testing.T) {
	cases := map[string]Key{
		"\x10": KeyUp,    // Ctrl+P
		"\x0e": KeyDown,  // Ctrl+N
		"\x06": KeyRight, // Ctrl+F
		"\x02": KeyLeft,  // Ctrl+B
		"\x04": KeyCtrlD,
		"\x17": KeyCtrlW,
	}
	for in, want := range cases {
		evs := collect(t, in)
		if len(evs) != 1 || evs[0].Key != want {
			t.Fatalf("%q => %+v, want key %v", in, evs, want)
		}
	}
}
