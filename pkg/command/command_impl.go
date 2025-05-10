package command

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

type getCommonImpl interface {
	common() *commonImpl
	newStep() (step, error)
}

type commonImpl struct {
	aliases      []string
	help         string
	metavar      string
	subCommands  CommandGroup
	logicalGroup bool
}

type commandImpl[T any] struct {
	commonImpl
	run            ExecuteFunc[T]
	postRun        ExecuteFunc[T]
	defaultOptions T
}

var _ Command = &commandImpl[any]{}

var aliasFormat = regexp.MustCompile(`^[-a-z0-9]+$`)

func parsedAliases(aliases string) []string {
	parsedAliases := strings.Split(aliases, "|")
	for i, alias := range parsedAliases {
		if !aliasFormat.MatchString(alias) {
			panic(fmt.Sprintf("alias %s is not valid, must match %s", alias, aliasFormat.String()))
		}
		parsedAliases[i] = strings.TrimSpace(alias)
	}
	return parsedAliases
}

func NewCommand[T any](aliases, help string, execute ExecuteFunc[T], defaults *T, modfiers ...Middeware) Command {
	if defaults == nil {
		panic("defaults cannot be nil")
	}
	if aliases == "" {
		aliases = path.Base(os.Args[0])
		aliases = strings.TrimSuffix(aliases, path.Ext(aliases))
		aliases = strings.Replace(aliases, "_", "-", -1)
		aliases = strings.Trim(aliases, "-")

	}
	impl := &commandImpl[T]{
		commonImpl: commonImpl{
			aliases: parsedAliases(aliases),
		},
		run:            execute,
		defaultOptions: *defaults,
	}
	var err error
	impl.metavar, impl.help, err = extractMetaVar(help)
	if err != nil {
		panic(fmt.Sprintf("failed to extract metavar from help: %v", err))
	}
	for _, modifier := range modfiers {
		modifier(impl)
	}
	return impl
}

func (c *commandImpl[T]) Name() string {
	if len(c.aliases) > 0 {
		return c.aliases[0]
	}
	return "unamed"
}

func (c *commandImpl[T]) Help() string {
	return c.help
}

func (c *commandImpl[T]) SubCommands() *CommandGroup {
	return &c.subCommands
}

func (c *commandImpl[T]) Flags() ([]Flag, error) {
	definedArgs, err := getFlagDefinitions(&c.defaultOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get defined args: %v", err)
	}
	return definedArgs, nil
}

func (c *commandImpl[T]) newStep() (step, error) {
	step := &stepImpl[T]{
		cmd: c,
	}
	if step.cmd == nil {
		return nil, fmt.Errorf("command cannot be nil")
	}
	return step, nil
}

func (c *commandImpl[T]) With(modifiers ...Middeware) Command {
	new := *c
	for _, modifier := range modifiers {
		modifier(&new)
	}
	return &new
}

func (c *commandImpl[T]) Aliases() []string {
	return c.aliases
}

func (c *commandImpl[T]) common() *commonImpl {
	return &c.commonImpl
}
