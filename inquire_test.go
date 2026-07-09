package inquire_test

import (
	"errors"
	"testing"

	"github.com/burl/inquire/v2"
)

func TestRunNotTerminal(t *testing.T) {
	ctx := t.Context()
	err := inquire.Query().Run(ctx)
	if !errors.Is(err, inquire.ErrNotTerminal) {
		t.Fatalf("got %v want ErrNotTerminal", err)
	}
}

func TestErrSentinelsDistinct(t *testing.T) {
	if inquire.ErrNotTerminal == nil || inquire.ErrInterrupted == nil {
		t.Fatal("sentinels must be non-nil")
	}
	if errors.Is(inquire.ErrNotTerminal, inquire.ErrInterrupted) {
		t.Fatal("sentinels must differ")
	}
}

func TestErrNotTerminalIsStable(t *testing.T) {
	err := inquire.ErrNotTerminal
	if !errors.Is(err, inquire.ErrNotTerminal) {
		t.Fatal("sentinel must match itself via errors.Is")
	}
}

func TestErrInterruptedIsStable(t *testing.T) {
	err := inquire.ErrInterrupted
	if !errors.Is(err, inquire.ErrInterrupted) {
		t.Fatal("sentinel must match itself via errors.Is")
	}
}
