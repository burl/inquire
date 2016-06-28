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

// SelectItemBoolVar - return item
func SelectItemBoolVar(bound *bool, name string) SelectableItem {
	return SelectableItem{bound, name}
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

// function to draw items
func (w *Select) drawItem(item int, active bool) {
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

func (w *Select) draw() {
	w.drawPrompt()
	for i := 0; i < len(w.Items); i++ {
		if i == 0 {
			w.drawItem(i, true)
		} else {
			w.drawItem(i, false)
		}
	}
}

// Render Select widget
func (w *Select) Render(flush func()) {
	curItem := 0

	w.draw()
	flush()

	// get input, draw/re draw menu, items
EventLoop:
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case 3:
				dieFromCtlc()
			case 13:
				break EventLoop
			case termbox.KeySpace, 9, termbox.KeyArrowLeft, termbox.KeyArrowRight:
				*w.Items[curItem].Selected = !*w.Items[curItem].Selected
				w.drawItem(curItem, true)
			case termbox.KeyArrowUp:
				if curItem > 0 {
					w.drawItem(curItem, false)
					curItem = curItem - 1
					w.drawItem(curItem, true)
				}
			case termbox.KeyArrowDown:
				if curItem < len(w.Items)-1 {
					w.drawItem(curItem, false)
					curItem = curItem + 1
					w.drawItem(curItem, true)
				}
			}
		}
		flush()
	}

	w.drawResult(strings.Join(w.selectList(), ", "))
	flush()
}
