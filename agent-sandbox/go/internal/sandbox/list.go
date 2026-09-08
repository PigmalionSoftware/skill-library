package sandbox

import (
	"fmt"
	"io"
	"os"

	"agent-sandbox/internal/git"
)

// WorktreeList prints the worktrees attached to the repository the command was
// run in, one per line: the branch and the path. It reads git and nothing else,
// so it needs neither Docker nor the network.
func WorktreeList(args []string, out io.Writer) error {
	if len(args) > 0 {
		return usageErrorf("worktree-list takes no arguments")
	}

	executionDir, err := os.Getwd()
	if err != nil {
		return err
	}

	repo, err := git.Open(executionDir)
	if err != nil {
		return err
	}

	worktrees, err := repo.Worktrees()
	if err != nil {
		return err
	}

	// The main worktree is the caller's own working copy, not one the sandbox
	// added beside it, so it is not part of the list.
	var added []git.Worktree
	for _, worktree := range worktrees {
		if !worktree.Main {
			added = append(added, worktree)
		}
	}

	// The branch column is as wide as the longest branch name; nothing is
	// printed at all when there is no worktree to list.
	width := 0
	for _, worktree := range added {
		if length := len(branchColumn(worktree)); length > width {
			width = length
		}
	}
	for _, worktree := range added {
		fmt.Fprintf(out, "%-*s  %s\n", width, branchColumn(worktree), worktree.Path)
	}
	return nil
}

// branchColumn names the branch checked out in the worktree. A detached HEAD
// has no branch, and the column still needs something in it.
func branchColumn(worktree git.Worktree) string {
	if worktree.Branch == "" {
		return "-"
	}
	return worktree.Branch
}
