package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	claudeTaskCronSchedule string `yaml:"claude_tas_cro_schedule"`
	claudeTaskHandlerName  string `yaml:"claude_task_handler_name"`
}

func LoadConfig(path string) (*Config, error) {
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
	if cfg.claudeTaskCronSchedule == "" {
		return nil, fmt.Errorf("missing config value: claudeTaskCronSchedule")
	}

	if cfg.claudeTaskCronSchedule == "" {
		return nil, fmt.Errorf("missing config value: claudeTaskCronSchedule")
	}

	return cfg, nil
}
