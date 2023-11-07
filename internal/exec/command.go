package exec

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

func WithWorkdir(dir string) func(*exec.Cmd) {
	return func(c *exec.Cmd) {
		c.Dir = dir
	}
}

func WithStdin(r io.Reader) func(*exec.Cmd) {
	return func(c *exec.Cmd) {
		c.Stdin = r
	}
}

func WithStdout(w io.Writer) func(*exec.Cmd) {
	return func(c *exec.Cmd) {
		c.Stdout = w
	}
}

func ExecuteCommand(command string, args []string, opts ...func(*exec.Cmd)) error {
	cmd := exec.Command(command, args...)
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	for _, opt := range opts {
		opt(cmd)
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v: %v", err, stderr.String())
	}
	return nil
}
