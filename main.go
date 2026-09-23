package main

import (
	"cmp"
	"errors"
	"flag"
	"github.com/HuskySlava/sm-jelly/internal/claude"
	"github.com/HuskySlava/sm-jelly/internal/config"
	"github.com/HuskySlava/sm-jelly/internal/git"
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
		RunDir:  cfg.RunDir,
		Timeout: time.Duration(cfg.ClaudeTimeoutSeconds) * time.Second,
	})

	// Create prompt jobs based on config
	var jobs []runner.Job
	for _, cj := range cfg.ClaudeJobs {
		j := runner.NewJob(cj.JobID, cj.CronSchedule, func() {
			r, err := c.Prompt(cj.JobPrompt)
			if err != nil {
				slog.Error("Unable to prompt Claude", "JobID", cj.JobID, "err", err)
				return
			}
			slog.Info("Job result", "JobID", cj.JobID, "result", r)
		})
		jobs = append(jobs, j)
	}

	g := git.New(&git.Config{
		RunDir:  cfg.RunDir,
		Timeout: time.Duration(cfg.GitTimeoutSeconds) * time.Second,
	})

	for _, gj := range cfg.GitJobs {
		j := runner.NewJob(gj.JobID, gj.CronSchedule, func() {
			switch gj.Command {
			case "pull":
				if err := g.Pull(); err != nil {
					slog.Error("Unable to git pull", "JobID", gj.JobID, "err", err)
				}
			case "commit":
				if err := g.Add(gj.Scope); err != nil {
					slog.Error("Unable to git add", "JobID", gj.JobID, "err", err)
					return
				}
				if err := g.Commit(gj.Message); err != nil {
					if errors.Is(err, git.ErrNothingToCommit) {
						slog.Info("Nothing to commit", "JobID", gj.JobID)
						return
					}
					slog.Error("Unable to git commit", "JobID", gj.JobID, "err", err)
				}
			case "push":
				if err := g.Push(); err != nil {
					slog.Error("Unable to git push", "JobID", gj.JobID, "err", err)
				}
			}
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
