package widget

import (
	"sort"
	"strings"
)

const maxCompleteHints = 6

// CompleteFrom returns a prefix matcher over fixed candidates (case-insensitive).
func CompleteFrom(candidates []string) func(prefix string) []string {
	cp := append([]string(nil), candidates...)
	sort.Strings(cp)
	return func(prefix string) []string {
		var out []string
		lower := strings.ToLower(prefix)
		for _, c := range cp {
			if strings.HasPrefix(strings.ToLower(c), lower) {
				out = append(out, c)
			}
		}
		return out
	}
}

type tabCompleteState struct {
	pending bool
	matches []string
	index   int
}

func (s *tabCompleteState) reset() {
	s.pending = false
	s.matches = nil
	s.index = 0
}

// applyTabCompletion updates ed from fn(prefix) and returns a faint hint for row 2.
func applyTabCompletion(ed *Editor, fn func(string) []string, st *tabCompleteState) string {
	if st.pending && len(st.matches) > 0 {
		st.index = (st.index + 1) % len(st.matches)
		ed.SetString(st.matches[st.index])
		return formatCompleteHint(st.matches, st.index)
	}

	prefix := string(ed.runes[:ed.pos])
	matches := fn(prefix)
	sort.Strings(matches)
	st.matches = matches
	st.index = 0
	st.pending = false

	switch len(matches) {
	case 0:
		return "no matches"
	case 1:
		ed.SetString(matches[0])
		return ""
	default:
		lcp := longestCommonPrefix(matches)
		if len(lcp) > len(prefix) {
			ed.SetString(lcp)
			st.pending = true
			return formatCompleteHint(matches, -1)
		}
		st.pending = true
		ed.SetString(matches[0])
		return formatCompleteHint(matches, 0)
	}
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, s := range strs[1:] {
		for len(prefix) > 0 && !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
}

func formatCompleteHint(matches []string, active int) string {
	if len(matches) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("matches (")
	b.WriteString(itoa(len(matches)))
	b.WriteString("): ")
	limit := len(matches)
	if limit > maxCompleteHints {
		limit = maxCompleteHints
	}
	for i := 0; i < limit; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		if i == active {
			b.WriteByte('[')
		}
		b.WriteString(matches[i])
		if i == active {
			b.WriteByte(']')
		}
	}
	if len(matches) > maxCompleteHints {
		b.WriteString(", … (tab cycles)")
	} else if len(matches) > 1 {
		b.WriteString(" (tab cycles)")
	}
	return b.String()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}