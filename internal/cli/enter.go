package cli

import (
	"fmt"
	"os"
	"strings"

	gocli "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"

	"github.com/89luca89/distrobox/pkg/commands"
	"github.com/89luca89/distrobox/pkg/ui"
)

type EnterCmd struct {
	Name            string   `cli:"name,n" help:"name for the distrobox"`
	DryRun          bool     `cli:"dry-run,d" help:"only print the container manager command generated"`
	CleanPath       bool     `cli:"clean-path,c" help:"use a clean PATH inside the container"`
	AdditionalFlags string   `cli:"additional-flags,a" help:"additional flags to pass to the container manager command"`
	NoTTY           bool     `cli:"no-tty,T" help:"disable TTY allocation"`
	NoWorkDir       bool     `cli:"no-workdir,nw" help:"always start the container from container home directory"`
	Args            []string `arg:"" help:"[container-name] [-- command...]"`
	gocli.Base
}

func (c *EnterCmd) Run() error {
	containerName := c.Name
	args := c.Args

	if containerName == "" && len(args) > 0 {
		containerName = args[0]
		args = args[1:]
	}

	var customCommand string
	if len(args) > 0 {
		customCommand = strings.Join(args, " ")
	}

	options := commands.EnterOptions{
		ContainerName:   containerName,
		AdditionalFlags: c.AdditionalFlags,
		CustomCommand:   customCommand,
		DryRun:          c.DryRun,
		NoTTY:           c.NoTTY,
		CleanPath:       c.CleanPath,
	}

	progress := ui.NewProgress(os.Stderr)
	printer := ui.NewPrinter(os.Stderr, true)

	enterCmd := commands.NewEnterCommand(containerManager, progress, printer)
	_, err := enterCmd.Execute(c.Ctx, options)
	if err != nil {
		return fmt.Errorf("failed to execute enter command: %w", err)
	}

	return nil
}
