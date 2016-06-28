package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/burl/inquire"
	"github.com/burl/inquire/widget"
)

func widgetTest() {

	var (
		name, quest, weight       string
		red, green, blue, proceed bool
	)

	name = "Sir Lancelot"
	quest = "grail"
	green = true

	inquire.Query().
		Input(&name, "What is your name", nil).
		Menu(&quest, "What is your quest", func(w *widget.Menu) {
			w.Hint("use arrow keys, pick one")
			w.Item("shrub", "find a shrubbery")
			w.Item("grail", "find the grail")
			w.Item("nuts", "find coconuts")
		}).
		Input(&weight, "What is the weight of an unladen swallow", func(w *widget.Input) {
			w.WhenEqual(&quest, "nuts")
			w.Valid(func(value string) string {
				n, err := strconv.Atoi(value)
				if err != nil || n < 1 {
					return "not good, you need to enter a number"
				}
				return ""
			})
		}).
		Select("what are your favorite colors", func(w *widget.Select) {
			w.Hint("use arrow/space, select multiple")
			w.Item(&red, "red")
			w.Item(&blue, "blue")
			w.Item(&green, "green")
		}).
		YesNo(&proceed, "Continue").
		Exec()

	if !proceed {
		fmt.Println("aborted.")
		os.Exit(1)
	}

	fmt.Printf("name  : %s\n", name)
	fmt.Printf("quest : %s\n", quest)
}

func main() {
	widgetTest()
}
