package cli

import (
	"bufio"
	"fmt"
	"os"

	gocli "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"

	"github.com/89luca89/distrobox/pkg/commands"
	"github.com/89luca89/distrobox/pkg/ui"
)

type RmCmd struct {
	All            bool     `cli:"all,a" help:"delete all distroboxes"`
	Force          bool     `cli:"force,f" help:"force deletion"`
	Yes            bool     `cli:"yes,Y" help:"non-interactive mode"`
	RmHome         bool     `cli:"rm-home" help:"Remove container home directory"`
	ContainerNames []string `arg:"" help:"container names to remove"`
	gocli.Base
}

func (c *RmCmd) Run() error {
	options := commands.RmOptions{
		NoTTY:          c.Yes,
		Force:          c.Force,
		All:            c.All,
		RemoveHome:     c.RmHome,
		ContainerNames: c.ContainerNames,
	}

	prompter := ui.NewPrompter(*bufio.NewReader(os.Stdin), os.Stdout)

	rmCmd := commands.NewRmCommand(containerManager, prompter)
	_, err := rmCmd.Execute(c.Ctx, options)
	if err != nil {
		return fmt.Errorf("failed to execute rm command: %w", err)
	}

	return nil
}
