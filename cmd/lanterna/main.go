package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/rs/zerolog"
)

type Args struct {
	log zerolog.Logger

	ConfigPath string `arg:"-c,--config" placeholder:"PATH" help:"Path to configuration file" default:"config.json"`

	Collect *CollectCmd `arg:"subcommand:collect" help:"Collect IP addresses and print them"`
	Init    *InitCmd    `arg:"subcommand:init" help:"Create a configuration file"`
	Run     *RunCmd     `arg:"subcommand:run" help:"Collect IP addresses and send them"`
}

type CollectCmd struct{}

type InitCmd struct{}

type RunCmd struct{}

func main() {
	wr := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true, TimeFormat: time.RFC3339}
	log := zerolog.New(wr).With().Timestamp().Logger()

	if err := drive(log); err != nil {
		fmt.Println("lanterna: error:", err)
		os.Exit(1)
	}
}

func drive(log zerolog.Logger) error {
	args := Args{log: log}
	p := arg.MustParse(&args)

	if p.Subcommand() == nil {
		return fmt.Errorf("missing command (try: lanterna -h)")
	}

	switch {
	case args.Collect != nil:
		return cmdCollect(args)
	case args.Init != nil:
		return cmdInit(args)
	case args.Run != nil:
		return cmdRun(args, collect, postJSON)
	default:
		return fmt.Errorf("unwired command: %s", p.SubcommandNames())
	}
}
