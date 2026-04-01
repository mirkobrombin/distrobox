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

type StopCmd struct {
	All            bool     `cli:"all,a" help:"stop all distroboxes"`
	Yes            bool     `cli:"yes,Y" help:"non-interactive, stop without asking"`
	ContainerNames []string `arg:"" help:"container names to stop"`
	gocli.Base
}

func (c *StopCmd) Run() error {
	options := &commands.StopOptions{
		ContainerNames: c.ContainerNames,
		NonInteractive: c.Yes,
		All:            c.All,
	}

	printer := ui.NewPrinter(os.Stdout, true)
	errPrinter := ui.NewPrinter(os.Stderr, true)
	prompter := ui.NewPrompter(*bufio.NewReader(os.Stdin), os.Stdout)

	stopCmd := commands.NewStopCommand(containerManager, prompter)
	err := stopCmd.Execute(c.Ctx, options)

	if errors.Is(err, commands.ErrStopAbortedByUserError) {
		printer.Println("Aborted.")
		return nil
	}

	if errors.Is(err, commands.ErrEmptyContainerList) {
		errPrinter.Println("No containers found.")
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to stop containers: %w", err)
	}

	return nil
}
