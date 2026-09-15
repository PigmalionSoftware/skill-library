package docker

import (
	"reflect"
	"testing"
)

func TestHostConfigKeepsExistingMountsWritableAndMakesAttachmentsReadOnly(t *testing.T) {
	config := hostConfig(RunOptions{Mounts: []Mount{
		{Host: "/host/worktree", Container: "/workspace"},
		{Host: "/host/image.png", Container: "/agent-sandbox-images/1.png", ReadOnly: true},
	}})
	want := []string{
		"/host/worktree:/workspace",
		"/host/image.png:/agent-sandbox-images/1.png:ro",
	}
	if !reflect.DeepEqual(config.Binds, want) {
		t.Fatalf("Binds = %q, want %q", config.Binds, want)
	}
}
