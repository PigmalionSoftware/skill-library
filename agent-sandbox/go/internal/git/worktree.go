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

// RemoveWorktree detaches the worktree at dir and deletes its directory. Git
// refuses a worktree with uncommitted changes unless force is set; the caller
// asks the same question first, so that refusal is a backstop rather than the
// message a user sees.
func (r *Repo) RemoveWorktree(dir string, force bool) error {
	args := []string{"worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	// The path goes after --, so a directory whose name begins with a dash is
	// still read as a path and not as an option.
	args = append(args, "--", dir)

	if err := r.run(args...); err != nil {
		return fmt.Errorf("git worktree remove: %w", err)
	}
	return nil
}

// PruneWorktrees forgets the worktrees whose directories are no longer on disk.
// It is the way back from a worktree deleted by hand, where the administrative
// record outlives the directory and RemoveWorktree has nothing left to remove.
func (r *Repo) PruneWorktrees() error {
	if err := r.run("worktree", "prune"); err != nil {
		return fmt.Errorf("git worktree prune: %w", err)
	}
	return nil
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
