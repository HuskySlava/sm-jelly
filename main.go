package main

import (
	"cmp"
	"flag"
	"fmt"
	"github.com/HuskySlava/sm-jelly/internal/config"
	"log/slog"
	"os"
)

type Flags struct {
	configPath string
	logLevel   slog.Level
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

func main() {
	flags := parseFlags()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: flags.logLevel,
	})))

	cfg, err := config.Load(flags.configPath)
	if err != nil {
		slog.Error("Unable to load config:", "err", err)
		panic(err)
	}

	fmt.Println(cfg)
}
