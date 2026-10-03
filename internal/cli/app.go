package cli

import (
	"github.com/urfave/cli/v3"
)

func NewApp() *cli.Command {
	app := &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		UsageText: "gendiff [global options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "format", Aliases: []string{"f"},
				Usage: "output format (default: \"stylish\")",
			},
		},
	}

	return app
}
