// Package sandbox creates the git worktree an agent works on and runs the
// agent against it inside the sandbox container.
package sandbox

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	agentsandbox "agent-sandbox"
	"agent-sandbox/internal/agent"
	"agent-sandbox/internal/docker"
	"agent-sandbox/internal/git"
)

// imageName is the tag of the sandbox image, shared by every agent.
const imageName = "agent-sandbox"

// workspace is where the worktree is mounted, matching the image's WORKDIR.
const workspace = "/workspace"

// semver matches the version inside an agent's --version output, which is
// rarely the bare number.
var semver = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)

// Run creates the worktree, runs the agent on it and optionally publishes the
// result. It returns the agent's own exit code.
func Run(ctx context.Context, opts Options, out io.Writer) (int, error) {
	executionDir, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	repo, err := git.Open(executionDir)
	if err != nil {
		return 0, err
	}

	client, err := docker.New(imageName, agentsandbox.Dockerfile)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	if err := ensureImageLatest(ctx, client, opts.Agent, out); err != nil {
		return 0, err
	}

	if err := git.CheckBranchName(opts.Branch); err != nil {
		// A bad branch name is a bad argument, so it exits like one, but git's
		// message needs no usage synopsis after it.
		return 0, StatusError{Status: ExitUsage, err: err}
	}

	// The worktree is a sibling of the directory the command was run from, not
	// of the repository root.
	worktreeDir := filepath.Join(filepath.Dir(executionDir), opts.Branch)
	if _, err := os.Stat(worktreeDir); err == nil {
		return 0, fmt.Errorf("worktree path already exists: %s", worktreeDir)
	}

	if err := os.MkdirAll(filepath.Dir(worktreeDir), 0o755); err != nil {
		return 0, err
	}
	if err := repo.AddWorktree(worktreeDir, opts.Branch, !repo.BranchExists(opts.Branch)); err != nil {
		return 0, err
	}

	runOpts, err := containerOptions(opts, worktreeDir)
	if err != nil {
		return 0, err
	}

	status, err := client.Run(ctx, runOpts)
	if err != nil {
		return 0, err
	}

	if err := publish(opts, worktreeDir, repo, out); err != nil {
		return 0, err
	}
	return status, nil
}

// containerOptions completes the agent's container configuration with the
// pieces that depend on this run: the workspace mount, the invoking user and
// the agent's command line.
func containerOptions(opts Options, worktreeDir string) (docker.RunOptions, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return docker.RunOptions{}, err
	}

	runOpts, err := opts.Agent.Container(home)
	if err != nil {
		return docker.RunOptions{}, err
	}

	runOpts.Entrypoint = opts.Agent.Binary()
	runOpts.Args = opts.Agent.Args(opts.Model, opts.FullPrompt())
	runOpts.User = strconv.Itoa(os.Getuid()) + ":" + strconv.Itoa(os.Getgid())
	runOpts.Mounts = append(runOpts.Mounts, docker.Mount{Host: worktreeDir, Container: workspace})

	return runOpts, nil
}

// publish commits and pushes what the agent produced, unless --push was left
// out, in which case the changes simply stay in the worktree.
func publish(opts Options, worktreeDir string, repo *git.Repo, out io.Writer) error {
	if !opts.Push {
		fmt.Fprintf(out, "Skipping commit and push. Changes left in %s\n", worktreeDir)
		return nil
	}

	worktree := repo.At(worktreeDir)
	if err := worktree.StageAll(); err != nil {
		return err
	}

	if !worktree.HasStagedChanges() {
		fmt.Fprintf(out, "Nothing to commit in %s\n", worktreeDir)
		return nil
	}

	// The commit message is the user's own prompt, without the sandbox rules
	// appended for the agent.
	if err := worktree.Commit(opts.Prompt); err != nil {
		return err
	}
	return worktree.Push(opts.Branch)
}

// ensureImageLatest builds the image when it is missing and rebuilds it when
// npm has a newer release of the agent than the one baked in.
func ensureImageLatest(ctx context.Context, client *docker.Client, target agent.Agent, out io.Writer) error {
	exists, err := client.ImageExists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		fmt.Fprintf(out, "Image %s not found. Building it.\n", imageName)
		if err := client.Build(ctx, "", ""); err != nil {
			return err
		}
	}

	installed, err := installedVersion(ctx, client, target)
	if err != nil {
		return err
	}
	latest, err := latestVersion(ctx, client, target)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%s in image: %s\n", target.Name(), installed)
	fmt.Fprintf(out, "%s latest:   %s\n", target.Name(), latest)

	if installed == latest {
		fmt.Fprintf(out, "%s is up to date.\n", target.Name())
		return nil
	}

	fmt.Fprintf(out, "Update available. Rebuilding %s with %s %s.\n", imageName, target.Name(), latest)
	if err := client.Build(ctx, target.BuildArg(), latest); err != nil {
		return err
	}

	rebuilt, err := installedVersion(ctx, client, target)
	if err != nil {
		return err
	}
	if rebuilt != latest {
		return fmt.Errorf("rebuilt %s is %s; expected %s", target.Name(), rebuilt, latest)
	}

	fmt.Fprintf(out, "%s rebuilt at version %s.\n", target.Name(), rebuilt)
	return nil
}

// installedVersion is the version the agent reports from inside the image.
func installedVersion(ctx context.Context, client *docker.Client, target agent.Agent) (string, error) {
	out, err := client.Output(ctx, target.Binary(), "--version")
	if err != nil {
		return "", err
	}

	version := semver.FindString(out)
	if version == "" {
		return "", fmt.Errorf("no version found in %s --version output: %q", target.Binary(), strings.TrimSpace(out))
	}
	return version, nil
}

// latestVersion is the version npm publishes for the agent's package.
func latestVersion(ctx context.Context, client *docker.Client, target agent.Agent) (string, error) {
	out, err := client.Output(ctx, "npm", "view", target.Package(), "version")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
