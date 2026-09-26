package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ClaudeJob struct {
	CronSchedule string `yaml:"cron_schedule"`
	JobID        string `yaml:"job_id"`
	JobPrompt    string `yaml:"prompt"`
}

type GitJob struct {
	CronSchedule string `yaml:"cron_schedule"`
	JobID        string `yaml:"job_id"`
	Command      string `yaml:"git_command"`
	Message      string `yaml:"git_message"`
	Scope        string `yaml:"git_scope"`
}

type Config struct {
	ClaudeModel          string      `yaml:"claude_model"`
	ClaudeJobs           []ClaudeJob `yaml:"claude_jobs"`
	GitJobs              []GitJob    `yaml:"git_jobs"`
	RunDir               string      `yaml:"run_dir"`
	ClaudeTimeoutSeconds int         `yaml:"claude_timeout_seconds"`
	GitTimeoutSeconds    int         `yaml:"git_timeout_seconds"`
}

func Load(path string) (*Config, error) {
	res, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from: %s, :%w", path, err)
	}

	var cfg *Config
	err = yaml.Unmarshal(res, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate config
	if cfg.ClaudeTimeoutSeconds < 5 {
		return nil, fmt.Errorf("claude_timeout_seconds value cannot be set below 5 seconds")
	}

	if cfg.GitTimeoutSeconds < 5 {
		return nil, fmt.Errorf("git_timeout_seconds value cannot be set below 5 seconds")
	}

	for i, cj := range cfg.ClaudeJobs {
		if cj.JobID == "" {
			return nil, fmt.Errorf("claude job #%d missing config value: job_id", i)
		}

		if cj.CronSchedule == "" {
			return nil, fmt.Errorf("claude job #%d missing config value: cron_schedule", i)
		}

		if cj.JobPrompt == "" {
			return nil, fmt.Errorf("claude job #%d missing config value: prompt", i)
		}
	}

	for i, gj := range cfg.GitJobs {
		if gj.JobID == "" {
			return nil, fmt.Errorf("git job #%d missing config value: job_id", i)
		}

		if gj.CronSchedule == "" {
			return nil, fmt.Errorf("git job #%d missing config value: cron_schedule", i)
		}
		if gj.Command != "pull" && gj.Command != "commit" && gj.Command != "push" {
			return nil, fmt.Errorf("git job #%d command \"%s\" not supported", i, gj.Command)
		}
		if gj.Command == "commit" && gj.Scope == "" {
			return nil, fmt.Errorf("git commit job #%d missing git_scope", i)
		}
		if gj.Command == "commit" && gj.Message == "" {
			return nil, fmt.Errorf("git commit job #%d missing git_message", i)
		}
	}

	return cfg, nil
}
