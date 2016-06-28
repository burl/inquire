package widget

import (
	"fmt"
	"os"

	"github.com/burl/termbox-go"
	"github.com/mattn/go-runewidth"
)

func dieFromCtlc() {
	termbox.Close()
	fmt.Println("")
	os.Exit(1)
}

func anyKey() {
	termbox.PollEvent()
}

func tbPrint(x, y int, fg, bg termbox.Attribute, msg string) int {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x += runewidth.RuneWidth(c)
	}
	return x
}

func tbClearToEOL(x, y int) {
	cols, _ := termbox.Size()
	for x < cols {
		termbox.SetCell(x, y, ' ', coldef, coldef)
		x = x + 1
	}
}
