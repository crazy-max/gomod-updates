package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/crazy-max/gomod-updates/pkg/gomodupdates"
	"golang.org/x/term"
)

type cli struct {
	Update  bool                      `kong:"name='update',help='List only modules with updates.'"`
	Direct  bool                      `kong:"name='direct',help='List only direct modules.'"`
	Major   bool                      `kong:"name='major',help='Check for major-version module path candidates.'"`
	CI      bool                      `kong:"name='ci',help='Non-zero exit code when at least one outdated dependency was found.'"`
	Format  gomodupdates.OutputFormat `kong:"name='format',default='default',enum='default,markdown',help='Output format (${enum}).'"`
	Mod     string                    `kong:"name='mod',default='mod',help='Module download mode for go list calls.'"`
	Version kong.VersionFlag          `kong:"name='version',help='Print version information.'"`
}

var (
	name    = "gomod-updates"
	desc    = "Report available Go module updates, including major-version candidates"
	url     = "https://github.com/crazy-max/gomod-updates"
	version = "dev"
)

func main() {
	if err := run(); errors.Is(err, gomodupdates.ErrOutdated) {
		os.Exit(1)
	} else if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flags := cli{}
	parser, err := kong.New(&flags,
		kong.Name(name),
		kong.Description(fmt.Sprintf("%s. More info: %s", desc, url)),
		kong.UsageOnError(),
		kong.Vars{
			"version": version,
		},
		kong.ConfigureHelp(kong.HelpOptions{}))
	if err != nil {
		return err
	}

	log.SetFlags(0)

	_, err = parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	ctx := context.Background()
	in, err := input(ctx, os.Stdin, flags.Mod, gomodupdates.GoListModules)
	if err != nil {
		return err
	}

	return gomodupdates.Run(ctx, in, os.Stdout, gomodupdates.Options{
		Update: flags.Update,
		Direct: flags.Direct,
		Major:  flags.Major,
		CI:     flags.CI,
		Format: flags.Format,
		Lister: gomodupdates.GoVersionLister{
			Mod: flags.Mod,
		},
	})
}

func input(ctx context.Context, stdin *os.File, mod string, listModules func(context.Context, string, string) ([]byte, error)) (io.Reader, error) {
	if !term.IsTerminal(int(stdin.Fd())) {
		in := bufio.NewReader(stdin)
		if _, err := in.Peek(1); err == nil {
			return in, nil
		} else if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}
	out, err := listModules(ctx, "", mod)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(out), nil
}
