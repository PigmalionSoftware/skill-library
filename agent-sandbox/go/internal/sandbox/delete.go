package sandbox

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"agent-sandbox/internal/git"
)

// WorktreeDelete removes the worktree checked out on a branch and deletes the
// branch with it, so the name is free for another run. Like WorktreeList it
// reads and writes git and nothing else, so it needs neither Docker nor the
// network.
//
// Deleting is not undoable, so every refusal happens before anything is
// removed: the command either does the whole job or leaves the repository
// exactly as it found it.
func WorktreeDelete(args []string, out io.Writer) error {
	branch, force, err := parseDeleteArgs(args)
	if err != nil {
		return err
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

	worktree, found := findWorktree(worktrees, branch)
	if !found {
		return fmt.Errorf("no worktree is checked out on branch %s; run worktree-list to see the ones there are", branch)
	}

	// The main worktree is the caller's own working copy, and the branch in it
	// is the one they are working on. Neither is the sandbox's to delete.
	if worktree.Main {
		return fmt.Errorf("%s is checked out in the main working copy, not in a sandbox worktree", branch)
	}

	// Removing the directory the command is running in would pull the ground
	// out from under it, and git refuses it anyway; saying so plainly beats
	// letting git's own message explain it.
	inside, err := isInside(executionDir, worktree.Path)
	if err != nil {
		return err
	}
	if inside {
		return fmt.Errorf("cannot delete %s from inside it; run this from the repository instead", worktree.Path)
	}

	if err := removeWorktree(repo, worktree.Path, force, out); err != nil {
		return err
	}

	if err := repo.DeleteBranch(branch); err != nil {
		return err
	}
	fmt.Fprintf(out, "Deleted branch %s\n", branch)
	return nil
}

// removeWorktree detaches one worktree from the repository: by removing its
// directory, or, when the directory is already gone, by pruning the record git
// kept of it.
func removeWorktree(repo *git.Repo, path string, force bool, out io.Writer) error {
	// A directory deleted by hand leaves the administrative record behind, and
	// there is nothing for git worktree remove to remove. Pruning is the way
	// back from that, and it is the only case where the directory is not what
	// gets deleted.
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := repo.PruneWorktrees(); err != nil {
			return err
		}
		fmt.Fprintf(out, "Pruned the record of %s, whose directory was already gone\n", path)
		return nil
	}

	// Uncommitted work is the other thing no reflog brings back. The check comes
	// before the removal so that a refusal costs nothing, and it says --force
	// where git's own refusal would arrive wrapped in an exit status.
	if !force {
		clean, err := repo.At(path).IsClean()
		if err != nil {
			return err
		}
		if !clean {
			return fmt.Errorf("%s has uncommitted changes; pass --force to delete it anyway", path)
		}
	}

	if err := repo.RemoveWorktree(path, force); err != nil {
		return err
	}
	fmt.Fprintf(out, "Removed worktree %s\n", path)
	return nil
}

// parseDeleteArgs reads the arguments of worktree-delete: the branch, named by
// -b because that is how the command is documented, and --force.
func parseDeleteArgs(args []string) (branch string, force bool, err error) {
	for len(args) > 0 {
		switch args[0] {
		case "-b", "--branch":
			if len(args) < 2 {
				return "", false, usageErrorf("%s needs a branch name", args[0])
			}
			if branch != "" {
				return "", false, usageErrorf("worktree-delete takes one branch")
			}
			branch, args = args[1], args[2:]
		case "--force":
			force, args = true, args[1:]
		default:
			return "", false, usageErrorf("unknown argument: %s", args[0])
		}
	}

	if branch == "" {
		return "", false, usageErrorf("-b <branch-name> is required")
	}
	return branch, force, nil
}

// findWorktree picks the worktree holding the branch. A branch is checked out
// in at most one worktree, so the first match is the only one.
func findWorktree(worktrees []git.Worktree, branch string) (git.Worktree, bool) {
	for _, worktree := range worktrees {
		if worktree.Branch == branch {
			return worktree, true
		}
	}
	return git.Worktree{}, false
}

// isInside reports whether dir is the directory at root or somewhere under it.
// Both are resolved first, because the worktree paths git reports have their
// symlinks resolved and the current directory may not.
func isInside(dir, root string) (bool, error) {
	dir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return false, err
	}

	// A root that is not there at all — the worktree whose directory was
	// deleted by hand — contains nothing, least of all the directory this
	// command is running in.
	root, err = filepath.EvalSymlinks(root)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	relative, err := filepath.Rel(root, dir)
	if err != nil {
		// The two are on different volumes, so one cannot be inside the other.
		return false, nil
	}

	// Rel answers "." for the directory itself and a path starting with ".."
	// for anything outside it; everything else is somewhere underneath.
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return false, nil
	}
	return true, nil
}
