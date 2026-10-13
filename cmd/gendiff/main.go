package main

import (
	"context"
	"fmt"
	"os"

	// Parser "code/internal/parser"
	"code/internal/app"
)

func main() {
	/*
	     var TestJSON string = `
	   {
	     "name": "Bohdan",
	     "age": 25,
	     "active": true,
	     "items": [1, 2, 3],
	     "settings": {
	       "theme": "dark",
	       "settings-theme": {
	         "color": "red"
	       }
	     }
	   }
	   `

	   	data := Parser.Parser()
	   	fmt.Println(data)
	*/

	app := app.New()

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Println()
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	os.Exit(0)
}
