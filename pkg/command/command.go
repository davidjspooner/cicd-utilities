package command

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
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
			if strings.HasPrefix(name, "$") {
				// Check if the environment variable is set
				envVar := strings.TrimPrefix(name, "$")
				if value, ok := os.LookupEnv(envVar); ok {
					f := rOptions.FieldByName(arg.field.Name)
					if f.Kind() == reflect.Bool {
						err := c.setField(f, "true")
						if err != nil {
							return fmt.Errorf("failed to set value for field %s: %v", arg.field.Name, err)
						}
					} else {
						err := c.setField(f, value)
						if err != nil {
							return fmt.Errorf("failed to set value for field %s: %v", arg.field.Name, err)
						}
					}
				}
				continue
			}
			for i := 0; i < len(expandedArgs); i++ {
				if expandedArgs[i] == name {
					f := rOptions.FieldByName(arg.field.Name)
					if f.Kind() == reflect.Bool {
						err := c.setField(f, "true")
						if err != nil {
							return fmt.Errorf("failed to set value for field %s: %v", arg.field.Name, err)
						}
					} else if i+1 < len(expandedArgs) {
						err := c.setField(f, expandedArgs[i+1])
						if err != nil {
							return fmt.Errorf("failed to set value for field %s: %v", arg.field.Name, err)
						}
						// Remove the next argument as it has been consumed
						expandedArgs = append(expandedArgs[:i+1], expandedArgs[i+2:]...)
					} else {
						return fmt.Errorf("missing value for argument %s", name)
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

func (c *command[T]) setField(rField reflect.Value, value string) error {
	switch rField.Kind() {
	case reflect.Bool:
		if value == "true" {
			rField.SetBool(true)
		} else if value == "false" {
			rField.SetBool(false)
		} else {
			return fmt.Errorf("invalid value for boolean field %s: %s", rField.Type().Name(), value)
		}
	case reflect.String:
		rField.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid value for integer field %s: %s", rField.Type().Name(), value)
		}
		rField.SetInt(intValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for float field %s: %s", rField.Type().Name(), value)
		}
		rField.SetFloat(floatValue)
	default:
		return fmt.Errorf("unsupported type for field %s: %s", rField.Type().Name(), rField.Kind())
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
