package widget

import (
	"fmt"

	"github.com/burl/termbox-go"
)

// MenuItemT - a menu entry
type MenuItemT struct {
	Value string
	Name  string
}

// Menu a menu widget
type Menu struct {
	WBase
	boundString *string
	items       []MenuItemT
}

// MenuItem - return item
func MenuItem(value, name string) MenuItemT {
	if value == "" {
		value = name
	}
	return MenuItemT{value, name}
}

// NewMenu - create new input widget
func NewMenu(prompt string, items []MenuItemT) *Menu {
	w := &Menu{
		WBase: WBase{prompt: prompt, lines: 2 + len(items)},
		items: items,
	}
	w.Init()
	return w
}

// MenuStringVar - return menu widget that will set a string var with a value
func MenuStringVar(bound *string, prompt string, items ...MenuItemT) *Menu {
	w := NewMenu(prompt, items)
	w.boundString = bound
	return w
}

// Item - add an item to the menu
func (m *Menu) Item(name, prompt string) *Menu {
	m.lines = m.lines + 1
	m.items = append(m.items, MenuItem(name, prompt))
	return m
}

func (m *Menu) drawItem(item int, active bool) {
	cursor := " "
	if active {
		cursor = string(CharChevronRight)
	}
	str := fmt.Sprintf("%s %s", cursor, m.items[item].Name)
	color := coldef
	if active {
		color = termbox.ColorCyan
	}
	tbPrint(0, 1+item, color, coldef, str)
}

func (m *Menu) draw() (curItem int) {
	var dfl string
	var nitems = len(m.items)
	m.drawPrompt()

	if m.boundString != nil {
		dfl = *m.boundString
	}

	for i := 0; i < nitems; i++ {
		if m.items[i].Value == dfl || (dfl == "" && i == 0) {
			m.drawItem(i, true)
			curItem = i
		} else {
			m.drawItem(i, false)
		}
	}
	return
}

// Render Menu
func (m *Menu) Render(flush func()) {

	curItem := m.draw()
	flush()

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
				m.drawItem(curItem, true)
			case termbox.KeyArrowUp:
				if curItem > 0 {
					m.drawItem(curItem, false)
					curItem = curItem - 1
					m.drawItem(curItem, true)
				}
			case termbox.KeyArrowDown:
				if curItem < len(m.items)-1 {
					m.drawItem(curItem, false)
					curItem = curItem + 1
					m.drawItem(curItem, true)
				}
			}
		}
		flush()
	}

	if m.boundString != nil {
		*m.boundString = m.items[curItem].Value
	}

	m.drawResult(m.items[curItem].Name)
	flush()
}
