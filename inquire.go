package inquire

import (
	"fmt"

	"github.com/burl/inquire/widget"
	"github.com/burl/termbox-go"
)

const coldef = termbox.ColorDefault

func doTermInit() {
	termbox.OwnConsole = false
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)
	termbox.SetOutputMode(termbox.Output256)
}

func doTermClose() {
	termbox.Close()
}

// Questions - a set of widgets
type Questions struct {
	widgets    []widget.Renderable
	lastMenu   *widget.Menu
	lastSelect *widget.Select
}

// Query - create and execute widgets
func Query() *Questions {
	inq := &Questions{
		widgets: []widget.Renderable{},
	}
	return inq
}

// YesNo - add a YesNo widget
func (inq *Questions) YesNo(value *bool, prompt string) *Questions {
	w := widget.YesNoBoolVar(value, prompt)
	inq.widgets = append(inq.widgets, w)
	return inq
}

// Input - a simple string-input widget
func (inq *Questions) Input(value *string, prompt string, more func(*widget.Input)) *Questions {
	w := widget.InputStringVar(value, prompt, "")
	if more != nil {
		more(w)
	}
	inq.widgets = append(inq.widgets, w)
	return inq
}

// Menu - a simple menu widget
func (inq *Questions) Menu(value *string, prompt string, more func(*widget.Menu)) *Questions {
	w := widget.MenuStringVar(value, prompt)
	if more != nil {
		more(w)
	}
	inq.lastMenu = w
	inq.widgets = append(inq.widgets, w)
	return inq
}

// MenuItem - Add Menu Item
func (inq *Questions) MenuItem(name, prompt string) *Questions {
	if inq.lastMenu == nil {
		panic("no previous menu defined, can not add menu item for: " + name + ", " + prompt)
	}
	inq.lastMenu.Item(name, prompt)
	return inq
}

// Select - a simple menu widget
func (inq *Questions) Select(prompt string, more func(*widget.Select)) *Questions {
	w := widget.SelectGroup(prompt)
	if more != nil {
		more(w)
	}
	inq.lastSelect = w
	inq.widgets = append(inq.widgets, w)
	return inq
}

// SelectItem - Add Menu Item
func (inq *Questions) SelectItem(value *bool, prompt string) *Questions {
	if inq.lastSelect == nil {
		panic("no previous select defined, can not add select item for: " + prompt)
	}
	inq.lastSelect.Item(value, prompt)
	return inq
}

// Exec the "ui"
func (inq *Questions) Exec() {
	doTermInit()
	defer doTermClose()
	for _, w := range inq.widgets {
		if !w.DoWhen() {
			continue
		}
		termbox.Clear(coldef, coldef)
		/*
			cols, _ := termbox.Size()
			for r := 0; r < w.Lines(); r++ {
				for c := 0; c < cols; c++ {
					termbox.SetCell(c, r, '.', coldef, coldef)
				}
			}
		*/
		row := termbox.SetLinesOut(w.Lines())
		w.SetRow(row)
		w.Render(func() {
			termbox.FlushRect(0, w.Lines(), row-1)
		})
		fmt.Printf("\033[%d;1H", row+1)
	}
}
