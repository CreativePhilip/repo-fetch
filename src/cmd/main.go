package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"repo-fetch/src/services/initialize"
	"repo-fetch/src/services/sync"
)

func main() {
	app := &cli.App{
		Name:        "rfetch",
		Description: "Sync your whole organization and maintain order",
		Commands: []*cli.Command{
			{
				Name:    "init",
				Usage:   "Initialize stuff",
				Aliases: []string{"i"},
				Action: func(ctx *cli.Context) error {
					return initialize.Init(ctx)
				},
			},
			{
				Name:    "sync",
				Usage:   "Sync repositories",
				Aliases: []string{"s"},
				Action: func(ctx *cli.Context) error {
					return sync.Sync(ctx)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
