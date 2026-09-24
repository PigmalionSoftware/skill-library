package main

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunAndResumeQuerySources(t *testing.T) {
	tests := []struct {
		name       string
		section    string
		config     string
		args       []string
		wantBranch string
		wantPrompt string
		wantCommit string
		wantPush   bool
		wantUsage  bool
	}{
		{
			name:       "JSON query and commit message stay separate",
			config:     `{"run":{"agent":"codex","query":"run go version","commit-message":"this is from json file"}}`,
			wantPrompt: "run go version",
			wantCommit: "this is from json file",
		},
		{
			name:      "commit message does not supply a query",
			config:    `{"run":{"agent":"codex","commit-message":"this is from json file"}}`,
			wantUsage: true,
		},
		{
			name:       "CLI query shorthand overrides JSON query and file prompt",
			config:     `{"run":{"agent":"codex","query":"JSON task","file-prompt":"prompt.md"}}`,
			args:       []string{"-q", "CLI task"},
			wantPrompt: "CLI task",
		},
		{
			name:       "long query flag overrides JSON query",
			config:     `{"run":{"agent":"codex","query":"JSON task"}}`,
			args:       []string{"--query", "long CLI task"},
			wantPrompt: "long CLI task",
		},
		{
			name:       "push shorthand keeps the JSON query",
			config:     `{"run":{"agent":"codex","query":"JSON task"}}`,
			args:       []string{"-p"},
			wantPrompt: "JSON task",
			wantPush:   true,
		},
		{
			name:      "run rejects positional prompt even with JSON query",
			config:    `{"run":{"agent":"codex","query":"JSON task"}}`,
			args:      []string{"CLI", "task"},
			wantUsage: true,
		},
		{
			name:      "resume rejects positional prompt even with JSON query",
			section:   "resume",
			config:    `{"resume":{"branch":"existing-branch","agent":"codex","query":"JSON continuation"}}`,
			args:      []string{"CLI", "task"},
			wantUsage: true,
		},
		{
			name:       "resume query shorthand overrides its JSON query and push remains a switch",
			section:    "resume",
			config:     `{"resume":{"branch":"existing-branch","agent":"codex","query":"JSON continuation"}}`,
			args:       []string{"-p", "-q", "CLI continuation"},
			wantBranch: "existing-branch",
			wantPrompt: "CLI continuation",
			wantPush:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := tt.section
			if section == "" {
				section = "run"
			}
			t.Chdir(t.TempDir())
			if err := os.WriteFile(localConfigFile, []byte(tt.config), 0o644); err != nil {
				t.Fatal(err)
			}

			var branch string
			var flags runFlags
			cmd := &cobra.Command{}
			cmd.Flags().StringVarP(&branch, "branch", "b", "", "")
			flags.bind(cmd)
			if err := cmd.Flags().Parse(tt.args); err != nil {
				t.Fatal(err)
			}
			promptArgs := cmd.Flags().Args()
			err := configRunArgs(section, &branch, &flags)(cmd, promptArgs)
			if tt.wantUsage {
				assertUsageError(t, err)
				return
			}
			if err != nil {
				t.Fatalf("run args: %v", err)
			}

			opts, err := flags.options(branch)
			if err != nil {
				t.Fatalf("resolve options: %v", err)
			}
			if opts.Prompt != tt.wantPrompt {
				t.Errorf("agent prompt = %q, want %q", opts.Prompt, tt.wantPrompt)
			}
			if opts.Branch != tt.wantBranch && tt.wantBranch != "" {
				t.Errorf("branch = %q, want %q", opts.Branch, tt.wantBranch)
			}
			if opts.CommitMessage != tt.wantCommit {
				t.Errorf("commit message = %q, want %q", opts.CommitMessage, tt.wantCommit)
			}
			if opts.Push != tt.wantPush {
				t.Errorf("push = %t, want %t", opts.Push, tt.wantPush)
			}
			if opts.FilePrompt != "" {
				t.Errorf("file prompt = %q, want empty", opts.FilePrompt)
			}
		})
	}
}
