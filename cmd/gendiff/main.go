package main

import (
	"context"
	"log/slog"
	"os"

	"code/internal/app"
)

func main() {
	cliApp := app.New()

	err := cliApp.Run(context.Background(), os.Args)
	if err != nil {
		slog.Error("gendiff failed", "error", err)
		os.Exit(1)
	}
}
