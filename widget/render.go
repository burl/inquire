package widget

import (
	"fmt"

	"github.com/burl/termbox-go"
)

const coldef = termbox.ColorDefault

func doTermInit() {
	err := termbox.ViewPortInit()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)
}

func doTermClose() {
	termbox.Close()
}

// Render the "ui"
func Render(widgets ...Renderable) {
	doTermInit()
	defer doTermClose()
	for _, w := range widgets {
		termbox.Clear(coldef, coldef)
		/*
			cols, _ := termbox.Size()
			for r := 0; r < w.Lines(); r++ {
				for c := 0; c < cols; c++ {
					termbox.SetCell(c, r, '.', coldef, coldef)
				}
			}
		*/
		row := termbox.ViewPortSetHeight(w.Lines())
		w.SetRow(row)
		w.Render(func() {
			termbox.ViewPortFlush(0, w.Lines(), row-1)
		})
		fmt.Printf("\033[%d;1H", row+1)
	}
}
