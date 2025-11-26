package normalize

import (
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

// OnlyDigits returns a version of s that contains digits only.
// Used to normalize phone numbers and CEP fields in the CSV.
func OnlyDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// MsToDate converts a millisecond Unix timestamp into YYYY-MM-DD.
// Returns an empty string when ms == 0.
func MsToDate(ms int64) string {
	if ms == 0 {
		return ""
	}
	t := time.Unix(0, ms*int64(time.Millisecond))
	return t.Format("2006-01-02")
}

// Slugify turns a human-readable string into a safe filename fragment:
// lowercased, spaces -> dashes, and only keeps letters/digits/-/_.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	// keep only letters, digits, '-' and '_'
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			out = append(out, r)
		}
	}
	// remove consecutive dashes
	clean := make([]rune, 0, len(out))
	var prevDash bool
	for _, r := range out {
		if r == '-' {
			if !prevDash {
				clean = append(clean, r)
				prevDash = true
			}
			continue
		}
		prevDash = false
		clean = append(clean, r)
	}
	if len(clean) == 0 {
		return "atividade"
	}
	return string(clean)
}

// EmptyIfNil safely dereferences optional string pointers for CSV output.
func EmptyIfNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BoolToStr converts a boolean to "true"/"false" for CSV cells.
func BoolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// IntPtrToStr renders optional integer pointers as strings ("" when nil).
func IntPtrToStr(v *int) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(*v)
}

// FixMojibake attempts to fix common mojibake where UTF-8 bytes were
// incorrectly interpreted as ISO-8859-1 (Latin1). If the string does
// not appear garbled, it returns the original. Otherwise it decodes
// the bytes as ISO-8859-1 into UTF-8.
func FixMojibake(s string) string {
	if s == "" {
		return s
	}

	// If string looks fine (valid UTF-8 and doesn't contain common mojibake markers), return early.
	if utf8.ValidString(s) && !strings.ContainsAny(s, "ÃÂÃ") && !strings.Contains(s, "Ã") && !strings.Contains(s, "Â") {
		return s
	}

	// Helper: score how many likely-correct accented chars exist
	score := func(str string) int {
		c := 0
		for _, r := range str {
			if r >= 0x80 && unicode.IsLetter(r) {
				c++
			}
		}
		return c
	}

	best := s
	bestScore := score(s)

	// Try Windows-1252
	if b, err := charmap.Windows1252.NewDecoder().Bytes([]byte(s)); err == nil {
		cand := string(b)
		sc := score(cand)
		if sc > bestScore {
			best = cand
			bestScore = sc
		}
	}

	// Try ISO-8859-1
	if b, err := charmap.ISO8859_1.NewDecoder().Bytes([]byte(s)); err == nil {
		cand := string(b)
		sc := score(cand)
		if sc > bestScore {
			best = cand
			bestScore = sc
		}
	}

	return best
}
