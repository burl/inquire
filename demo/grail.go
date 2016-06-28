package main

import (
	"fmt"
	"os"

	"github.com/burl/inquire"
	"github.com/burl/inquire/widget"
)

func widgetTest() {

	var (
		name, quest               string
		red, green, blue, proceed bool
	)

	name = os.Getenv("LOGNAME")
	quest = "grail"

	inquire.Query().
		Input(&name, "What is your name", nil).
		Menu(&quest, "What is your quest", func(m *widget.Menu) {
			m.Item("shrub", "find a shrubbery")
			m.Item("grail", "find the grail")
			m.Item("bridge", "find the bridge")
		}).
		Select("what are your favorite colors", func(s *widget.Select) {
			s.Item(&red, "red")
			s.Item(&blue, "blue")
			s.Item(&green, "green")
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
