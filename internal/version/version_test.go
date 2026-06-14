package version

import (
	"strings"
	"testing"
)

func TestStringNeverEmpty(t *testing.T) {
	// With no ldflags set in the test binary, must still produce something.
	if s := String(); !strings.HasPrefix(s, "aurscan ") {
		t.Fatalf("String() = %q, want prefix 'aurscan '", s)
	}
	if Short() == "" {
		t.Fatal("Short() returned empty")
	}
}

func TestStampedVersionWins(t *testing.T) {
	old := Version
	defer func() { Version = old }()
	Version = "v9.9.9"
	if got := Short(); got != "v9.9.9" {
		t.Fatalf("Short() = %q, want v9.9.9", got)
	}
}
