package command

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Object interface {
	Name() string
	Description() string
	Options() ([]Option, error)
	Execute(ctx context.Context, args []string) error
}

type Group []Object

func (g Group) Execute(ctx context.Context, name string, args []string) error {
	for _, subCommand := range g {
		if subCommand.Name() == name {
			return subCommand.Execute(ctx, args)
		}
	}
	//TODO: add "did you mean" functionality
	return fmt.Errorf("command %s not found", name)
}

type Function[T any] func(ctx context.Context, cmd Object, options *T, args []string) error

type command[T any] struct {
	name        string
	description string
	execute     Function[T]
	defaults    T
}

func New[T any](name, description string, execute Function[T], defaults *T) *command[T] {
	if defaults == nil {
		panic("defaults cannot be nil")
	}
	if name == "" {
		name = os.Args[0]
	}
	return &command[T]{
		name:        name,
		description: description,
		execute:     execute,
		defaults:    *defaults,
	}
}

func (c *command[T]) Name() string {
	return c.name
}

func (c *command[T]) Description() string {
	return c.description
}

func (c *command[T]) Execute(ctx context.Context, args []string) error {

	definedArgs, err := c.Options()
	if err != nil {
		return fmt.Errorf("failed to get defined args: %v", err)
	}

	expandedArgs, err := expandInputArgs(args)
	if err != nil {
		return fmt.Errorf("failed to expand input args: %v", err)
	}

	var options = c.defaults
	rOptions := reflect.ValueOf(&options).Elem()

	for _, arg := range definedArgs {
		for _, name := range arg.Names {
			for i := 0; i < len(expandedArgs); i++ {
				if expandedArgs[i] == name {
					switch arg.field.Type.Kind() {
					case reflect.Bool:
						// Set the boolean field to true
						f := rOptions.FieldByName(arg.field.Name)
						f.SetBool(true)
						expandedArgs = append(expandedArgs[:i], expandedArgs[i+1:]...)
					case reflect.String:
						// Set the string field to the next argument
						if i+1 < len(expandedArgs) {
							f := rOptions.FieldByName(arg.field.Name)
							f.SetString(expandedArgs[i+1])
							// Remove the next argument from the list
							expandedArgs = append(expandedArgs[:i], expandedArgs[i+2:]...)
						} else {
							return fmt.Errorf("missing value for argument %s", name)
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						// Set the int field to the next argument
						if i+1 < len(expandedArgs) {
							value, err := strconv.Atoi(expandedArgs[i+1])
							if err != nil {
								return fmt.Errorf("invalid value for argument %s: %v", name, err)
							}
							f := rOptions.FieldByName(arg.field.Name)
							f.SetInt(int64(value))
							// Remove the next argument from the list
							expandedArgs = append(expandedArgs[:i], expandedArgs[i+2:]...)
						} else {
							return fmt.Errorf("missing value for argument %s", name)
						}
					case reflect.Float32, reflect.Float64:
						// Set the float field to the next argument
						if i+1 < len(expandedArgs) {
							value, err := strconv.ParseFloat(expandedArgs[i+1], 64)
							if err != nil {
								return fmt.Errorf("invalid value for argument %s: %v", name, err)
							}
							f := rOptions.FieldByName(arg.field.Name)
							f.SetFloat(value)
							// Remove the next argument from the list
							expandedArgs = append(expandedArgs[:i], expandedArgs[i+2:]...)
						} else {
							return fmt.Errorf("missing value for argument %s", name)
						}
					default:
						return fmt.Errorf("unsupported type for argument %s: %s", name, arg.field.Type.Kind())
					}
				}
			}
		}
	}

	err = c.execute(ctx, c, &options, expandedArgs)
	if err != nil {
		return err
	}
	return nil
}

func (c *command[T]) Options() ([]Option, error) {
	definedArgs, err := getDefinedOptions(&c.defaults)
	if err != nil {
		return nil, fmt.Errorf("failed to get defined args: %v", err)
	}
	return definedArgs, nil
}
