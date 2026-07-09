package widget

import (
	"context"
	"fmt"

	"github.com/burl/inquire/v2/internal/termui"
)

type menuItem struct {
	tag  string
	name string
}

// Menu is a vertical single-select prompt.
type Menu struct {
	Base
	value  *string
	prompt string
	hint   string
	items  []menuItem
}

// NewMenu constructs a Menu bound to value.
func NewMenu(value *string, prompt string) *Menu {
	return &Menu{value: value, prompt: prompt}
}

// Hint sets hint text shown beside the prompt.
func (w *Menu) Hint(h string) *Menu {
	w.hint = h
	return w
}

// When registers a skip predicate.
func (w *Menu) When(fn func() bool) *Menu {
	w.Base.When(fn)
	return w
}

// Item appends a menu choice. tag is stored in value; name is displayed.
func (w *Menu) Item(tag, name string) {
	if tag == "" {
		tag = name
	}
	w.items = append(w.items, menuItem{tag: tag, name: name})
}

func (w *Menu) defaultIndex() int {
	if w.value == nil || *w.value == "" {
		return 0
	}
	for i, it := range w.items {
		if it.tag == *w.value {
			return i
		}
	}
	return 0
}

func (w *Menu) drawItems(band *termui.Band, cur int) {
	for i, it := range w.items {
		prefix := "  "
		st := termui.Style{}
		if i == cur {
			prefix = charChevronRight + " "
			st = styleActive
		}
		band.WriteString(0, 1+i, prefix+it.name, st)
	}
}

// Run interactively collects a menu selection.
func (w *Menu) Run(ctx context.Context, scr *termui.Screen) error {
	if len(w.items) == 0 {
		return fmt.Errorf("inquire: menu has no items")
	}

	lines := 1 + len(w.items) + 1 // prompt + items + footer
	band, err := scr.OpenBand(ctx, lines)
	if err != nil {
		return err
	}

	cur := w.defaultIndex()
	showHint := w.hint != ""

	draw := func() {
		band.Clear()
		h := w.hint
		if !showHint {
			h = ""
		}
		_ = drawPromptRow(band, 0, w.prompt, h)
		w.drawItems(band, cur)
		drawFooter(band, lines-1, footerMenu)
		_ = band.Flush()
	}
	draw()

	for {
		ev, err := PollKey(ctx, scr, band, draw)
		if err != nil {
			return err
		}
		if ev.Type != termui.EventKey {
			continue
		}
		showHint = false

		switch ev.Key {
		case termui.KeyCtrlC:
			return termui.ErrInterrupted
		case termui.KeyEnter, termui.KeySpace:
			if w.value != nil {
				*w.value = w.items[cur].tag
			}
			band.Clear()
			drawSettledRow(band, 0, w.prompt, w.items[cur].name, false, 0)
			return band.FinalizeStatic(1)
		case termui.KeyUp, termui.KeyLeft:
			if cur > 0 {
				cur--
				draw()
			}
		case termui.KeyDown, termui.KeyRight:
			if cur < len(w.items)-1 {
				cur++
				draw()
			}
		case termui.KeyTab:
			draw()
		}
	}
}
