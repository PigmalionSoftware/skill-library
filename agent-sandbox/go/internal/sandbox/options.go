package sandbox

import (
	"agent-sandbox/internal/agent"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"
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

	if branch == "" {
		randomBranch, err := generateBranchName()
		if err != nil {
			return Options{}, err
		}
		branch = randomBranch
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

// generateBranchName creates a lowercase-letter name with a six-digit suffix
// without a host dictionary or an additional package. A collision is reported
// by the sandbox rather than being silently retried.
func generateBranchName() (string, error) {
	const (
		letters    = "abcdefghijklmnopqrstuvwxyz"
		wordLength = 8
	)

	var word strings.Builder
	word.Grow(wordLength)
	letterLimit := big.NewInt(int64(len(letters)))
	for range wordLength {
		index, err := rand.Int(rand.Reader, letterLimit)
		if err != nil {
			return "", fmt.Errorf("generate random branch letters: %w", err)
		}
		word.WriteByte(letters[index.Int64()])
	}

	number, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", fmt.Errorf("generate random branch number: %w", err)
	}
	return word.String() + "-" + strconv.FormatInt(number.Int64()+100000, 10), nil
}
