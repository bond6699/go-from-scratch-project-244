package cli

import (
	"context"
	"errors"
	"fmt"

	"code/internal/parser"

	"github.com/urfave/cli/v3"
)

var errPathArgsCount = errors.New("invalid argument count")

func getPaths(cmd *cli.Command) (string, string, error) {
	filepath1 := cmd.StringArg("filepath1")
	filepath2 := cmd.StringArg("filepath2")

	if filepath1 == "" || filepath2 == "" {
		return "", "", errPathArgsCount
	}

	return filepath1, filepath2, nil
}

func actionLogic(ctx context.Context, cmd *cli.Command) error {
	fp1, fp2, err := getPaths(cmd)
	if err != nil {
		_ = cli.ShowRootCommandHelp(cmd)
		return err
	}

	model1, err := parser.Parse(fp1)
	if err != nil {
		_ = cli.ShowRootCommandHelp(cmd)
		return err
	}

	model2, err := parser.Parse(fp2)
	if err != nil {
		_ = cli.ShowRootCommandHelp(cmd)
		return err
	}

	data, _ := parser.ToJSON(model1)
	fmt.Println(string(data))
	fmt.Println()

	data, _ = parser.ToJSON(model2)
	fmt.Println(string(data))

	return nil
}

func NewApp() *cli.Command {
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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "format", Aliases: []string{"f"},
				Usage: "output format (default: \"stylish\")",
			},
		},
		Action: actionLogic,
	}

	return app
}
