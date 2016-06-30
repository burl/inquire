// Package inquire provides a set of widgets for creating simple
// line-oriented terminal user interfaces.
//
// User interfaces consist of line-oriented "widgets" of the following types:
//     YesNo  - prompt for yes/no answer
//     Input  - single line text input with optional default
//     Menu   - vertical menu of choices, chose one
//     Select - vertical menu of checkboxes, chose many
//
// Widgets may be easily declared and bound to variables to receive
// input from the user.
//
// Optional validation may be tied to widgets
// and validation error text will be displayed in the (always available
// and otherwise blank) line below each widget.
//
// Widgets may also be conditionally skipped based on conditions that
// msut be met in a 'When' callback.
//
// inquire uses a modified termbox-go library that is made to render
// part of the termbox buffer into a region of lines as a "viewport,"
// without taking over the entire screen.
//
package inquire

import (
	"fmt"

	"github.com/burl/inquire/widget"
	"github.com/burl/termbox-go"
)

// Questions is an opaque type that represents a set of widgets for interacting
// with the user.
type Questions struct {
	widgets    []widget.Renderable
	lastMenu   *widget.Menu
	lastSelect *widget.Select
}

const coldef = termbox.ColorDefault

func doTermInit() {
	err := termbox.ViewPortInit()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)
	termbox.SetOutputMode(termbox.Output256)
}

func doTermClose() {
	termbox.Close()
}

// Query creates a new Questions type
func Query() *Questions {
	inq := &Questions{
		widgets: []widget.Renderable{},
	}
	return inq
}

// YesNo adds a YesNo widget to a set of questions.  value will be
// set to true for yes and false for no, based on user input.
func (inq *Questions) YesNo(value *bool, prompt string) *Questions {
	w := widget.YesNoBoolVar(value, prompt)
	inq.widgets = append(inq.widgets, w)
	return inq
}

// Input adds an Input widget to a set of questions.  The argument 'value'
// should point to a string var, which will be assigned the result of the
// user interaction with the Input widget.  If the initial content of 'value'
// is non-empty, then it will be used as the default answer.
//
// A callback function may be passed as the final argument in order to
// register validation callbacks or set up conditional inqury.  If this
// is not needed, then ``nil`` should be passed.
func (inq *Questions) Input(value *string, prompt string, more func(*widget.Input)) *Questions {
	w := widget.InputStringVar(value, prompt, "")
	if more != nil {
		more(w)
	}
	inq.widgets = append(inq.widgets, w)
	return inq
}

// Menu displays a veritcal menu of choices for the user to select one of.
// The final argument may be nil or a function that will be called back
// with the menu widget for further configuration. The value string pointed
// to by the first argument will receive the value of the "tag" string
// for the menu item chosen
func (inq *Questions) Menu(value *string, prompt string, more func(*widget.Menu)) *Questions {
	w := widget.MenuStringVar(value, prompt)
	if more != nil {
		more(w)
	}
	inq.lastMenu = w
	inq.widgets = append(inq.widgets, w)
	return inq
}

// MenuItem may be used to append a menu item to the most recently added
// 'Menu' widget.
func (inq *Questions) MenuItem(tag, prompt string) *Questions {
	if inq.lastMenu == nil {
		panic("no previous menu defined, can not add menu item for: " + tag + ", " + prompt)
	}
	inq.lastMenu.Item(tag, prompt)
	return inq
}

// Select displays a vertical menu of "checkbox" type entries for the user
// to select zero or more of.  The variables bound to these entries must
// be bool vars and are associated with each entry.
func (inq *Questions) Select(prompt string, more func(*widget.Select)) *Questions {
	w := widget.SelectGroup(prompt)
	if more != nil {
		more(w)
	}
	inq.lastSelect = w
	inq.widgets = append(inq.widgets, w)
	return inq
}

// SelectItem may be used to append a select menu item to the most recently
// declared Select widget.
func (inq *Questions) SelectItem(value *bool, prompt string) *Questions {
	if inq.lastSelect == nil {
		panic("no previous select defined, can not add select item for: " + prompt)
	}
	inq.lastSelect.Item(value, prompt)
	return inq
}

// Exec will execute the event loop, prompting the user for input
// for each question defined
func (inq *Questions) Exec() {
	doTermInit()
	defer doTermClose()
	for _, w := range inq.widgets {
		if !w.DoWhen() {
			continue
		}
		termbox.Clear(coldef, coldef)
		row := termbox.ViewPortSetHeight(w.Lines())
		w.SetRow(row)
		w.Render(func() {
			termbox.ViewPortFlush(0, w.Lines(), row-1)
		})
		fmt.Printf("\033[%d;1H", row+1)
	}
}
