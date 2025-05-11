package man

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/davidjspooner/cicd-utilities/pkg/command"
)

func Commands() []command.Command {
	arg0 := path.Base(os.Args[0])
	manCmd := command.NewCommand(
		"man",
		fmt.Sprintf("Manual for %s", arg0),
		manPage,
		&command.NoopOptions{},
	)
	return []command.Command{manCmd}
}

func manPage(ctx context.Context, options *command.NoopOptions, args []string) error {
	print("#TODO")
	return nil
}
