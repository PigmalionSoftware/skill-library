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
func WorktreeList(out io.Writer) error {
	added, err := addedWorktrees()
	if err != nil {
		return err
	}

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

// WorktreeBranches names the branches the sandbox worktrees hold, for the shell
// completion of the flags that take one. A worktree on a detached HEAD has no
// branch to offer and is left out.
func WorktreeBranches() ([]string, error) {
	added, err := addedWorktrees()
	if err != nil {
		return nil, err
	}

	var branches []string
	for _, worktree := range added {
		if worktree.Branch != "" {
			branches = append(branches, worktree.Branch)
		}
	}
	return branches, nil
}

// addedWorktrees are the worktrees the sandbox added beside the caller's
// working copy. The main worktree is that working copy itself, so it is never
// one of them.
func addedWorktrees() ([]git.Worktree, error) {
	executionDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repo, err := git.Open(executionDir)
	if err != nil {
		return nil, err
	}

	worktrees, err := repo.Worktrees()
	if err != nil {
		return nil, err
	}

	var added []git.Worktree
	for _, worktree := range worktrees {
		if !worktree.Main {
			added = append(added, worktree)
		}
	}
	return added, nil
}

// branchColumn names the branch checked out in the worktree. A detached HEAD
// has no branch, and the column still needs something in it.
func branchColumn(worktree git.Worktree) string {
	if worktree.Branch == "" {
		return "-"
	}
	return worktree.Branch
}
