package inquire_test

import (
	"github.com/burl/inquire"
	"github.com/burl/inquire/widget"
)

func Example() {
	// Variables for the answers to questions:
	name := ""
	really := false

	// Create a new list of questions to ask:
	questions := inquire.Query()

	// add questions
	questions.Input(&name, "what is your name", nil)
	questions.YesNo(&really, "is that really your name")

	// ask all the questions...
	questions.Exec()

}

func Example_chaining() {

	// All widget methods return the questions object, so
	// you can also just chain calls to the Questions() method:
	name := ""
	really := false
	inquire.Query().
		Input(&name, "what is your name", nil).
		YesNo(&really, "is that really your name").
		Exec()
}

func Example_when() {

	// Widgets may be optionally displayed based on
	// arbitrary conditions from a 'When' callback:
	mayImbibe := false
	beer := ""
	wine := ""

	// In this example, the user will only be asked about their favorite
	// beer or wine if they answered 'Yes' to the first question.

	inquire.Query().
		YesNo(&mayImbibe, "are you over 21 years of age").
		Input(&beer, "what is your favorite beer", func(w *widget.Input) {
			w.When(func() bool {
				return mayImbibe
			})
		}).
		Input(&wine, "what is your favorite wine", func(w *widget.Input) {
			// since testing equality is common, this is a shortcut
			// for the above/generic Wnen() callback
			w.WhenEqual(&mayImbibe, true)
		}).
		Exec()
}

func Example_valid() {

	// Example of input validation

	planet := ""

	inquire.Query().
		Input(&planet, "what planet do you live on", func(w *widget.Input) {
			w.Valid(func(value string) string {
				if value != "earth" {
					// return a non-empty string if data is invalid, it will
					// be used as error text for the user
					return "nope, I don't believe it, you could not breathe on " + value
				}
				// return an empty string if the data is valid
				return ""
			})
		}).
		Exec()
}
