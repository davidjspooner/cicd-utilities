package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/davidjspooner/cicd-utilities/pkg/command"
)

type SignOptions struct {
	// The name of the file to sign
	KeyFile    string `arg:"--keyfile,PGP key file to use for signing"`
	Key        string `arg:"$PGP_PRIVATE_KEY,PGP key to use for signing"`
	Passphrase string `arg:"$PGP_PASSPHRASE,Passphrase for the PGP key"`
	Extension  string `arg:"--extension,File extension override to use for the signed file"`
}

func init() {
	// Add the sign command to the root command
	cmd := command.New("pgp-sign",
		"Sign files with PGP",
		pgpSignFiles,
		&SignOptions{
			Extension: ".sig",
		},
	)
	commands = append(commands, cmd)
}

func pgpSignFiles(ctx context.Context, cmd command.Object, option *SignOptions, args []string) error {
	files, err := globFiles(args)
	if err != nil {
		return fmt.Errorf("failed to glob files: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found")
	}
	for _, file := range files {
		if err := pgpSignFile(file, option); err != nil {
			return fmt.Errorf("failed to sign file %s: %w", file, err)
		}
	}
	return nil
}

func pgpSignFile(file string, option *SignOptions) error {

	var entityList openpgp.EntityList
	var err error

	if option.KeyFile != "" {
		keyFile, err := os.Open(option.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to open key file: %w", err)
		}
		defer keyFile.Close()
		entityList, err = openpgp.ReadArmoredKeyRing(keyFile)
		if err != nil {
			return fmt.Errorf("failed to read armored key ring: %w", err)
		}
	} else if option.Key != "" {
		// Read the key from the environment variable
		r := bytes.NewReader([]byte(option.Key))
		entityList, err = openpgp.ReadArmoredKeyRing(r)
		if err != nil {
			return fmt.Errorf("failed to read armored key ring: %w", err)
		}
	} else {
		return fmt.Errorf("no key file or key provided")
	}

	if len(entityList) == 0 {
		return fmt.Errorf("no PGP keys found")
	}
	if len(entityList) > 1 {
		return fmt.Errorf("multiple PGP keys found, please specify one")
	}
	if entityList[0].PrivateKey == nil {
		return fmt.Errorf("PGP key is not private")
	}
	if err := entityList[0].PrivateKey.Decrypt([]byte(option.Passphrase)); err != nil {
		return fmt.Errorf("failed to decrypt private key: %w", err)
	}

	extension := option.Extension
	if extension == "" {
		extension = ".sig"
	}
	signedFile, err := os.Create(file + extension)
	if err != nil {
		return fmt.Errorf("failed to create signed file: %w", err)
	}
	defer signedFile.Close()

	plainFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	signWriter, err := openpgp.Sign(signedFile, entityList[0], nil, nil)
	if err != nil {
		return fmt.Errorf("failed to sign file: %w", err)
	}
	io.Copy(signWriter, plainFile)
	signWriter.Close()
	if err := signedFile.Close(); err != nil {
		return fmt.Errorf("failed to close signed file: %w", err)
	}
	return nil
}
