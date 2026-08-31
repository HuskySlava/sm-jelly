package main

import (
	"cmp"
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	configPath string
}

func parseFlags() *Flags {
	envConfigPath := os.Getenv("CONFIG_PATH")
	flagConfigPath := flag.String("config", envConfigPath, "config file path")
	flag.Parse()

	configPath := cmp.Or(*flagConfigPath, "config.yaml")

	return &Flags{
		configPath: configPath,
	}
}

func main() {
	flags := parseFlags()
	fmt.Println(flags)
}
