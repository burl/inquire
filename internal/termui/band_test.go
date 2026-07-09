package termui

import "testing"

func TestLineLastContent(t *testing.T) {
	b := &Band{cols: 10}
	row := make([]cell, 10)
	row[0] = cell{Ch: 'h'}
	row[1] = cell{Ch: 'i'}
	for i := 2; i < 10; i++ {
		row[i] = cell{Ch: ' '}
	}
	if got := b.lineLastContent(row); got != 1 {
		t.Fatalf("got %d want 1", got)
	}

	blank := make([]cell, 5)
	for i := range blank {
		blank[i] = cell{Ch: ' '}
	}
	if got := b.lineLastContent(blank); got != -1 {
		t.Fatalf("blank got %d want -1", got)
	}

	row[5] = cell{Ch: ' ', St: Style{Rev: true}}
	if got := b.lineLastContent(row); got != 5 {
		t.Fatalf("styled space got %d want 5", got)
	}
}
