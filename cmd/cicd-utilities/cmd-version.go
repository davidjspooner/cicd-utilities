package main

import (
	"context"
	"fmt"

	"github.com/davidjspooner/cicd-utilities/pkg/command"
)

type VersionOptions struct {
	Short bool `arg:"--short,Print only the version number"`
}

const (
	// Version is the version of the cicd-utilities
	BuildName = "snapshot"
	// VersionDate is the date of the version of the cicd-utilities
	BuildDate = "2023-10-01"
	// VersionCommit is the commit of the version of the cicd-utilities
	BuildBy = "tbd"
	// VersionBuild is the build of the version of the cicd-utilities
)

func init() {
	cmd := command.New(
		"version",
		"Print the version of the cicd-utilities",
		versionCommand,
		&VersionOptions{},
	)
	commands = append(commands, cmd)
}

func versionCommand(ctx context.Context, cmd command.Object, options *VersionOptions, args []string) error {
	if options.Short {
		fmt.Printf("%s\n", BuildName)
		return nil
	}
	fmt.Printf("version: %s\n", BuildName)
	fmt.Printf("date: %s\n", BuildDate)
	fmt.Printf("by: %s\n", BuildBy)
	return nil
}

//
