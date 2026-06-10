package redactor

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
)

// Mode controls how detected PII is handled.
type Mode int

const (
	Mask Mode = iota
	Hash
	Drop
)

// piePattern pairs a PII type name with its sub-pattern.
// Order matters: more specific patterns first.
var piePatterns = []struct {
	name string
	expr string
}{
	{"SSN", `\b\d{3}-\d{2}-\d{4}\b`},
	{"CREDIT_CARD", `\b\d{4}[ -]?\d{4}[ -]?\d{4}[ -]?\d{4}\b`},
	{"EMAIL", `\b[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}\b`},
	{"API_TOKEN", `\b(?:sk|pk|tok)_[A-Za-z0-9_]{16,}\b`},
	{"PHONE", `\b(?:\+?1[-.\s]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}\b`},
	{"IP", `\b(?:\d{1,3}\.){3}\d{1,3}\b`},
}

// combined is a single regex with one named group per PII type.
// One scan of the text finds all PII types at once.
var combined = buildCombined()
var groupNames []string

func buildCombined() *regexp.Regexp {
	pattern := ""
	for i, p := range piePatterns {
		if i > 0 {
			pattern += "|"
		}
		// named capture group, e.g. (?P<EMAIL>...)
		pattern += "(?P<" + p.name + ">" + p.expr + ")"
		groupNames = append(groupNames, p.name)
	}
	return regexp.MustCompile(pattern)
}

type Redactor struct {
	mode Mode
}

func New(mode Mode) *Redactor {
	return &Redactor{mode: mode}
}

// Redact scans the text once, replacing every PII match per the mode.
// Redact scans the text once and rebuilds it with PII replaced.
func (r *Redactor) Redact(text string) (string, int) {
	names := combined.SubexpNames()
	// indices: for each match, [start,end, g1s,g1e, g2s,g2e, ...]
	matches := combined.FindAllStringSubmatchIndex(text, -1)
	if matches == nil {
		return text, 0
	}

	var b []byte
	last := 0
	count := 0

	for _, m := range matches {
		start, end := m[0], m[1]
		// find which named group matched (first non-empty submatch)
		typ := "PII"
		for g := 1; g < len(names); g++ {
			gs, ge := m[2*g], m[2*g+1]
			if gs >= 0 && ge > gs && names[g] != "" {
				typ = names[g]
				break
			}
		}

		b = append(b, text[last:start]...) // text before the match
		switch r.mode {
		case Hash:
			b = append(b, ("[" + typ + ":" + shortHash(text[start:end]) + "]")...)
		case Drop:
			// append nothing
		default:
			b = append(b, ("[REDACTED:" + typ + "]")...)
		}
		last = end
		count++
	}
	b = append(b, text[last:]...) // remaining tail
	return string(b), count
}

func shortHash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])[:8]
}
