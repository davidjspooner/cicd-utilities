package command

import (
	"context"
	"fmt"
	"os"
	"path"
)

type VersionOptions struct {
	Short bool `flag:"--short,Print only the version number"`
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

func VersionCommand() Command {
	toolname := path.Base(os.Args[0])
	return NewCommand("version", fmt.Sprintf("Print the version of %s", toolname),
		func(ctx context.Context, options *VersionOptions, args []string) error {
			if options.Short {
				fmt.Printf("%s\n", BuildName)
				return nil
			}
			fmt.Printf("version: %s\n", BuildName)
			fmt.Printf("date: %s\n", BuildDate)
			fmt.Printf("by: %s\n", BuildBy)
			return nil
		}, &VersionOptions{})
}
