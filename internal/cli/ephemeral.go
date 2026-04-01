package cli

import (
	"fmt"
	"os"

	gocli "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"

	"github.com/89luca89/distrobox/pkg/commands"
	"github.com/89luca89/distrobox/pkg/ui"
)

type EphemeralCmd struct {
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
	DryRun             bool     `cli:"dry-run,d" help:"only print the container manager command generated"`
	Nopasswd           bool     `cli:"absolutely-disable-root-password-i-am-really-positively-sure"`
	gocli.Base
}

func (c *EphemeralCmd) Run() error {
	opts := commands.EphemeralOptions{
		CreateOptions: commands.CreateOptions{
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
			GenerateEntry:           false,
			Rootful:                 c.Rootful,
			ContainerAlwaysPull:     c.Pull,
			NonInteractive:          c.Yes,
		},
		DryRun: c.DryRun,
	}

	progress := ui.NewProgress(os.Stderr)
	printer := ui.NewPrinter(os.Stderr, true)

	ephemeralCmd := commands.NewEphemeralCommand(containerManager, progress, printer)
	err := ephemeralCmd.Execute(c.Ctx, opts)
	if err != nil {
		return fmt.Errorf("ephemeral command failed: %w", err)
	}

	return nil
}
