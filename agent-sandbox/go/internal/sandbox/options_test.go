package sandbox

import (
	"regexp"
	"testing"
)

func TestNewOptionsGeneratesBranchWhenOmitted(t *testing.T) {
	opts, err := NewOptions("", "codex", "", "", "create a file", false)
	if err != nil {
		t.Fatalf("NewOptions() error = %v", err)
	}

	if !regexp.MustCompile(`^[a-z]{8}-[0-9]{6}$`).MatchString(opts.Branch) {
		t.Fatalf("NewOptions() branch = %q, want eight lowercase letters and six digits", opts.Branch)
	}
}
