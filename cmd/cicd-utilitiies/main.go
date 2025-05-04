package main

import (
	"flag"
	"fmt"
	"os"
)

type Command struct {
	Name        string
	Description string
	// Function to execute the command
	Execute func(args []string) error
	// Function to validate the command arguments
}

var Commands = make(map[string]Command)

func registerCommand(name string, description string, execute func(args []string) error) Command {
	command := Command{
		Name:        name,
		Description: description,
		Execute:     execute,
	}
	if _, exists := Commands[name]; exists {
		panic(fmt.Sprintf("Command %s is already registered", name))
	}
	Commands[name] = command
	return command
}

var verbose bool

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cicd-utility <command> --options")
		os.Exit(1)
	}

	mainCommand := flag.NewFlagSet("main", flag.ExitOnError)
	verbosePtr := mainCommand.Bool("verbose", false, "Enable verbose output")
	help := mainCommand.Bool("help", false, "Display help information")
	mainCommand.Parse(os.Args[1:])
	verbose = *verbosePtr
	if *help {
		helpCommand(os.Args[1:])
		os.Exit(0)
	}

	if len(mainCommand.Args()) < 1 {
		fmt.Println("Usage: cicd-utility --verbose <command> <args_and_options>")
		helpCommand(os.Args[:])
		os.Exit(1)
	}
	command := os.Args[1]
	if cmd, exists := Commands[command]; exists {
		if err := cmd.Execute(mainCommand.Args()[1:]); err != nil {
			fmt.Printf("Error executing command %s: %v\n", command, err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown command: %s\n", command)
		helpCommand(os.Args[:])
		os.Exit(1)
	}
}
