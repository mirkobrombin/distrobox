package main

import (
	"log"

	"github.com/89luca89/distrobox/internal/cli"
	"github.com/89luca89/distrobox/internal/config"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	cmd := cli.NewRootCommand()
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
