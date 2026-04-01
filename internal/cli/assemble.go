package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"

	gocli "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"

	"github.com/89luca89/distrobox/pkg/commands"
	"github.com/89luca89/distrobox/pkg/manifest"
	"github.com/89luca89/distrobox/pkg/ui"
)

type AssembleCmd struct {
	Create AssembleCreateCmd `cmd:"create" help:"Create distroboxes from manifest"`
	Rm     AssembleRmCmd     `cmd:"rm" help:"Remove distroboxes from manifest"`
	gocli.Base
}

type AssembleCreateCmd struct {
	File    string `cli:"file" help:"path to manifest file" default:"./distrobox.ini"`
	Name    string `cli:"name,n" help:"run against a single entry"`
	Replace bool   `cli:"replace,R" help:"replace existing distroboxes"`
	DryRun  bool   `cli:"dry-run,d" help:"only print the container manager command generated"`
	gocli.Base
}

func (c *AssembleCreateCmd) Run() error {
	return assembleRun(c.Ctx, c.File, c.Name, c.DryRun, false, c.Replace)
}

type AssembleRmCmd struct {
	File   string `cli:"file" help:"path to manifest file" default:"./distrobox.ini"`
	Name   string `cli:"name,n" help:"run against a single entry"`
	DryRun bool   `cli:"dry-run,d" help:"only print the container manager command generated"`
	gocli.Base
}

func (c *AssembleRmCmd) Run() error {
	return assembleRun(c.Ctx, c.File, c.Name, c.DryRun, true, false)
}

func assembleRun(ctx context.Context, filePath, name string, dryRun, doDelete, replace bool) error {
	if filePath == "" {
		filePath = "./distrobox.ini"
	}

	parsed, err := manifest.Parse(ctx, filePath)
	if err != nil {
		return fmt.Errorf("failed to parse manifest file: %w", err)
	}

	opts := commands.AssembleOptions{
		Items:   parsed,
		Boxname: name,
		DryRun:  dryRun,
	}
	if doDelete {
		opts.Delete = true
	} else {
		opts.Replace = replace
	}

	prompter := ui.NewPrompter(*bufio.NewReader(os.Stdin), os.Stdout)
	progress := ui.NewProgress(os.Stderr)
	printer := ui.NewPrinter(os.Stdout, true)

	assembleCmd := commands.NewAssembleCommand(containerManager, prompter, progress, printer)
	err = assembleCmd.Execute(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to execute assemble command: %w", err)
	}
	return nil
}
