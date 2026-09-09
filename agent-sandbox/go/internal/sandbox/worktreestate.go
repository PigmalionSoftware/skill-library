package sandbox

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// worktreeRecord is one worktree the sandbox created.
type worktreeRecord struct {
	Repo    string    `json:"repo"`
	Path    string    `json:"path"`
	Branch  string    `json:"branch"`
	Created time.Time `json:"created"`
}

// stateFile is the list of worktrees the sandbox created, kept in the user's
// configuration directory so nothing is written inside the repository.
func stateFile() (string, error) {
	config, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(config, "agent-sandbox")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "worktrees.jsonl"), nil
}

// recordWorktree adds a worktree to the state file. A run only ever appends one
// line, so two runs starting at once cannot lose each other's entry.
func recordWorktree(record worktreeRecord) error {
	path, err := stateFile()
	if err != nil {
		return err
	}

	line, err := json.Marshal(record)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	if _, err := file.Write(append(line, '\n')); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

// recordedWorktrees are the worktrees the sandbox created in repoDir, in the
// order they were created.
func recordedWorktrees(repoDir string) ([]worktreeRecord, error) {
	records, err := readRecords()
	if err != nil {
		return nil, err
	}

	repoDir = resolvePath(repoDir)

	var mine []worktreeRecord
	for _, record := range records {
		if resolvePath(record.Repo) == repoDir {
			mine = append(mine, record)
		}
	}
	return mine, nil
}

// forgetWorktrees drops the given paths of repoDir from the state file. Lines
// belonging to other repositories are kept as they are, since nothing here can
// speak for them.
func forgetWorktrees(repoDir string, paths []string) error {
	if len(paths) == 0 {
		return nil
	}

	records, err := readRecords()
	if err != nil {
		return err
	}

	forget := make(map[string]bool, len(paths))
	for _, path := range paths {
		forget[resolvePath(path)] = true
	}

	repoDir = resolvePath(repoDir)

	var kept []worktreeRecord
	for _, record := range records {
		if resolvePath(record.Repo) == repoDir && forget[resolvePath(record.Path)] {
			continue
		}
		kept = append(kept, record)
	}
	return writeRecords(kept)
}

// readRecords is every line of the state file. A line that does not parse is
// skipped, so one bad write cannot make the worktree verbs unusable.
func readRecords() ([]worktreeRecord, error) {
	path, err := stateFile()
	if err != nil {
		return nil, err
	}

	// No file yet is an empty list, which is what a first run finds.
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []worktreeRecord
	lines := bufio.NewScanner(file)
	for lines.Scan() {
		var record worktreeRecord
		if err := json.Unmarshal(lines.Bytes(), &record); err != nil {
			continue
		}
		records = append(records, record)
	}
	return records, lines.Err()
}

// writeRecords replaces the state file. The new content is written beside it
// and renamed over it, so an interrupted write cannot leave the list truncated.
func writeRecords(records []worktreeRecord) error {
	path, err := stateFile()
	if err != nil {
		return err
	}

	temp, err := os.CreateTemp(filepath.Dir(path), "worktrees-*.jsonl")
	if err != nil {
		return err
	}

	for _, record := range records {
		line, err := json.Marshal(record)
		if err == nil {
			_, err = temp.Write(append(line, '\n'))
		}
		if err != nil {
			temp.Close()
			os.Remove(temp.Name())
			return err
		}
	}

	if err := temp.Close(); err != nil {
		os.Remove(temp.Name())
		return err
	}
	return os.Rename(temp.Name(), path)
}
