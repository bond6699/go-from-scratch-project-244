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

	fmt.Println()

	if err != nil {
		slog.Error("gendiff failed", "error", err)
		os.Exit(1)
	}
}
