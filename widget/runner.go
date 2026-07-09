package widget

import (
	"context"

	"github.com/burl/inquire/v2/internal/termui"
)

// Runner is a prompt widget executed against a terminal screen.
type Runner interface {
	DoWhen() bool
	Run(ctx context.Context, scr *termui.Screen) error
}
