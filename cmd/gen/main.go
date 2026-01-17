package main

import (
	"fmt"
	"io"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/ranna-go/ranna/internal/config"
)

type Args struct {
	CmdDotEnv *CmdDotEnv `arg:"subcommand:dotenv" help:"Generate default .env file"`
}

func (Args) Description() string {
	return "Tool to generate configs for ranna"
}

func main() {
	var args Args
	parser := arg.MustParse(&args)

	var err error
	switch {
	case args.CmdDotEnv != nil:
		err = args.CmdDotEnv.Run()
	default:
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

type CmdDotEnv struct {
	Output string `arg:"positional" default:".env" help:"Output file (or '-' to print it to stdout)"`
	Force  bool   `arg:"-f,--force" help:"Override file if it already exists"`
}

func (t *CmdDotEnv) Run() (err error) {
	var w io.Writer

	if t.Output == "-" {
		w = os.Stdout
	} else {
		if !t.Force {
			stat, err := os.Stat(t.Output)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed checking file %s: %s", t.Output, err)
			}
			if stat != nil {
				return fmt.Errorf("%s already exists", t.Output)
			}
		}

		f, err := os.Create(t.Output)
		if err != nil {
			return fmt.Errorf("failed creating file %s: %s", t.Output, err)
		}
		defer f.Close()
		w = f
	}

	return config.GenerateDotEnv(w)
}
