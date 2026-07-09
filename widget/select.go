package widget

import (
	"context"
	"fmt"
	"strings"

	"github.com/burl/inquire/v2/internal/termui"
)

type selectItem struct {
	selected *bool
	name     string
}

// Select is a vertical multi-select checkbox prompt.
type Select struct {
	Base
	prompt     string
	hint       string
	items      []selectItem
	proxyBools []bool
}

// NewSelect constructs a Select group.
func NewSelect(prompt string) *Select {
	return &Select{prompt: prompt}
}

// Hint sets hint text shown beside the prompt.
func (w *Select) Hint(h string) *Select {
	w.hint = h
	return w
}

// When registers a skip predicate.
func (w *Select) When(fn func() bool) *Select {
	w.Base.When(fn)
	return w
}

// Item appends a checkbox entry bound to value.
func (w *Select) Item(value *bool, name string) {
	if value == nil {
		w.proxyBools = append(w.proxyBools, false)
		value = &w.proxyBools[len(w.proxyBools)-1]
	}
	w.items = append(w.items, selectItem{selected: value, name: name})
}

func (w *Select) selectedNames() string {
	var names []string
	for _, it := range w.items {
		if it.selected != nil && *it.selected {
			names = append(names, it.name)
		}
	}
	if len(names) == 0 {
		return "(none)"
	}
	return strings.Join(names, ", ")
}

func (w *Select) drawItems(band *termui.Band, cur int) {
	for i, it := range w.items {
		cursor := " "
		if i == cur {
			cursor = charChevronRight
		}
		mark := charCircle
		if it.selected != nil && *it.selected {
			mark = charCircleFilled
		}
		st := termui.Style{}
		if i == cur {
			st = styleActive
		}
		band.WriteString(0, 1+i, cursor+" "+mark+" "+it.name, st)
	}
}

// Run interactively collects checkbox selections.
func (w *Select) Run(ctx context.Context, scr *termui.Screen) error {
	if len(w.items) == 0 {
		return fmt.Errorf("inquire: select has no items")
	}

	lines := 1 + len(w.items) + 1 // prompt + items + footer
	band, err := scr.OpenBand(ctx, lines)
	if err != nil {
		return err
	}

	cur := 0
	showHint := w.hint != ""

	draw := func() {
		band.Clear()
		h := w.hint
		if !showHint {
			h = ""
		}
		_ = drawPromptRow(band, 0, w.prompt, h)
		w.drawItems(band, cur)
		drawFooter(band, lines-1, footerSelect)
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
		case termui.KeyEnter:
			band.Clear()
			drawSettledRow(band, 0, w.prompt, w.selectedNames(), false, 0)
			return band.FinalizeStatic(1)
		case termui.KeyUp:
			if cur > 0 {
				cur--
				draw()
			}
		case termui.KeyDown:
			if cur < len(w.items)-1 {
				cur++
				draw()
			}
		case termui.KeySpace, termui.KeyTab:
			if w.items[cur].selected != nil {
				*w.items[cur].selected = !*w.items[cur].selected
			}
			draw()
		case termui.KeyRune:
			switch ev.Rune {
			case 'a', 'A':
				for i := range w.items {
					if w.items[i].selected != nil {
						*w.items[i].selected = true
					}
				}
				draw()
			case 'i', 'I':
				for i := range w.items {
					if w.items[i].selected != nil {
						*w.items[i].selected = !*w.items[i].selected
					}
				}
				draw()
			}
		}
	}
}
