// Command termui exercises the inquire band layer: an inline menu that does
// not take over the full screen. Arrow keys move, Enter selects, Ctrl+C aborts.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/burl/inquire/v2/internal/termui"
)

func main() {
	if err := run(); err != nil {
		if err == termui.ErrInterrupted {
			fmt.Fprintln(os.Stderr, "\ninterrupted")
			os.Exit(130)
		}
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("termui demo — inline band (no alt-screen)")
	fmt.Println("use ↑/↓, Enter to choose; Ctrl+C to abort; try resizing")
	fmt.Println()

	scr, err := termui.OpenScreen(os.Stdin, os.Stdout)
	if err != nil {
		return err
	}
	defer func() { _ = scr.Close() }()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	items := []string{
		"find a shrubbery",
		"find the grail",
		"find coconuts",
	}
	cur := 1

	band, err := scr.OpenBand(ctx, 1+len(items))
	if err != nil {
		return err
	}

	draw := func() {
		band.Clear()
		stPrompt := termui.Style{Fg: termui.ColorGreen, Bold: true}
		stQ := termui.Style{Bold: true}
		stHint := termui.Style{Faint: true}
		stActive := termui.Style{Fg: termui.ColorCyan}
		stIdle := termui.Style{}

		band.WriteString(0, 0, "? ", stPrompt)
		band.WriteString(2, 0, "What is your quest?", stQ)
		band.WriteString(22, 0, " (arrows, enter)", stHint)

		for i, it := range items {
			st := stIdle
			prefix := "  "
			if i == cur {
				st = stActive
				prefix = "❯ "
			}
			band.WriteString(0, 1+i, prefix+it, st)
		}
		_ = band.Flush()
	}

	draw()

	for {
		ev, err := scr.Poll(ctx)
		if err != nil {
			_ = band.Close()
			return err
		}
		switch ev.Type {
		case termui.EventResize:
			_ = band.OnResize(ctx, ev.Cols)
			draw()
		case termui.EventKey:
			switch ev.Key {
			case termui.KeyUp:
				if cur > 0 {
					cur--
					draw()
				}
			case termui.KeyDown:
				if cur < len(items)-1 {
					cur++
					draw()
				}
			case termui.KeyEnter:
				band.Clear()
				band.WriteString(0, 0, "✔ ", termui.Style{Fg: termui.ColorGreen})
				band.WriteString(2, 0, "What is your quest? ", termui.Style{})
				band.WriteString(22, 0, items[cur], termui.Style{Fg: termui.ColorCyan})
				if err := band.FinalizeStatic(1); err != nil {
					return err
				}
				fmt.Printf("selected: %s\n", items[cur])
				return nil
			}
		}
	}
}
