package widget

// WhenEqual returns a predicate that is true when *ptr equals want.
// Use with [Base.When] or widget When methods:
//
//	w.When(WhenEqual(&quest, "nuts"))
func WhenEqual[T comparable](ptr *T, want T) func() bool {
	return func() bool { return *ptr == want }
}