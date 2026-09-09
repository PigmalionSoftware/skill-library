package main

import (
	"errors"

	"github.com/spf13/cobra"

	"agent-sandbox/internal/sandbox"
)

// newWorktreeListCmd is the verb that shows what the sandbox has left behind:
// the worktrees it added to this repository, one per line.
func newWorktreeListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "worktree-list",
		Short: "List the worktrees the sandbox added to this repository",
		Args:  usageArgs(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return sandbox.WorktreeList(cmd.OutOrStdout())
		},
	}
}

// newWorktreeDeleteCmd is the verb that cleans one of them up, freeing the
// branch name for another run.
func newWorktreeDeleteCmd() *cobra.Command {
	var (
		branch string
		force  bool
	)

	cmd := &cobra.Command{
		Use:   "worktree-delete -b <branch-name>",
		Short: "Remove a sandbox worktree and delete the branch with it",
		Args:  usageArgs(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if branch == "" {
				return sandbox.NewUsageError(errors.New("-b <branch-name> is required"))
			}
			return sandbox.WorktreeDelete(branch, force, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&branch, "branch", "b", "", "branch whose worktree to delete")
	cmd.Flags().BoolVar(&force, "force", false, "delete even if the worktree has uncommitted changes")

	completeFlag(cmd, "branch", func(string) ([]string, error) {
		return sandbox.WorktreeBranches()
	})

	return cmd
}
