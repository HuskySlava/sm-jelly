package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Job struct {
	ClaudeTaskID           string `yaml:"claude_task_id"`
	ClaudeTaskCronSchedule string `yaml:"claude_tas_cron_schedule"`
	ClaudeTaskHandlerName  string `yaml:"claude_task_handler_name"`
}

type Config struct {
	Jobs                 []Job  `yaml:"jobs"`
	ClaudeRunDir         string `yaml:"claude_run_dir"`
	ClaudeTimeoutSeconds int    `yaml:"claude_timeout_seconds"`
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
		return nil, fmt.Errorf("claude_timeout_seconds value cannot be set below 5 seconds: %w", err)
	}

	for i, j := range cfg.Jobs {
		if j.ClaudeTaskID == "" {
			return nil, fmt.Errorf("job #%d missing config value: claude_task_id", i)
		}

		if j.ClaudeTaskCronSchedule == "" {
			return nil, fmt.Errorf("job #%d missing config value: claude_tas_cron_schedule", i)
		}

		if j.ClaudeTaskHandlerName == "" {
			return nil, fmt.Errorf("job #%d, missing config value: claude_task_handler_name", i)
		}
	}

	return cfg, nil
}
