package widget_test

import (
	"testing"

	"github.com/burl/inquire/v2/widget"
)

func TestWhenEqual(t *testing.T) {
	flag := false
	yn := widget.NewYesNo(&flag, "prompt")
	yn.When(widget.WhenEqual(&flag, true))
	if yn.DoWhen() {
		t.Fatal("expected skip when flag is false")
	}
	flag = true
	if !yn.DoWhen() {
		t.Fatal("expected run when flag is true")
	}
}
