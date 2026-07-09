// Grail demo: Monty Python quest through every inquire widget.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/burl/inquire/v2"
	"github.com/burl/inquire/v2/widget"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var (
		name, quest, weight, passwd string
		red, green, blue, proceed   bool
	)

	name = "Sir Lancelot"
	quest = "grail"
	green = true

	err := inquire.Query().
		Note("You enter the realm of Monty Python and approach the Bridge of Death.", nil).
		AnyKey("A knight in black armour blocks your path", func(w *widget.AnyKey) {
			w.Hint("press any key to face the Bridge Keeper")
		}).
		Input(&name, "What is your name", nil).
		Menu(&quest, "What is your quest", func(w *widget.Menu) {
			w.Hint("use arrow keys, pick one")
			w.Item("shrub", "find a shrubbery")
			w.Item("grail", "find the grail")
			w.Item("nuts", "find coconuts")
		}).
		Input(&weight, "What is the weight of an unladen swallow", func(w *widget.Input) {
			w.When(widget.WhenEqual(&quest, "nuts"))
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
		Note("If you fail the next question, you shall be cast into the Gorge of Eternal Peril.", nil).
		Input(&passwd, "What is the capital of Assyria", func(w *widget.Input) {
			w.MaskInput()
			w.Hint("shhh")
		}).
		AnyKey("The Bridge Keeper regards your answers with grave suspicion", func(w *widget.AnyKey) {
			w.Hint("press any key for judgment")
		}).
		YesNo(&proceed, "May you cross the bridge", func(w *widget.YesNo) {
			w.Hint("Yes/No")
		}).
		Run(ctx)

	if err != nil {
		if errors.Is(err, inquire.ErrInterrupted) {
			fmt.Fprintln(os.Stderr, "\ninterrupted")
			os.Exit(130)
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if !proceed {
		fmt.Println("Auuuuuuuugh! (cast into the gorge)")
		os.Exit(1)
	}

	fmt.Printf("\nHere are the answers:\n---------------------\n")
	fmt.Printf("name  : %s\n", name)
	fmt.Printf("quest : %s\n", quest)
	if quest == "nuts" {
		fmt.Printf("weight: %s\n", weight)
	}
	fmt.Printf("colors: red:%v, green:%v, blue:%v\n", red, green, blue)
	fmt.Printf("secret: %s (shhh!)\n", passwd)
	fmt.Println("\nRight. Off you go.")
}