package certdays

import (
	"testing"
	"time"
)

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"example.com":              "example.com:443",
		"https://example.com/path": "example.com:443",
		"example.com:8443":         "example.com:8443",
	}
	for in, want := range cases {
		if got, _ := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDaysUntil(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if d := DaysUntil(now, now.Add(36*time.Hour)); d != 1 {
		t.Fatalf("got %d", d)
	}
	if d := DaysUntil(now, now.Add(-48*time.Hour)); d != -2 {
		t.Fatalf("got %d", d)
	}
}
