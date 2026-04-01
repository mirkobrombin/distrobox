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

type CreateCmd struct {
	Image              string   `cli:"image,i" help:"image to use for the container" default:"registry.fedoraproject.org/fedora-toolbox:latest"`
	Name               string   `cli:"name,n" help:"name for the distrobox"`
	Hostname           string   `cli:"hostname" help:"hostname for the distrobox"`
	Pull               bool     `cli:"pull,p" help:"pull the image even if it exists locally"`
	Yes                bool     `cli:"yes,Y" help:"non-interactive, pull images without asking"`
	Rootful            bool     `cli:"root,r" help:"launch with root privileges"`
	Clone              string   `cli:"clone,c" help:"name of the distrobox container to use as base for a new container"`
	Home               string   `cli:"home,H" help:"select a custom HOME directory for the container"`
	Volume             []string `cli:"volume" help:"additional volumes to add to the container"`
	AdditionalFlags    []string `cli:"additional-flags,a" help:"additional flags to pass to the container manager command"`
	AdditionalPackages []string `cli:"additional-packages,ap" help:"additional packages to install during initial container setup"`
	InitHooks          string   `cli:"init-hooks" help:"additional commands to execute at the end of container initialization"`
	PreInitHooks       string   `cli:"pre-init-hooks" help:"additional commands to execute at the start of container initialization"`
	Init               bool     `cli:"init,I" help:"use init system inside the container"`
	Nvidia             bool     `cli:"nvidia" help:"try to integrate host nVidia drivers in the guest"`
	Platform           string   `cli:"platform" help:"specify which platform to use, eg: linux/arm64"`
	UnshareDevsys      bool     `cli:"unshare-devsys" help:"do not share host devices and sysfs dirs from host"`
	UnshareGroups      bool     `cli:"unshare-groups" help:"do not forward user additional groups into the container"`
	UnshareIpc         bool     `cli:"unshare-ipc" help:"do not share ipc namespace with host"`
	UnshareNetns       bool     `cli:"unshare-netns" help:"do not share the net namespace with host"`
	UnshareProcess     bool     `cli:"unshare-process" help:"do not share process namespace with host"`
	UnshareAll         bool     `cli:"unshare-all" help:"activate all the unshare flags below"`
	NoEntry            bool     `cli:"no-entry" help:"do not generate a container entry in the application list"`
	DryRun             bool     `cli:"dry-run,d" help:"only print the container manager command generated"`
	Nopasswd           bool     `cli:"absolutely-disable-root-password-i-am-really-positively-sure"`
	Compatibility      bool     `cli:"compatibility,C" help:"show compatibility information and exit"`
	gocli.Base
}

func (c *CreateCmd) Run() error {
	if c.Compatibility {
		return showCompatibility()
	}

	opts := commands.CreateOptions{
		ContainerImage:          c.Image,
		ContainerName:           c.Name,
		ContainerHostname:       c.Hostname,
		ContainerClone:          c.Clone,
		UnshareNetNs:            c.UnshareNetns || c.UnshareAll,
		UnshareDevsys:           c.UnshareDevsys || c.UnshareAll,
		UnshareGroups:           c.UnshareGroups || c.UnshareAll || c.Init,
		UnshareIpc:              c.UnshareIpc || c.UnshareAll,
		UnshareProcess:          c.UnshareProcess || c.UnshareAll || c.Init,
		AdditionalFlags:         c.AdditionalFlags,
		AdditionalVolumes:       c.Volume,
		AdditionalPackages:      c.AdditionalPackages,
		Nopasswd:                c.Nopasswd,
		ContainerUserCustomHome: c.Home,
		Init:                    c.Init,
		Nvidia:                  c.Nvidia,
		ContainerInitHook:       c.InitHooks,
		ContainerPreInitHook:    c.PreInitHooks,
		ContainerPlatform:       c.Platform,
		DryRun:                  c.DryRun,
		GenerateEntry:           !c.NoEntry,
		Rootful:                 c.Rootful,
		ContainerAlwaysPull:     c.Pull,
		NonInteractive:          c.Yes,
	}

	progress := ui.NewProgress(os.Stderr)
	prompter := ui.NewPrompter(*bufio.NewReader(os.Stdin), os.Stdout)

	createCmd := commands.NewCreateCommand(containerManager, progress, prompter)
	err := createCmd.Execute(c.Ctx, opts)

	var containerAlreadyExistsErr *commands.ContainerAlreadyExistsError
	if errors.As(err, &containerAlreadyExistsErr) {
		printContainerAlreadyExists(progress, containerAlreadyExistsErr.ContainerName, opts.Rootful)
	}

	if errors.Is(err, commands.ErrImagePullAbortedByUser) {
		progress.Finalize("next time, pull the image first")
		return nil
	}

	if err != nil {
		return fmt.Errorf("create command failed: %w", err)
	}

	if !opts.DryRun {
		printCreateCompleted(progress, opts.ContainerName, opts.Rootful)
	}

	return nil
}

func showCompatibility() error {
	// TODO: fetch compatibility
	return nil
}

func printCreateCompleted(progress *ui.Progress, containerName string, rootful bool) {
	rootFlag := ""
	if rootful {
		rootFlag = "--root "
	}
	msg := "Distrobox '%s' successfully created.\nTo enter, run:\n\ndistrobox enter %s%s\n\n"
	progress.Finalize(msg, containerName, rootFlag, containerName)
}

func printContainerAlreadyExists(progress *ui.Progress, containerName string, rootful bool) {
	rootFlag := ""
	if rootful {
		rootFlag = "--root "
	}
	msg := "Distrobox named '%s' already exists.\nTo enter, run:\n\ndistrobox enter %s%s\n\n"
	progress.Finalize(msg, containerName, rootFlag, containerName)
}
