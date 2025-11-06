package normalize

import (
	"fmt"
	"strings"
	"time"
	"unicode"
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

