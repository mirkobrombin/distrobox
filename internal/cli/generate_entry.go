package cli

import (
	"fmt"
	"os"
	"path/filepath"

	gocli "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"

	"github.com/89luca89/distrobox/pkg/commands"
)

type GenerateEntryCmd struct {
	Delete bool   `cli:"delete,d" help:"delete the entry"`
	Icon   string `cli:"icon,i" help:"specify a custom icon" default:"auto"`
	All    bool   `cli:"all,a" help:"perform for all distroboxes"`
	Root   bool   `cli:"root,r" help:"perform on rootful distroboxes"`
	Name   string `arg:"" help:"container name"`
	gocli.Base
}

func (c *GenerateEntryCmd) Run() error {
	distroboxPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get distrobox executable path: %w", err)
	}

	listCmd := commands.NewListCommand(containerManager)

	opts := &commands.GenerateEntryOptions{
		Delete:              c.Delete,
		Root:                c.Root,
		DesktopEntryBaseDir: getDesktopEntryDir(),
		DistroboxPath:       distroboxPath,
	}
	if c.All {
		opts.All = true
	} else {
		opts.ContainerName = c.Name
		opts.Icon = c.Icon
	}

	genEntryCmd := commands.NewGenerateEntryCommand(listCmd)
	err = genEntryCmd.Execute(c.Ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to execute generate entry command: %w", err)
	}

	return nil
}

func getDesktopEntryDir() string {
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome == "" {
		home := os.Getenv("HOME")
		return filepath.Join(home, ".local", "share")
	}
	return xdgDataHome
}
