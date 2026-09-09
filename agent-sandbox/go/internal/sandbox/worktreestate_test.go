package sandbox

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestStateFile(t *testing.T) {
	configDir := setupWorktreeState(t)

	path, err := stateFile()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(configDir, "agent-sandbox", "worktrees.jsonl")
	if path != want {
		t.Fatalf("stateFile() = %q, want %q", path, want)
	}
	if info, err := os.Stat(filepath.Dir(path)); err != nil || !info.IsDir() {
		t.Fatalf("state directory = (%v, %v), want existing directory", info, err)
	}
}

func TestRecordWorktreeAndReadRecords(t *testing.T) {
	setupWorktreeState(t)

	empty, err := readRecords()
	if err != nil {
		t.Fatal(err)
	}
	if empty != nil {
		t.Fatalf("readRecords() before writing = %#v, want nil", empty)
	}

	records := []worktreeRecord{
		{Repo: "/repos/one", Path: "/worktrees/one", Branch: "one", Created: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)},
		{Repo: "/repos/two", Path: "/worktrees/two", Branch: "two", Created: time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)},
	}
	for _, record := range records {
		if err := recordWorktree(record); err != nil {
			t.Fatal(err)
		}
	}

	got, err := readRecords()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, records) {
		t.Fatalf("readRecords() = %#v, want %#v", got, records)
	}
}

func TestReadRecordsSkipsMalformedLines(t *testing.T) {
	setupWorktreeState(t)
	record := worktreeRecord{Repo: "/repo", Path: "/worktree", Branch: "feature"}
	if err := recordWorktree(record); err != nil {
		t.Fatal(err)
	}

	path, err := stateFile()
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("not json\n"); err != nil {
		file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	got, err := readRecords()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []worktreeRecord{record}) {
		t.Fatalf("readRecords() = %#v, want only %#v", got, record)
	}
}

func TestRecordedWorktreesFiltersByResolvedRepository(t *testing.T) {
	setupWorktreeState(t)
	repoDir := t.TempDir()
	repoLink := filepath.Join(t.TempDir(), "repo")
	if err := os.Symlink(repoDir, repoLink); err != nil {
		t.Fatal(err)
	}

	mine := []worktreeRecord{
		{Repo: repoDir, Path: "/worktrees/one", Branch: "one"},
		{Repo: repoLink, Path: "/worktrees/two", Branch: "two"},
	}
	for _, record := range append(mine, worktreeRecord{Repo: t.TempDir(), Path: "/worktrees/other", Branch: "other"}) {
		if err := recordWorktree(record); err != nil {
			t.Fatal(err)
		}
	}

	got, err := recordedWorktrees(repoDir)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, mine) {
		t.Fatalf("recordedWorktrees() = %#v, want %#v", got, mine)
	}
}

func TestForgetWorktreesRemovesOnlyMatchingRepositoryAndPath(t *testing.T) {
	setupWorktreeState(t)
	repoDir := t.TempDir()
	pathToRemove := filepath.Join(t.TempDir(), "remove")
	pathToKeep := filepath.Join(t.TempDir(), "keep")
	records := []worktreeRecord{
		{Repo: repoDir, Path: pathToRemove, Branch: "remove"},
		{Repo: repoDir, Path: pathToKeep, Branch: "keep"},
		{Repo: t.TempDir(), Path: pathToRemove, Branch: "other-repo"},
	}
	for _, record := range records {
		if err := recordWorktree(record); err != nil {
			t.Fatal(err)
		}
	}

	if err := forgetWorktrees(repoDir, []string{pathToRemove}); err != nil {
		t.Fatal(err)
	}

	got, err := readRecords()
	if err != nil {
		t.Fatal(err)
	}
	want := []worktreeRecord{records[1], records[2]}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("readRecords() after forget = %#v, want %#v", got, want)
	}
}

func TestForgetWorktreesDoesNothingWithoutPaths(t *testing.T) {
	setupWorktreeState(t)
	record := worktreeRecord{Repo: "/repo", Path: "/worktree", Branch: "feature"}
	if err := recordWorktree(record); err != nil {
		t.Fatal(err)
	}

	if err := forgetWorktrees(record.Repo, nil); err != nil {
		t.Fatal(err)
	}
	got, err := readRecords()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []worktreeRecord{record}) {
		t.Fatalf("readRecords() = %#v, want %#v", got, record)
	}
}

func TestWriteRecordsReplacesState(t *testing.T) {
	setupWorktreeState(t)
	if err := recordWorktree(worktreeRecord{Repo: "/repo", Path: "/worktree/old", Branch: "old"}); err != nil {
		t.Fatal(err)
	}

	want := []worktreeRecord{{Repo: "/repo", Path: "/worktree/new", Branch: "new"}}
	if err := writeRecords(want); err != nil {
		t.Fatal(err)
	}
	got, err := readRecords()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("readRecords() after replacement = %#v, want %#v", got, want)
	}

	if err := writeRecords(nil); err != nil {
		t.Fatal(err)
	}
	got, err = readRecords()
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("readRecords() after clearing = %#v, want nil", got)
	}
}

func setupWorktreeState(t *testing.T) string {
	t.Helper()

	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)
	return configDir
}
