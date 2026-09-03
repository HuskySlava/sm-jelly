package git

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

type Config struct {
	RunDir  string
	Timeout time.Duration
}

type Git struct {
	Config *Config
}

func New(config *Config) *Git {
	return &Git{
		Config: config,
	}
}

func (g *Git) run(sub string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.Config.Timeout)
	defer cancel()

	a := append([]string{sub}, args...)
	cmd := exec.CommandContext(ctx, "git", a...)

	cmd.Dir = g.Config.RunDir
	cmd.WaitDelay = 5 * time.Second

	out, err := cmd.Output()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("git %s timed out after %s", sub, g.Config.Timeout)
		}
		if errors.Is(err, exec.ErrWaitDelay) {
			return nil, fmt.Errorf("git %s hangs past %s", sub, cmd.WaitDelay)
		}
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return nil, fmt.Errorf("git %s exited %d: stderr=%q stdout=%.300s", sub, ee.ExitCode(), ee.Stderr, out)
		}

		return nil, fmt.Errorf("git %s failed: %w", sub, err)
	}

	return out, nil
}

func (g *Git) Pull() error {
	_, err := g.run("pull")
	return err
}

func (g *Git) Add(scopes ...string) error {
	_, err := g.run("add", scopes...)
	return err
}

func (g *Git) Commit(message string) error {
	_, err := g.run("commit", "-m", message)
	return err
}

func (g *Git) Push() error {
	_, err := g.run("push")
	return err
}
