package git

import (
	"fmt"
	"strings"
)

// Worktree is one working directory attached to the repository.
type Worktree struct {
	// Path is the absolute path of the working directory.
	Path string
	// Branch is the short branch name, empty when HEAD is detached.
	Branch string
	// Main marks the repository's own working copy, as opposed to one added
	// beside it.
	Main bool
}

// Worktrees lists every working directory attached to the repository, the main
// one first. It answers the same from any of them, including from inside a
// worktree added by AddWorktree.
func (r *Repo) Worktrees() ([]Worktree, error) {
	out, err := r.output("worktree", "list", "--porcelain")
	if err != nil {
		return nil, fmt.Errorf("git worktree list: %w", err)
	}
	return parseWorktrees(out), nil
}

// parseWorktrees reads git worktree list --porcelain: one record per worktree,
// separated by a blank line, each a list of "label value" lines. Only the path
// and the branch are taken; the commit and the flags git offers alongside them
// are left alone, as is any label a newer git adds.
func parseWorktrees(out string) []Worktree {
	var list []Worktree

	for record := range strings.SplitSeq(out, "\n\n") {
		var worktree Worktree

		for line := range strings.SplitSeq(record, "\n") {
			label, value, _ := strings.Cut(strings.TrimSpace(line), " ")
			switch label {
			case "worktree":
				worktree.Path = value
			case "branch":
				worktree.Branch = strings.TrimPrefix(value, "refs/heads/")
			}
		}

		// A record with no worktree line is not a worktree.
		if worktree.Path != "" {
			list = append(list, worktree)
		}
	}

	// Git lists the repository's own working copy first, always.
	if len(list) > 0 {
		list[0].Main = true
	}
	return list
}
