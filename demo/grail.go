package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/burl/inquire"
	"github.com/burl/inquire/widget"
)

func widgetTest() {

	// vars that will receive the answers to the questions
	var (
		name, quest, weight, passwd string
		red, green, blue, proceed   bool
	)

	name = "Sir Lancelot" // if you want a default value
	quest = "grail"       // for any of the widgets, then
	green = true          // just assign values to the vars

	inquire.Query().
		// if you just have a plain old question, pass nil as the final arg
		Input(&name, "What is your name", nil).
		Menu(&quest, "What is your quest", func(w *widget.Menu) {
			// if you want to do a bit more, pass a callback - each kind
			// of widget will pass a type into the callback where you can...
			w.Hint("use arrow keys, pick one")  // set up custom hint text
			w.Item("shrub", "find a shrubbery") // set up the values and prompts
			w.Item("grail", "find the grail")   // the &quest var will be set
			w.Item("nuts", "find coconuts")     // to one of: shrub, grail or nuts
		}).
		Input(&weight, "What is the weight of an unladen swallow", func(w *widget.Input) {
			// this question will only be shown when the value of quest is "nuts"
			w.WhenEqual(&quest, "nuts") // there is also a generic form
			w.Valid(func(value string) string {
				// and things can be validated ...
				n, err := strconv.Atoi(value)
				if err != nil || n < 1 {
					// invalid input should return a non-empty error message
					return "not good, you need to enter a number"
				}
				// if the data is valid, return an empty string
				return ""
			})
		}).
		Select("what are your favorite colors", func(w *widget.Select) {
			w.Hint("use arrow/space, select multiple")
			w.Item(&red, "red")     // the select or "checkbox" widget will
			w.Item(&blue, "blue")   // toggle the value of the referenced
			w.Item(&green, "green") // boolean variable
		}).
		Input(&passwd, "What is your secret", func(w *widget.Input) {
			w.MaskInput() // or, w.MaskInput('*')
		}).
		YesNo(&proceed, "Continue"). // simple yes/no
		Exec()                       // render all the questions.

	if !proceed {
		fmt.Println("aborted.")
		os.Exit(1)
	}

	fmt.Printf("\nHere are the answers:\n---------------------\n")
	fmt.Printf("name  : %s\n", name)
	fmt.Printf("quest : %s\n", quest)
	fmt.Printf("colors: red:%v, green:%v, blue:%v\n", red, green, blue)
	fmt.Printf("secret: %s (shhh!)\n", passwd)
}

func main() {
	widgetTest()
}
