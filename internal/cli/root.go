package cli

import (
	"fmt"
	"os"
	"os/exec"

	gocli "github.com/mirkobrombin/go-cli-builder/v2/pkg/cli"

	"github.com/89luca89/distrobox/pkg/containermanager"
	"github.com/89luca89/distrobox/pkg/containermanager/providers"
	"github.com/mirkobrombin/go-foundation/pkg/adapters"
)

// Root is the root CLI struct.
type Root struct {
	ContainerManager string `cli:"container-manager" env:"DBX_CONTAINER_MANAGER" default:""`
	SudoCommand      string `cli:"sudo-command" env:"DBX_SUDO_COMMAND" default:"sudo"`
	Root             bool   `cli:"root,r"`
	Verbose          bool   `cli:"verbose,v" env:"DBX_VERBOSE"`

	List          ListCmd          `cmd:"list" help:"List distroboxes" aliases:"ls"`
	GenerateEntry GenerateEntryCmd `cmd:"generate-entry" help:"Generate or delete distrobox entries"`
	Create        CreateCmd        `cmd:"create" help:"Create a new distrobox container"`
	Enter         EnterCmd         `cmd:"enter" help:"Enter a distrobox"`
	Assemble      AssembleCmd      `cmd:"assemble" help:"Create or remove distroboxes from a manifest file"`
	Rm            RmCmd            `cmd:"rm" help:"Remove distroboxes"`
	Stop          StopCmd          `cmd:"stop" help:"Stop running distrobox containers"`
	Ephemeral     EphemeralCmd     `cmd:"ephemeral" help:"Create a temporary distrobox that is removed on exit"`
	Upgrade       UpgradeCmd       `cmd:"upgrade" help:"Upgrade packages inside distrobox containers"`

	gocli.Base
}

func (r *Root) Before() error {
	if r.Root {
		if err := validateSudo(); err != nil {
			return fmt.Errorf("cannot run in root mode: %w", err)
		}
	}

	registry := buildContainerManagerRegistry(r.Root, r.SudoCommand, r.Verbose)

	cmType := r.ContainerManager
	if cmType != "" && cmType != "podman-static" {
		cm, ok := registry.Get(cmType)
		if !ok {
			return fmt.Errorf("unsupported container manager: %s", cmType)
		}
		containerManager = cm
	} else {
		containerManager = registry.Default()
	}

	return nil
}

func buildContainerManagerRegistry(root bool, sudoCommand string, verbose bool) *adapters.Registry[containermanager.ContainerManager] {
	registry := adapters.NewRegistry[containermanager.ContainerManager]()
	registry.Register("docker", providers.NewDocker(root, sudoCommand, verbose))
	registry.Register("podman", providers.NewPodman(root, sudoCommand, verbose))
	registry.SetDefault("podman")
	return registry
}

// NewRootCommand returns the go-cli-builder App for main.go compatibility.
func NewRootCommand() *gocli.App {
	app, _ := gocli.New(&Root{})
	app.SetName("distrobox")
	return app
}

func validateSudo() error {
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to validate sudo: %w", err)
	}

	return nil
}
