package agent

import (
	"path/filepath"

	"agent-sandbox/internal/docker"
)

// codex runs OpenAI Codex. It needs a TTY, so the container is started with -it.
type codex struct{}

func (codex) Name() string         { return "codex" }
func (codex) Binary() string       { return "codex" }
func (codex) Package() string      { return "@openai/codex" }
func (codex) BuildArg() string     { return "CODEX_VERSION" }
func (codex) DefaultModel() string { return "gpt-5.6-terra" }
func (codex) SupportsImages() bool { return true }

func (c codex) Container(home string) (docker.RunOptions, error) {
	codexHome := filepath.Join(home, ".codex")
	if err := requireDir(codexHome, c.Name()); err != nil {
		return docker.RunOptions{}, err
	}

	return docker.RunOptions{
		TTY: true,
		Env: []string{"CODEX_HOME=/codex-home"},
		Mounts: []docker.Mount{
			{Host: codexHome, Container: "/codex-home"},
		},
	}, nil
}

func (codex) Args(model, prompt string, images []string) []string {
	args := []string{"--dangerously-bypass-approvals-and-sandbox"}
	if model != "" {
		args = append(args, "--model", model)
	}
	args = append(args, "--config", `model_reasoning_effort="high"`, "exec")
	for _, image := range images {
		args = append(args, "-i", image)
	}
	// Codex accepts one or more values after -i. The option terminator keeps the
	// text prompt from being consumed as another attachment.
	if len(images) > 0 {
		args = append(args, "--")
	}
	return append(args, prompt)
}
