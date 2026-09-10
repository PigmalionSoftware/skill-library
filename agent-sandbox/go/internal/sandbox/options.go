package sandbox

import (
	"agent-sandbox/internal/agent"
)

// Options is one invocation of the sandbox.
type Options struct {
	Branch    string
	Agent     agent.Agent
	Model     string
	BaseImage string
	Push      bool
	Prompt    string
}

// NewOptions turns the values the command line carried into one invocation,
// resolving the agent name and filling in the model the agent defaults to. The
// command line itself is parsed by the caller; what is left here is the part
// that needs to know the agents, which is why an unusable name comes back as a
// UsageError rather than as a plain failure.
func NewOptions(branch, agentName, model, baseImage, prompt string, push bool) (Options, error) {
	if agentName == "" {
		return Options{}, usageErrorf("--agent is required (%s)", agent.NamesProse())
	}

	selected, err := agent.Lookup(agentName)
	if err != nil {
		return Options{}, UsageError{err}
	}
	// An empty model is not an omission to report: it means the agent resolves
	// the model itself, and DefaultModel says so by returning an empty string.
	if model == "" {
		model = selected.DefaultModel()
	}

	return Options{
		Branch:    branch,
		Agent:     selected,
		Model:     model,
		BaseImage: baseImage,
		Push:      push,
		Prompt:    prompt,
	}, nil
}

// AgentNames are the agent names the --agent flag accepts, in the order they
// are offered. The command line needs them for its help text and its
// completions, and this keeps it from having to know the agent package.
func AgentNames() []string {
	return agent.Names()
}

// FullPrompt is the prompt handed to the agent: the user's instruction plus the
// house rules for a sandbox run.
func (o Options) FullPrompt() string {
	return o.Prompt + " do not use superpowerer or any spec skills, you decide all,  do not commit or push changes to git"
}
