package inquire_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/burl/inquire/v2"
)

func ExampleQuery_notTerminal() {
	err := inquire.Query().Run(context.TODO())
	fmt.Println(errors.Is(err, inquire.ErrNotTerminal))
	// Output: true
}

func ExampleQuery_chain() {
	var name string
	var ok bool

	// Widgets chain on Query; Run needs a real TTY in production.
	_ = inquire.Query().
		Input(&name, "what is your name", nil).
		YesNo(&ok, "continue", nil)

	// Output:
}
