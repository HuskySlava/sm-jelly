package claude

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	Timeout time.Duration
	RunDir  string
}

type Claude struct {
	Config *Config
}

func New(cfg *Config) *Claude {
	return &Claude{
		Config: cfg,
	}
}

func (c *Claude) Prompt(prompt string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Config.Timeout)
	defer cancel()

	cmd := exec.CommandContext(
		ctx, "claude",
		"-p", prompt,
		"--output-format", "json",
		"--permission-mode",
		"acceptEdits")

	cmd.Dir = c.Config.RunDir // Adjust

	// TODO: Remove / fix before prod
	cmd.Env = append(os.Environ(), "CLAUDE_CONFIG_DIR="+os.Getenv("HOME")+"/.claude-personal")

	cmd.WaitDelay = 5 * time.Second

	out, err := cmd.Output()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("claude timed out after %s", c.Config.Timeout)
		}
		if errors.Is(err, exec.ErrWaitDelay) {
			return nil, fmt.Errorf("claude hangs past %s", cmd.WaitDelay)
		}
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return nil, fmt.Errorf("claude exited %d: stderr=%q stdout=%.300s", ee.ExitCode(), ee.Stderr, out)
		}

		return nil, fmt.Errorf("unable to get output from claude: %w", err)
	}

	return out, nil
}
