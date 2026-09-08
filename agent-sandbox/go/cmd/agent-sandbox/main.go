// Command agent-sandbox runs a coding agent inside the agent-sandbox container,
// on a git worktree of its own, so it never touches the current working copy.
//
// Usage:
//
//	agent-sandbox <branch> --agent <codex|claude|opencode|pi> [--model <model>] [--push] <prompt...>
//	agent-sandbox worktree-list
//	agent-sandbox worktree-delete -b <branch> [--force]
//
// The branch name always comes first, before any option. Authentication comes
// from the agent's configuration directory on the host, which is mounted into
// the container; no credentials are passed as environment variables.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"agent-sandbox/internal/sandbox"
)

func main() {
	program := filepath.Base(os.Args[0])
	args := os.Args[1:]

	// The first argument is a verb when it names one; otherwise it is the
	// branch name of a run, the form this command started out with. The verbs
	// only ever talk to git, so they return before the signal handling and the
	// Docker client a run needs.
	if len(args) > 0 {
		switch args[0] {
		case "worktree-list":
			if err := sandbox.WorktreeList(args[1:], os.Stdout); err != nil {
				fail(program, err)
			}
			return
		case "worktree-delete":
			if err := sandbox.WorktreeDelete(args[1:], os.Stdout); err != nil {
				fail(program, err)
			}
			return
		}
	}

	opts, err := sandbox.ParseArgs(args)
	if err != nil {
		fail(program, err)
	}

	// An interrupt cancels the context, which stops the container and waits for
	// it: the agent holds the worktree open, so it must not outlive this
	// process. A second interrupt kills this process outright.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	status, err := sandbox.Run(ctx, opts, os.Stdout)
	if err != nil {
		fail(program, err)
	}

	// The agent's own status is the command's status.
	os.Exit(status)
}

// fail reports err and exits with the status it asks for, defaulting to a
// plain failure.
func fail(program string, err error) {
	// A malformed command line is the one error the synopsis can answer.
	var usageErr sandbox.UsageError
	if errors.As(err, &usageErr) {
		if message := usageErr.Error(); message != "" {
			fmt.Fprintf(os.Stderr, "error: %s\n", message)
		}
		fmt.Fprintln(os.Stderr, sandbox.Usage(program))
		os.Exit(sandbox.ExitUsage)
	}

	// An interrupt is not a failure to explain: it is what the user asked for,
	// and it exits the way a signalled process does.
	if errors.Is(err, context.Canceled) {
		fmt.Fprintln(os.Stderr, "interrupted; the container was stopped")
		os.Exit(sandbox.ExitInterrupted)
	}

	fmt.Fprintf(os.Stderr, "error: %s\n", err)

	status := sandbox.ExitFailure
	var statusErr sandbox.StatusError
	if errors.As(err, &statusErr) {
		status = statusErr.Status
	}
	os.Exit(status)
}
