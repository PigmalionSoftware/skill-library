package sandbox

import (
	"fmt"
	"strings"

	"agent-sandbox/internal/agent"
)

// Options is one invocation of the sandbox.
type Options struct {
	Branch string
	Agent  agent.Agent
	Model  string
	Push   bool
	Prompt string
}

// Usage is the synopsis printed on a usage error, one line per form.
func Usage(program string) string {
	return fmt.Sprintf(
		"Usage: %s <branch-name> --agent <%s> [--model <model>] [--push] <prompt...>\n"+
			"       %s worktree-list",
		program, strings.Join(agent.Names(), "|"), program,
	)
}

// ParseArgs parses the arguments after the program name. The branch name comes
// first, then the options, then the prompt.
func ParseArgs(args []string) (Options, error) {
	if len(args) < 2 {
		return Options{}, UsageError{}
	}
	if strings.HasPrefix(args[0], "--") {
		return Options{}, usageErrorf("the branch name must come first, before any option")
	}

	opts := Options{Branch: args[0]}
	args = args[1:]
	agentName := ""

	for len(args) > 0 && strings.HasPrefix(args[0], "--") {
		switch args[0] {
		case "--agent":
			if len(args) < 2 {
				return Options{}, UsageError{}
			}
			agentName, args = args[1], args[2:]
		case "--model":
			if len(args) < 2 {
				return Options{}, UsageError{}
			}
			opts.Model, args = args[1], args[2:]
		case "--push":
			opts.Push, args = true, args[1:]
		default:
			return Options{}, usageErrorf("unknown option: %s", args[0])
		}
	}

	if len(args) < 1 {
		return Options{}, UsageError{}
	}
	opts.Prompt = strings.Join(args, " ")

	if agentName == "" {
		return Options{}, usageErrorf("--agent is required (%s)", agent.NamesProse())
	}

	selected, err := agent.Lookup(agentName)
	if err != nil {
		return Options{}, UsageError{err}
	}
	opts.Agent = selected

	if opts.Model == "" {
		opts.Model = selected.DefaultModel()
	}
	return opts, nil
}

// FullPrompt is the prompt handed to the agent: the user's instruction plus the
// house rules for a sandbox run.
func (o Options) FullPrompt() string {
	return o.Prompt + " do not use superpowerer or any spec skills, you decide all,  do not commit or push changes to git"
}
