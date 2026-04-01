package cli

import (
	"fmt"
	"os"

	gocli "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"

	"github.com/89luca89/distrobox/pkg/commands"
	"github.com/89luca89/distrobox/pkg/ui"
)

type ListCmd struct {
	NoColor bool `cli:"no-color" help:"Disable color output"`
	gocli.Base
}

func (c *ListCmd) Run() error {
	listCmd := commands.NewListCommand(containerManager)
	result, err := listCmd.Execute(c.Ctx)
	if err != nil {
		return fmt.Errorf("failed to execute list command: %w", err)
	}

	noColor := c.NoColor || !isTerminal()
	printResult(result, noColor)

	return nil
}

func printResult(result *commands.ListResult, noColor bool) {
	rowFormat := "%-12s | %-20s | %-18s | %-30s\n"

	//nolint:forbidigo // Using fmt.Printf is acceptable here for CLI output
	fmt.Printf(rowFormat, "ID", "NAME", "STATUS", "IMAGE")

	for _, cont := range result.Containers {
		var line string
		switch {
		case noColor:
			line = rowFormat
		case cont.IsRunning():
			line = ui.Green(rowFormat)
		default:
			line = ui.Yellow(rowFormat)
		}

		//nolint:forbidigo // Using fmt.Printf is acceptable here for CLI output
		fmt.Printf(line, cont.ID, cont.Name, cont.Status, cont.Image)
	}
}

func isTerminal() bool {
	stat, _ := os.Stdout.Stat()
	return (stat.Mode() & os.ModeCharDevice) != 0
}
