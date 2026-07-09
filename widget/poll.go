package widget

import (
	"context"

	"github.com/burl/inquire/v2/internal/termui"
)

// PollKey waits for a key event, repainting the band after terminal resize.
func PollKey(ctx context.Context, scr *termui.Screen, band *termui.Band, repaint func()) (termui.Event, error) {
	for {
		ev, err := scr.Poll(ctx)
		if err != nil {
			return ev, err
		}
		switch ev.Type {
		case termui.EventResize:
			_ = band.OnResize(ctx, ev.Cols)
			repaint()
		case termui.EventKey:
			return ev, nil
		}
	}
}
