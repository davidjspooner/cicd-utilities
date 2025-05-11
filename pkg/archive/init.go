package archive

import (
	"github.com/davidjspooner/cicd-utilities/pkg/command"
)

func Commands() []command.Command {

	archiveCommand := command.NewCommand(
		"archive",
		"Archive commands",
		nil,
		&command.NoopOptions{},
	)

	cmd1 := command.NewCommand(
		"pgp-sign",
		"Sign files with PGP",
		pgpSignFiles,
		&SignOptions{
			Extension: ".sig",
		},
	)
	cmd2 := command.NewCommand(
		"checksum",
		"Generate checksum(s) for file(s) using a specified algorithm",
		executeChecksum,
		&ChecksumOptions{
			Algorithm: "sha256",
		},
	)
	cmd3 := command.NewCommand(
		"compress",
		"Compress files or directories into zip or tar.gz formats",
		compressCommand,
		&CompressOptions{
			Format:  "tar.gz",
			Replace: false,
		},
	)
	archiveCommand.SubCommands().MustAdd(cmd1, cmd2, cmd3)
	return []command.Command{archiveCommand}
}
