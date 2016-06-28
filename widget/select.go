package widget

import (
	"fmt"
	"strings"

	"github.com/burl/termbox-go"
)

// strings for decorating select menu
const (
	StrSelectItemCursor    = string(CharChevronRight)
	StrSelectItemUnChecked = string(CharCircle)
	StrSelectItemChecked   = string(CharCircleFilled)
	// StrSelectItemCursor = "➤"
	//StrSelectItemCursor = "\u261e"
	// StrSelectItemUnChecked = "▢"
	// StrSelectItemChecked   = "\u2713"
)

// SelectableItem - a selectable item, bound to a bool var
type SelectableItem struct {
	Selected *bool
	Name     string
}

// Select an input widget, single line/string of text
type Select struct {
	WBase
	Items []SelectableItem
}

// SelectItemBoolVar - return item
func SelectItemBoolVar(bound *bool, name string) SelectableItem {
	return SelectableItem{bound, name}
}

// WSelect - create new input widget
func WSelect(prompt string, items []SelectableItem) *Select {
	proxy := make([]SelectableItem, len(items))
	proxyBools := make([]bool, len(items))
	for i, item := range items {
		proxy[i].Name = item.Name
		proxy[i].Selected = item.Selected
		if proxy[i].Selected == nil {
			proxy[i].Selected = &proxyBools[i]
		}
	}
	w := &Select{
		WBase: WBase{prompt: prompt, lines: 2 + len(items)},
		Items: proxy,
	}
	w.Init()
	return w
}

// SelectGroup - select group
func SelectGroup(prompt string, items ...SelectableItem) *Select {
	return WSelect(prompt, items)
}

// Item - add an item to the menu
func (w *Select) Item(value *bool, prompt string) *Select {
	w.lines = w.lines + 1
	w.Items = append(w.Items, SelectItemBoolVar(value, prompt))
	return w
}

// DoValidate validation routine
func (w *Select) DoValidate() (msg string) {
	return "validation failed"
}

func (w *Select) selectList() (list []string) {
	for _, choice := range w.Items {
		if *choice.Selected {
			list = append(list, choice.Name)
		}
	}
	return
}

// Render Select widget
func (w *Select) Render(flush func()) {
	w.drawPrompt()
	//row := w.row
	//col := w.rightCol - 1

	nitems := len(w.Items)

	// function to draw items
	draw := func(item int, active bool) {
		cursor := " "
		if active {
			cursor = StrSelectItemCursor
		}
		radio := StrSelectItemUnChecked
		if *w.Items[item].Selected {
			radio = StrSelectItemChecked
		}
		str := fmt.Sprintf("%s %s %s", cursor, radio, w.Items[item].Name)
		color := coldef
		if active {
			color = termbox.ColorCyan
		}
		tbPrint(0, 1+item, color, coldef, str)
	}

	// draw initial menu
	for i := 0; i < nitems; i++ {
		if i == 0 {
			draw(i, true)
		} else {
			draw(i, false)
		}
	}

	flush()

	// set up for loop
	curItem := 0
	done := false
	pollno := 0

	// get input, draw/re draw menu, items
	for {
		pollno = pollno + 1
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case 3:
				dieFromCtlc()
			case 13:
				done = true
				break
			case termbox.KeySpace, 9, termbox.KeyArrowLeft, termbox.KeyArrowRight:
				*w.Items[curItem].Selected = !*w.Items[curItem].Selected
				draw(curItem, true)
			case termbox.KeyArrowUp:
				if curItem > 0 {
					draw(curItem, false)
					curItem = curItem - 1
					draw(curItem, true)
				}
			case termbox.KeyArrowDown:
				if curItem < nitems-1 {
					draw(curItem, false)
					curItem = curItem + 1
					draw(curItem, true)
				}
			}
		}
		flush()
		if done {
			break
		}
	}

	w.drawResult(strings.Join(w.selectList(), ", "))
	flush()
}
