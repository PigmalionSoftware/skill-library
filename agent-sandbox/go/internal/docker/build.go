package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/pkg/jsonmessage"
	"golang.org/x/term"
)

const dockerfileName = "Dockerfile"

// Build builds the image. buildArg, when both it and value are non-empty, pins
// one Dockerfile ARG so a single agent can be upgraded.
func (c *Client) Build(ctx context.Context, buildArg, value string) error {
	buildContext, err := c.buildContext()
	if err != nil {
		return err
	}

	opts := build.ImageBuildOptions{
		Tags:       []string{c.image},
		Dockerfile: dockerfileName,
		Remove:     true,
	}
	if buildArg != "" && value != "" {
		opts.BuildArgs = map[string]*string{buildArg: &value}
	}

	resp, err := c.api.ImageBuild(ctx, buildContext, opts)
	if err != nil {
		return fmt.Errorf("building image %s: %w", c.image, err)
	}
	defer resp.Body.Close()

	if err := c.displayProgress(resp.Body); err != nil {
		return fmt.Errorf("building image %s: %w", c.image, err)
	}
	return nil
}

// buildContext tars the Dockerfile. It is the whole context: the Dockerfile
// copies nothing from the host, so a single entry is a complete build context.
func (c *Client) buildContext() (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	header := &tar.Header{
		Name:    dockerfileName,
		Mode:    0o644,
		Size:    int64(len(c.dockerfile)),
		ModTime: time.Now(),
	}
	if err := tw.WriteHeader(header); err != nil {
		return nil, fmt.Errorf("writing build context: %w", err)
	}
	if _, err := tw.Write([]byte(c.dockerfile)); err != nil {
		return nil, fmt.Errorf("writing build context: %w", err)
	}
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("writing build context: %w", err)
	}

	return &buf, nil
}

// displayProgress renders the build's JSON message stream. It also reports the
// build failure itself, which the daemon sends as the last message rather than
// as an HTTP error.
func (c *Client) displayProgress(body io.Reader) error {
	fd, isTerminal := terminal(c.stderr)
	return jsonmessage.DisplayJSONMessagesStream(body, c.stderr, fd, isTerminal, nil)
}

// terminal reports the file descriptor behind w and whether it is a terminal,
// so progress can be rendered in place instead of line by line.
func terminal(w io.Writer) (uintptr, bool) {
	file, ok := w.(*os.File)
	if !ok {
		return 0, false
	}

	fd := file.Fd()
	return fd, term.IsTerminal(int(fd))
}
