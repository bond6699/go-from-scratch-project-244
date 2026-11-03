package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	"code"

	"github.com/urfave/cli/v3"
)

var (
	errInvalidArguments  = fmt.Errorf("expected exactly 2 file paths")
	errValidatePath      = fmt.Errorf("failed to validate path")
	errUnsupportedFormat = fmt.Errorf(
		"format is not supported; select one of the possible options (stylish, plain, json)",
	)
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

	result, err := code.GenDiff(filepath1, filepath2, format)
	if err != nil {
		return fmt.Errorf("failed to execute Gendiff: %w", err)
	}

	if result != "" {
		_, err = fmt.Fprintln(cmd.Writer, result)
		if err != nil {
			return fmt.Errorf("failed to write diff: %w", err)
		}
	}

	return nil
}

func validatePath(path string) error {
	path = strings.TrimSpace(path)

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory: %s", path)
	}

	if !strings.HasSuffix(path, ".json") && !strings.HasSuffix(path, ".yml") && !strings.HasSuffix(path, ".yaml") {
		return fmt.Errorf("unsupported format")
	}

	return nil
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
			args := cmd.Args().Slice()

			if len(args) != 2 {
				_ = cli.ShowRootCommandHelp(cmd)
				return fmt.Errorf("%w, got %d", errInvalidArguments, cmd.Args().Len())
			}

			for _, path := range args {
				if err := validatePath(path); err != nil {
					return fmt.Errorf("%w %q: %w", errValidatePath, path, err)
				}
			}

			return nil
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "format", Aliases: []string{"f"},
				Usage: "output format (default: \"stylish\")",
				Value: "stylish",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			format := cmd.String("format")

			switch format {
			case "stylish", "plain", "json":
				return ctx, nil
			default:
				return ctx, fmt.Errorf("%q %w", format, errUnsupportedFormat)
			}
		},
		Action: runDiff,
	}

	return app
}
