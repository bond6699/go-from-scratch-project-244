package main

import (
	"context"
	"fmt"
	"os"
	// Parser "code/internal/parser"
	"code/internal/cli"
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

	app := cli.NewApp()

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	os.Exit(0)
}
