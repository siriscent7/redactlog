package redactor

import (
	"strings"
	"testing"
)

func TestRedactEmail(t *testing.T) {
	r := New(Mask)
	out, n := r.Redact("contact me at jane.doe@example.com please")
	if strings.Contains(out, "jane.doe@example.com") {
		t.Fatal("email was not redacted")
	}
	if n != 1 {
		t.Fatalf("expected 1 redaction, got %d", n)
	}
	if !strings.Contains(out, "[REDACTED:EMAIL]") {
		t.Fatalf("expected mask label, got %q", out)
	}
}

func TestRedactSSN(t *testing.T) {
	r := New(Mask)
	out, n := r.Redact("SSN is 123-45-6789")
	if strings.Contains(out, "123-45-6789") || n != 1 {
		t.Fatalf("SSN not redacted: %q (n=%d)", out, n)
	}
}

func TestRedactMultiple(t *testing.T) {
	r := New(Mask)
	in := "email a@b.com, ssn 111-22-3333, token sk_abcdefghij1234567890"
	out, n := r.Redact(in)
	if n != 3 {
		t.Fatalf("expected 3 redactions, got %d: %q", n, out)
	}
}

func TestHashModeStable(t *testing.T) {
	r := New(Hash)
	out1, _ := r.Redact("email a@b.com")
	out2, _ := r.Redact("email a@b.com")
	if out1 != out2 {
		t.Fatal("hash mode should be deterministic for the same input")
	}
	if strings.Contains(out1, "a@b.com") {
		t.Fatal("hash mode leaked the original value")
	}
}

func TestDropMode(t *testing.T) {
	r := New(Drop)
	out, _ := r.Redact("user a@b.com logged in")
	if strings.Contains(out, "a@b.com") {
		t.Fatal("drop mode left PII")
	}
}

func TestNoFalsePositiveOnPlainText(t *testing.T) {
	r := New(Mask)
	out, n := r.Redact("the quick brown fox jumps over the lazy dog")
	if n != 0 {
		t.Fatalf("expected 0 redactions on clean text, got %d: %q", n, out)
	}
}

func BenchmarkRedact(b *testing.B) {
	r := New(Mask)
	line := "user jane.doe@example.com from 10.0.0.5 ssn 123-45-6789 token sk_abcdefghij1234567890 ok ok ok ok"
	b.SetBytes(int64(len(line)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Redact(line)
	}
}

func TestRecallOnSyntheticPII(t *testing.T) {
	r := New(Mask)
	// each line has exactly one known PII item -> 5 total
	lines := []string{
		"login from user alice@corp.com today",
		"ssn on file 222-33-4444 verified",
		"card 4111111111111111 charged",
		"call me at 415-555-0199 tomorrow",
		"api key sk_test_abcdefghij1234567890 rotated",
	}
	expected := len(lines)
	found := 0
	for _, l := range lines {
		_, n := r.Redact(l)
		found += n
	}
	recall := float64(found) / float64(expected)
	if recall < 1.0 {
		t.Logf("recall = %.0f%% (%d/%d)", recall*100, found, expected)
	}
	if found < expected {
		t.Fatalf("missed PII: found %d, expected %d", found, expected)
	}
}

func BenchmarkRedactRealistic(b *testing.B) {
	r := New(Mask)
	// realistic log line: mostly non-PII, occasional PII
	line := `2026-06-07T08:15:42Z INFO request_id=abc123 method=GET path=/api/users status=200 latency_ms=42 user=jane.doe@example.com result=success cache=hit region=us-west`
	b.SetBytes(int64(len(line)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Redact(line)
	}
}
