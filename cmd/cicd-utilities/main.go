package main

import (
	"context"
	"fmt"
	"os"

	"github.com/davidjspooner/cicd-utilities/pkg/command"
)

type GlobalOptions struct {
	Verbose bool `arg:"--verbose,Enable verbose output"`
}

var global GlobalOptions

var commands = command.Group{}

func main() {

	mainCommand := command.New("", "A utility for CI/CD operations",
		func(ctx context.Context, cmd command.Object, options *GlobalOptions, args []string) error {

			global = *options
			// Parse the command line arguments
			if len(args) < 1 {
				fmt.Println("Usage: cicd-utility <command> --options")
				return nil
			}
			err := commands.Execute(context.Background(), args[0], args[1:])
			return err

		}, &global)

	err := mainCommand.Execute(context.Background(), os.Args[1:])
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}
