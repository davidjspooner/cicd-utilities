package main

import "fmt"

func helpCommand(args []string) error {
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
		helpCommand,
	)
}
