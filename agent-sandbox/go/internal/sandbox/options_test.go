package sandbox

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestNewOptionsGeneratesBranchWhenOmitted(t *testing.T) {
	opts, err := NewOptions(Options{
		AgentName: "codex",
		Prompt:    "create a file",
	})
	if err != nil {
		t.Fatalf("NewOptions() error = %v", err)
	}

	if !regexp.MustCompile(`^[a-z]{8}-[0-9]{6}$`).MatchString(opts.Branch) {
		t.Fatalf("NewOptions() branch = %q, want eight lowercase letters and six digits", opts.Branch)
	}
}

func TestNewOptionsResolvesSupportedAgentImagePaths(t *testing.T) {
	image := filepath.Join(t.TempDir(), "mockup.png")
	if err := os.WriteFile(image, []byte("not a real image"), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"codex", "claude"} {
		t.Run(name, func(t *testing.T) {
			opts, err := NewOptions(Options{AgentName: name, Prompt: "inspect it", Images: []string{image}})
			if err != nil {
				t.Fatal(err)
			}
			if got, want := opts.Images, []string{image}; len(got) != len(want) || got[0] != want[0] {
				t.Fatalf("Images = %q, want %q", got, want)
			}
		})
	}
}

func TestNewOptionsRejectsNonRegularSupportedAgentImages(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"codex", "claude"} {
		for _, image := range []string{filepath.Join(dir, "missing.png"), dir} {
			t.Run(name+"/"+filepath.Base(image), func(t *testing.T) {
				_, err := NewOptions(Options{AgentName: name, Prompt: "inspect it", Images: []string{image}})
				var usageErr UsageError
				if !errors.As(err, &usageErr) {
					t.Fatalf("error = %v, want UsageError", err)
				}
			})
		}
	}
}

func TestNewOptionsIgnoresImagesForUnsupportedAgents(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.png")
	for _, name := range []string{"opencode", "pi"} {
		t.Run(name, func(t *testing.T) {
			opts, err := NewOptions(Options{AgentName: name, Prompt: "inspect it", Images: []string{missing}})
			if err != nil {
				t.Fatal(err)
			}
			if got, want := opts.Images, []string{missing}; len(got) != len(want) || got[0] != want[0] {
				t.Fatalf("Images = %q, want ignored input %q", got, want)
			}
		})
	}
}
