package main

import (
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

func help(args []string) error {
	fmt.Println("Available commands:")
	longestName := 0
	for name := range Commands {
		longestName = max(longestName, len(name))
	}
	longestName += 2 // for padding
	for name, cmd := range Commands {
		fmt.Printf("  %s%*s : %s\n", name, (longestName - len(name)), " ", cmd.Description)
	}
	return nil
}

func init() {
	registerCommand(
		"help",
		"Display help information",
		help,
	)
}

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cicd-utility <command> --options")
		os.Exit(1)
	}

	command := os.Args[1]
	if cmd, exists := Commands[command]; exists {
		if err := cmd.Execute(os.Args[2:]); err != nil {
			fmt.Printf("Error executing command %s: %v\n", command, err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown command: %s\n", command)
		help(os.Args[1:])
		os.Exit(1)
	}
}
