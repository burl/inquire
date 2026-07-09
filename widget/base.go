package widget

// Base holds shared widget configuration.
type Base struct {
	when func() bool
}

// When registers a predicate; the widget is skipped when DoWhen returns false.
func (b *Base) When(fn func() bool) {
	b.when = fn
}

// DoWhen reports whether this widget should run.
func (b *Base) DoWhen() bool {
	if b.when == nil {
		return true
	}
	return b.when()
}
