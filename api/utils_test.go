package main

import (
	"strings"
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	code := generateShortCode()

	if len(code) != shortCodeLength {
		t.Errorf("expected length %d, got %d", shortCodeLength, len(code))
	}

	for _, ch := range code {
		if !strings.ContainsRune(charset, ch) {
			t.Errorf("character %c not in allowed charset", ch)
		}
	}
}

func TestGenerateShortCodeUniqueness(t *testing.T) {
	seen := make(map[string]bool)

	for i := 0; i < 1000; i++ {
		code := generateShortCode()
		if seen[code] {
			t.Errorf("duplicate code generated: %s", code)
		}
		seen[code] = true
	}
}
