package main

import (
	"context"
	"fmt"
	"os"

	"code/internal/app"
)

func main() {
	cliApp := app.New()

	if err := cliApp.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
