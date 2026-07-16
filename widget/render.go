package widget

import (
	"strings"

	"github.com/burl/inquire/v2/internal/termui"
	"github.com/mattn/go-runewidth"
)

var (
	stylePrompt   = termui.Style{Fg: termui.ColorGreen, Bold: true}
	styleQuestion = termui.Style{Bold: true}
	styleHint     = termui.Style{Faint: true}
	styleAnswer   = termui.Style{Fg: termui.ColorCyan}
	styleError    = termui.Style{Fg: termui.ColorRed}
	styleActive   = termui.Style{Fg: termui.ColorCyan}
	styleCursor   = termui.Style{Rev: true}
)

const (
	charChevronRight = "❯"
	charCircle       = "◯"
	charCircleFilled = "◉"
)

// drawPromptRow paints "? prompt? (hint)" and returns the column for the value.
func drawPromptRow(band *termui.Band, y int, prompt, hint string) int {
	band.WriteString(0, y, "? ", stylePrompt)
	x := 2
	x += writeStyled(band, x, y, prompt+"?", styleQuestion)
	if hint != "" {
		x += writeStyled(band, x, y, " ("+hint+")", styleHint)
	}
	return x + 1
}

// drawSettledRow paints "✔ prompt? answer" on one line.
func drawSettledRow(band *termui.Band, y int, prompt, answer string, masked bool, mask rune) {
	band.WriteString(0, y, "✔ ", stylePrompt)
	x := 2
	x += writeStyled(band, x, y, prompt+"? ", termui.Style{})
	display := answer
	if masked && mask != 0 {
		var b strings.Builder
		for _, r := range answer {
			_ = r
			b.WriteRune(mask)
		}
		display = b.String()
	}
	_ = writeStyled(band, x, y, display, styleAnswer)
}

// drawErrorRow paints a validation error on row y.
func drawErrorRow(band *termui.Band, y int, msg string) {
	band.WriteString(0, y, "✖ ", styleError)
	band.WriteString(2, y, "error: "+msg, termui.Style{})
}

// drawHintRow paints a faint secondary hint on row y.
func drawHintRow(band *termui.Band, y int, msg string) {
	band.WriteString(0, y, msg, styleHint)
}

func writeStyled(band *termui.Band, x, y int, s string, st termui.Style) int {
	band.WriteString(x, y, s, st)
	return runewidth.StringWidth(s)
}

const (
	footerMenu   = "↑/↓ move · ←/→ move · space/enter confirm"
	footerSelect = "↑/↓ move · space toggle · a all · i invert · enter confirm"
)

// drawFooter paints a faint hint row.
func drawFooter(band *termui.Band, y int, msg string) {
	band.WriteString(0, y, msg, styleHint)
}
