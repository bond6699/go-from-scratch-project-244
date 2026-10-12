package cli

import (
	"context"
	"fmt"

	"code/internal/analyzer"
	"code/internal/formatter"
	"code/internal/parser"

	"github.com/urfave/cli/v3"
)

var formatsMap map[string]struct{} = map[string]struct{}{
	"stylish": {},
	"plain":   {},
	"json":    {},
}

func getArgs(cmd *cli.Command) (string, string, string, error) {
	filepath1 := cmd.StringArg("filepath1")
	filepath2 := cmd.StringArg("filepath2")
	format := cmd.String("format")

	if filepath1 == "" || filepath2 == "" {
		return "", "", "", fmt.Errorf("invalid argument count")
	}

	if format == "" {
		format = "stylish"
	}

	if _, ok := formatsMap[format]; !ok {
		return filepath1, filepath2, format, fmt.Errorf("invalid format")
	}

	return filepath1, filepath2, format, nil
}

func actionLogic(ctx context.Context, cmd *cli.Command) error {
	fp1, fp2, format, err := getArgs(cmd)
	if err != nil {
		_ = cli.ShowRootCommandHelp(cmd)
		return err
	}

	data1, err := parser.Parse(fp1)
	if err != nil {
		_ = cli.ShowRootCommandHelp(cmd)
		return err
	}
	data1 = parser.Normalize(data1)

	data2, err := parser.Parse(fp2)
	if err != nil {
		_ = cli.ShowRootCommandHelp(cmd)
		return err
	}
	data2 = parser.Normalize(data2)

	diff := analyzer.Analyze(data1, data2)

	result, err := formatter.Format(format, diff)
	if err != nil {
		fmt.Print(err)
		return err
	}

	fmt.Println(result)

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
