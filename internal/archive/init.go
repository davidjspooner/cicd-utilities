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

	checksumCmd := command.NewCommand(
		"checksum",
		"Generate checksum(s) for file(s) using a specified algorithm",
		executeChecksum,
		&ChecksumOptions{
			Algorithm: "sha256",
		},
	)
	compressCmd := command.NewCommand(
		"compress",
		"Compress files or directories into zip or tar.gz formats",
		compressCommand,
		&CompressOptions{
			Format:  "tar.gz",
			Replace: false,
		},
	)
	archiveCommand.SubCommands().MustAdd(checksumCmd, compressCmd)
	return []command.Command{archiveCommand}
}
