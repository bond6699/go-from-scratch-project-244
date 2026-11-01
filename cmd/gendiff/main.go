package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"code/internal/app"
)

func main() {
	cliApp := app.New()

	err := cliApp.Run(context.Background(), os.Args)
	if err != nil {
		logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

		fmt.Println()
		logger.Error("gendiff failed", "error", err)

		os.Exit(1)
	}
}
