package archive

import "github.com/davidjspooner/cicd-utilities/pkg/command"

func Commands() []command.Command {
	commands := []command.Command{}
	cmd1 := command.NewCommand(
		"pgp-sign",
		"Sign files with PGP",
		pgpSignFiles,
		&SignOptions{
			Extension: ".sig",
		},
	)
	commands = append(commands, cmd1)
	cmd2 := command.NewCommand(
		"checksum",
		"Generate checksum(s) for file(s) using a specified algorithm",
		executeChecksum,
		&ChecksumOptions{
			Algorithm: "sha256",
		},
	)
	commands = append(commands, cmd2)
	cmd3 := command.NewCommand(
		"compress",
		"Compress files or directories into zip or tar.gz formats",
		compressCommand,
		&CompressOptions{
			Format:  "tar.gz",
			Replace: false,
		},
	)
	commands = append(commands, cmd3)
	return commands
}
