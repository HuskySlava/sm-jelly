package main

import (
	"cmp"
	"flag"
	"fmt"
	"github.com/HuskySlava/sm-jelly/internal/claude"
	"github.com/HuskySlava/sm-jelly/internal/config"
	"github.com/HuskySlava/sm-jelly/internal/runner"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Flags struct {
	configPath string
	logLevel   slog.Level
}

func main() {
	// Gracefully shutdown on SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Handle flags
	flags := parseFlags()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: flags.logLevel,
	})))

	// Load config
	cfg, err := config.Load(flags.configPath)
	if err != nil {
		slog.Error("Unable to load config", "err", err)
		panic(err)
	}

	// Init claude
	c := claude.New(&claude.Config{
		RunDir:  cfg.ClaudeRunDir,
		Timeout: time.Duration(cfg.ClaudeTimeoutSeconds) * time.Second,
	})

	// Create prompt jobs based on config
	var jobs []runner.Job
	for i, cj := range cfg.Jobs {
		j := runner.NewJob(cj.ClaudeTaskID, cj.ClaudeTaskCronSchedule, func() {
			r, err := c.Prompt("Edit index.html, improve it slightly add a bit of code - it should be small update that improves it")
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := os.WriteFile("./test.json", r, 0644); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Test %d \n", i)
		})
		jobs = append(jobs, j)
	}

	// Init job runner
	r, err := runner.New(jobs)
	if err != nil {
		slog.Error("failed to create runner", "err", err)
		panic(err)
	}
	err = r.Run()
	if err != nil {
		slog.Error("failed to start runner", "err", err)
		panic(err)
	}

	// Wait for SIGTERM
	<-quit
	err = r.Stop()
	if err != nil {
		slog.Error("failed to stop runner", "err", err)
	}
}

func parseFlags() *Flags {
	envConfigPath := os.Getenv("CONFIG_PATH")
	flagConfigPath := flag.String("config", envConfigPath, "config file path")

	var envLogLevel slog.Level
	if err := envLogLevel.UnmarshalText([]byte(os.Getenv("LOG_LEVEL"))); err != nil {
		envLogLevel = slog.LevelInfo
	}

	var logLevel slog.Level
	flag.TextVar(&logLevel, "level", envLogLevel, "log level")

	flag.Parse()

	configPath := cmp.Or(*flagConfigPath, "config.yaml")

	return &Flags{
		configPath: configPath,
		logLevel:   logLevel,
	}
}
