package agent

import (
	"path/filepath"

	"agent-sandbox/internal/docker"
)

// claude runs Claude Code. Its HOME is redirected so the mounted config and
// credentials are the only ones it sees.
type claude struct{}

func (claude) Name() string         { return "claude" }
func (claude) Binary() string       { return "claude" }
func (claude) Package() string      { return "@anthropic-ai/claude-code" }
func (claude) BuildArg() string     { return "CLAUDE_CODE_VERSION" }
func (claude) DefaultModel() string { return "opus" }

func (c claude) Container(home string) (docker.RunOptions, error) {
	claudeHome := filepath.Join(home, ".claude")
	if err := requireDir(claudeHome, c.Name()); err != nil {
		return docker.RunOptions{}, err
	}

	claudeConfig := filepath.Join(home, ".claude.json")
	if err := requireFile(claudeConfig, c.Name()); err != nil {
		return docker.RunOptions{}, err
	}

	return docker.RunOptions{
		Interactive: true,
		Env:         []string{"HOME=/claude-home"},
		Mounts: []docker.Mount{
			{Host: claudeHome, Container: "/claude-home/.claude"},
			{Host: claudeConfig, Container: "/claude-home/.claude.json"},
		},
	}, nil
}

func (claude) Args(model, prompt string) []string {
	args := []string{"--print", "--permission-mode", "bypassPermissions", "--effort", "high"}
	if model != "" {
		args = append(args, "--model", model)
	}
	return append(args, prompt)
}
