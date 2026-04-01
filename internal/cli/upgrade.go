package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	gocli "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"

	"github.com/89luca89/distrobox/pkg/commands"
	"github.com/89luca89/distrobox/pkg/ui"
)

type UpgradeCmd struct {
	All            bool     `cli:"all,a" help:"upgrade all distroboxes"`
	Running        bool     `cli:"running" help:"upgrade only running distroboxes (requires --all)"`
	Yes            bool     `cli:"yes,Y" help:"non-interactive, upgrade without asking"`
	ContainerNames []string `arg:"" help:"container names to upgrade"`
	gocli.Base
}

func (c *UpgradeCmd) Run() error {
	options := &commands.UpgradeOptions{
		ContainerNames: c.ContainerNames,
		All:            c.All,
		Running:        c.Running,
		NonInteractive: c.Yes,
	}

	printer := ui.NewPrinter(os.Stdout, true)
	errPrinter := ui.NewPrinter(os.Stderr, true)
	progress := ui.NewProgress(os.Stderr)
	prompter := ui.NewPrompter(*bufio.NewReader(os.Stdin), os.Stdout)

	upgradeCmd := commands.NewUpgradeCommand(containerManager, progress, printer, prompter)
	err := upgradeCmd.Execute(c.Ctx, options)

	if errors.Is(err, commands.ErrUpgradeAbortedByUser) {
		printer.Println("Aborted.")
		return nil
	}

	if errors.Is(err, commands.ErrEmptyContainerList) {
		errPrinter.Println("No containers found.")
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to upgrade containers: %w", err)
	}

	return nil
}
