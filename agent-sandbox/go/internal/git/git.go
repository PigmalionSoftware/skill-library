// Package git wraps the git CLI calls the sandbox needs to set up a worktree
// for the agent and to publish its result.
package git

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Repo is a git repository addressed by one of its working directories.
type Repo struct {
	Dir string

	Stdout io.Writer
	Stderr io.Writer
}

// Open returns the repository containing dir, resolved to its top level.
func Open(dir string) (*Repo, error) {
	repo := &Repo{Dir: dir, Stdout: os.Stdout, Stderr: os.Stderr}

	root, err := repo.output("rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("%s is not inside a Git repository", dir)
	}

	repo.Dir = root
	return repo, nil
}

// At returns the same repository addressed through another working directory,
// such as a worktree created by AddWorktree.
func (r *Repo) At(dir string) *Repo {
	return &Repo{Dir: dir, Stdout: r.Stdout, Stderr: r.Stderr}
}

// CheckBranchName reports whether name is a valid branch name.
func CheckBranchName(name string) error {
	cmd := exec.Command("git", "check-ref-format", "--branch", name)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("invalid branch name: %s", name)
	}
	return nil
}

// BranchExists reports whether the branch is already present in the repository.
func (r *Repo) BranchExists(branch string) bool {
	return r.run("show-ref", "--verify", "--quiet", "refs/heads/"+branch) == nil
}

// AddWorktree checks the branch out at dir, creating the branch when create is set.
func (r *Repo) AddWorktree(dir, branch string, create bool) error {
	args := []string{"worktree", "add"}
	if create {
		args = append(args, "-b", branch, dir)
	} else {
		args = append(args, dir, branch)
	}

	if err := r.run(args...); err != nil {
		return fmt.Errorf("git worktree add: %w", err)
	}
	return nil
}

// StageAll stages every change in the working directory.
func (r *Repo) StageAll() error {
	if err := r.run("add", "-A"); err != nil {
		return fmt.Errorf("git add: %w", err)
	}
	return nil
}

// HasStagedChanges reports whether anything is staged for commit.
func (r *Repo) HasStagedChanges() bool {
	return r.run("diff", "--cached", "--quiet") != nil
}

// Commit commits the staged changes.
func (r *Repo) Commit(message string) error {
	if err := r.run("commit", "-m", message); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}
	return nil
}

// Push pushes the branch and sets its upstream to origin.
func (r *Repo) Push(branch string) error {
	if err := r.run("push", "--set-upstream", "origin", branch); err != nil {
		return fmt.Errorf("git push: %w", err)
	}
	return nil
}

func (r *Repo) run(args ...string) error {
	cmd := exec.Command("git", append([]string{"-C", r.Dir}, args...)...)
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	return cmd.Run()
}

// output runs git and returns its trimmed stdout.
func (r *Repo) output(args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", r.Dir}, args...)...)
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}
