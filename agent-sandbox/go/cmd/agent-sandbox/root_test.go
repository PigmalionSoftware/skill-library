package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"agent-sandbox/internal/sandbox"
)

func TestRunArgsAcceptsFilePromptWithoutPositionalPrompt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "prompt.md")
	want := "# Fix login.\n"
	if err := os.WriteFile(path, []byte(want), 0o644); err != nil {
		t.Fatalf("write prompt file: %v", err)
	}

	cmd, _ := newRootCmd()
	if err := cmd.Flags().Set("file-prompt", path); err != nil {
		t.Fatalf("set file-prompt: %v", err)
	}

	if err := runArgs(cmd, nil); err != nil {
		t.Errorf("runArgs() error = %v, want nil", err)
	}
}

func TestRunArgsRejectsBothPromptSources(t *testing.T) {
	cmd, _ := newRootCmd()
	if err := cmd.Flags().Set("file-prompt", "prompt.md"); err != nil {
		t.Fatalf("set file-prompt: %v", err)
	}

	err := runArgs(cmd, []string{"different", "task"})
	assertUsageError(t, err)
}

func TestRunArgsRejectsMissingPromptSource(t *testing.T) {
	cmd, _ := newRootCmd()
	assertUsageError(t, runArgs(cmd, nil))
}

func assertUsageError(t *testing.T, err error) {
	t.Helper()
	var usageErr sandbox.UsageError
	if !errors.As(err, &usageErr) {
		t.Fatalf("error = %v, want UsageError", err)
	}
}
