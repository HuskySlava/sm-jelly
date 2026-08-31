package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ClaudeTaskCronSchedule string `yaml:"claude_tas_cron_schedule"`
	ClaudeTaskHandlerName  string `yaml:"claude_task_handler_name"`
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
	if cfg.ClaudeTaskCronSchedule == "" {
		return nil, fmt.Errorf("missing config value: ClaudeTaskCronSchedule")
	}

	if cfg.ClaudeTaskHandlerName == "" {
		return nil, fmt.Errorf("missing config value: ClaudeTaskHandlerName")
	}

	return cfg, nil
}
