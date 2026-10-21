package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	"code"

	"github.com/urfave/cli/v3"
)

func getArgs(cmd *cli.Command) (filepath1, filepath2, format string) {
	filepath1 = cmd.StringArg("filepath1")
	filepath2 = cmd.StringArg("filepath2")
	format = cmd.String("format")

	if format == "" {
		format = "stylish"
	}

	return filepath1, filepath2, format
}

func runDiff(ctx context.Context, cmd *cli.Command) error {
	filepath1, filepath2, format := getArgs(cmd)

	result, err := code.Diff(filepath1, filepath2, format)
	if err != nil {
		_ = cli.ShowRootCommandHelp(cmd)
		return err
	}

	fmt.Println(result)

	return nil
}

func validatePath(path string) bool {
	path = strings.TrimSpace(path)
	if !strings.HasSuffix(path, ".json") && !strings.HasSuffix(path, ".yml") && !strings.HasSuffix(path, ".yaml") {
		return false
	}

	_, err := os.Stat(path)

	return err == nil
}

func New() *cli.Command {
	app := &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		UsageText: "gendiff [global options] <filepath1> <filepath2>",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "filepath1",
				UsageText: "<filepath1>",
			},
			&cli.StringArg{
				Name:      "filepath2",
				UsageText: "<filepath2>",
			},
		},

		ArgValidator: func(_ context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 2 {
				return fmt.Errorf("expected exactly 2 file paths, got %d", cmd.Args().Len())
			}

			for _, path := range cmd.Args().Slice() {
				if !validatePath(path) {
					return fmt.Errorf("expected incorrect file path, got %s", path)
				}
			}

			return nil
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "format", Aliases: []string{"f"},
				Usage: "output format (default: \"stylish\")",
				Validator: func(format string) error {
					switch format {
					case "stylish", "plain", "json":
						return nil
					default:
						return fmt.Errorf("format '%s' is not supported", format)
					}
				},
			},
		},
		Action: runDiff,
	}

	return app
}
