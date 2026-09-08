// Command agent-sandbox runs a coding agent inside the agent-sandbox container,
// on a git worktree of its own, so it never touches the current working copy.
//
// Usage:
//
//	agent-sandbox <branch> --agent <codex|claude|opencode|pi> [--model <model>] [--push] <prompt...>
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

	opts, err := sandbox.ParseArgs(os.Args[1:])
	if err != nil {
		var usageErr sandbox.UsageError
		if errors.As(err, &usageErr) {
			if message := usageErr.Error(); message != "" {
				fmt.Fprintf(os.Stderr, "error: %s\n", message)
			}
			fmt.Fprintln(os.Stderr, sandbox.Usage(program))
			os.Exit(sandbox.ExitUsage)
		}
		fail(err)
	}

	// An interrupt cancels the context, which stops the container and waits for
	// it: the agent holds the worktree open, so it must not outlive this
	// process. A second interrupt kills this process outright.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	status, err := sandbox.Run(ctx, opts, os.Stdout)
	if err != nil {
		fail(err)
	}

	// The agent's own status is the command's status.
	os.Exit(status)
}

// fail reports err and exits with the status it asks for, defaulting to a
// plain failure.
func fail(err error) {
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
