package main

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"

	"agent-sandbox/internal/sandbox"
)

func newRootCmd() (*cobra.Command, *int) {
	var (
		status        int
		branchName    string
		agentName     string
		model         string
		baseImage     string
		push          bool
		commitMessage string
		filePrompt    string
	)

	cmd := &cobra.Command{
		Use:   "agent-sandbox [-b <branch-name>] -a <agent> [flags] (<prompt...> | -f <prompt-file>)",
		Short: "Run a coding agent in a container, on a git worktree of its own",
		Long: "Run a coding agent inside the agent-sandbox container, on a git worktree of\n" +
			"its own, so it never touches the current working copy. Optionally supply the\n" +
			"worktree branch with -b or --branch; otherwise one is generated.\n\n" +
			"Supply the agent prompt directly as positional arguments or with -f or\n" +
			"--file-prompt; exactly one source is required.\n\n" +
			"Authentication comes from the agent's configuration directory on the host,\n" +
			"which is mounted into the container; no credentials are passed as environment\n" +
			"variables, so the agent must already be authenticated on the host.",
		Example: "  agent-sandbox -a codex \"fix the login redirect loop\"\n" +
			"  agent-sandbox -b fix-go-tests -a codex -i golang:1.26-alpine \"run go test ./...\"\n" +
			"  agent-sandbox --branch fix-login --agent claude --model sonnet --push \"add a test for it\"\n" +
			"  agent-sandbox -a codex -f prompt.md",

		// Args has to be set even where cobra's default would do, because a nil
		// Args makes cobra reject the first argument of a root command that has
		// subcommands as an unknown command — and here that argument is the
		// branch name.
		Args: runArgs,

		// fail is the only thing that prints an error or a synopsis, so that the
		// exit status and the message it goes with are decided in one place.
		SilenceUsage:  true,
		SilenceErrors: true,

		Version: buildVersion(),

		RunE: func(cmd *cobra.Command, args []string) error {
			// Everything after the branch name is the prompt, joined back into
			// the sentence it was before the shell took it apart, so that a
			// prompt of more than one word need not be quoted. A prompt holding
			// a word that starts with a dash does, or the flag parser claims it.
			opts, err := sandbox.NewOptions(sandbox.Options{
				Branch:        branchName,
				AgentName:     agentName,
				Model:         model,
				BaseImage:     baseImage,
				Push:          push,
				Prompt:        strings.Join(args, " "),
				CommitMessage: commitMessage,
				FilePrompt:    filePrompt,
			})
			if err != nil {
				return err
			}

			status, err = sandbox.Run(cmd.Context(), opts, cmd.OutOrStdout())
			return err
		},
	}

	cmd.Flags().StringVarP(&agentName, "agent", "a", "", "agent to run ("+strings.Join(sandbox.AgentNames(), "|")+")")
	cmd.Flags().StringVarP(&branchName, "branch", "b", "", "worktree branch name (default: generated)")
	cmd.Flags().StringVarP(&model, "model", "m", "", "model to use (default: the agent's own)")
	cmd.Flags().StringVarP(&baseImage, "base-image", "i", "", "Alpine base image for the agent sandbox (for example golang:1.26-alpine)")
	cmd.Flags().BoolVarP(&push, "push", "p", false, "commit the agent's work and push the branch")
	cmd.Flags().StringVarP(&commitMessage, "commit-message", "c", "", "commit message (default: resolved prompt)")
	cmd.Flags().StringVarP(&filePrompt, "file-prompt", "f", "", "path to a file containing the agent prompt")

	// pflag reports a malformed flag through the error func of the command it
	// was parsing, or of the nearest parent that has one, so this covers the
	// subcommands as well as the run.
	cmd.SetFlagErrorFunc(func(_ *cobra.Command, err error) error {
		return sandbox.NewUsageError(err)
	})

	completeFlag(cmd, "agent", func(string) ([]string, error) {
		return sandbox.AgentNames(), nil
	})

	cmd.AddCommand(newWorktreeListCmd(), newWorktreeDeleteCmd(), newWorktreeDeleteAllCmd())
	return cmd, &status
}

// runArgs requires exactly one prompt source. A file prompt removes shell
// quoting from longer instructions, but accepting it alongside positional text
// would silently discard the latter when Options resolves the file contents.
func runArgs(cmd *cobra.Command, args []string) error {
	filePrompt, err := cmd.Flags().GetString("file-prompt")
	if err != nil {
		return sandbox.NewUsageError(err)
	}

	hasPositionalPrompt := len(args) > 0
	hasFilePrompt := filePrompt != ""
	if !hasPositionalPrompt && !hasFilePrompt {
		return sandbox.UsageError{}
	}
	if hasPositionalPrompt && hasFilePrompt {
		return sandbox.NewUsageError(errors.New("prompt and --file-prompt cannot be used together"))
	}
	return nil
}

// usageArgs reports a wrong number of arguments as a malformed command line.
// Cobra hands the errors an argument validator produces straight back from
// Execute, with no hook of their own to pass them through, so the wrapping that
// SetFlagErrorFunc does for flags happens here for positional arguments.
func usageArgs(validate cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := validate(cmd, args); err != nil {
			return sandbox.NewUsageError(err)
		}
		return nil
	}
}

// completeFlag offers the values of a flag to the shell. The completions are
// filtered to what the user has typed so far and never fall back to file names,
// because none of these flags names a file. A registration only fails when the
// same flag is registered twice, which is a mistake in this file rather than
// anything a run can do, so there is nothing to report at run time.
func completeFlag(cmd *cobra.Command, name string, values func(prefix string) ([]string, error)) {
	_ = cmd.RegisterFlagCompletionFunc(name, func(_ *cobra.Command, _ []string, prefix string) ([]string, cobra.ShellCompDirective) {
		candidates, err := values(prefix)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var matching []string
		for _, candidate := range candidates {
			if strings.HasPrefix(candidate, prefix) {
				matching = append(matching, candidate)
			}
		}
		return matching, cobra.ShellCompDirectiveNoFileComp
	})
}
