package command

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// step represents a single step in the execution plan.
type step interface {
	// command returns the command associated with the frame.
	command() Command
	// options returns the options associated with the frame.
	options() any
	// execute runs the command with the given context and arguments.
	run(ctx context.Context, args []string) error
	postRun(ctx context.Context, args []string) error

	parseEnvAndArgs(args []string) ([]string, error)
}

// stepImpl is a generic implementation of the frame interface.
type stepImpl[T any] struct {
	cmd  *commandImpl[T] // The command associated with the frame.
	opts T               // The options for the command.
}

var _ step = &stepImpl[any]{}

func (f *stepImpl[T]) command() Command {
	return f.cmd
}

func (f *stepImpl[T]) options() any {
	return &f.opts
}

func (f *stepImpl[T]) run(ctx context.Context, args []string) error {
	if f.cmd == nil || f.cmd.run == nil {
		return nil
	}
	return f.cmd.run(ctx, &f.opts, args)
}

func (f *stepImpl[T]) postRun(ctx context.Context, args []string) error {
	if f.cmd == nil || f.cmd.postRun == nil {
		return nil
	}
	return f.cmd.postRun(ctx, &f.opts, args)
}

func (step *stepImpl[T]) parseEnvAndArgs(args []string) ([]string, error) {
	definedArgs, err := getFlagDefinitions(&step.cmd.defaultOptions)
	if err != nil {
		return args, fmt.Errorf("failed to get defined args: %v", err)
	}
	rOpts := reflect.ValueOf(&step.opts).Elem()

	err = step.ParseEnv(definedArgs, rOpts)
	if err != nil {
		return args, fmt.Errorf("failed to parse env: %v", err)
	}
	args, err = step.parseArgs(definedArgs, args, rOpts)
	if err != nil {
		return args, fmt.Errorf("failed to parse args: %v", err)
	}
	return args, nil
}

func (step *stepImpl[T]) ParseEnv(flags []Flag, rOpts reflect.Value) error {
	for _, flag := range flags {
		for _, name := range flag.aliases {
			envValue, exists := os.LookupEnv(name)
			if exists {
				var err error
				if flag.field.Type.Kind() == reflect.Bool {
					// If the flag is a boolean, set it to true if the env variable is set
					err = setFieldValue(flag, "true", rOpts.FieldByName(flag.field.Name))
				} else {
					err = setFieldValue(flag, envValue, rOpts.FieldByName(flag.field.Name))
				}
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (step *stepImpl[T]) parseArgs(flags []Flag, args []string, rOpts reflect.Value) ([]string, error) {

	for _, flag := range flags {
		for _, name := range flag.aliases {
			for i := 0; i < len(args); i++ {
				if args[i] == "--" {
					break //dont check any more args
				}
				if args[i] == name {
					removed, err := setFieldForArg(flag, name, i, args, rOpts)
					if err != nil {
						return nil, err
					}
					if removed > 0 {
						// Remove the argument from the list
						args = append(args[:i], args[i+removed:]...)
						i-- // Adjust index after removal
					}
				}
			}
		}
	}

	return args, nil
}

func setFieldForArg(flag Flag, name string, i int, args []string, rOpts reflect.Value) (removed int, err error) {
	rField := rOpts.FieldByName(flag.field.Name)
	if !rField.IsValid() {
		return 0, fmt.Errorf("field %s not found in options", flag.field.Name)
	}
	if !rField.CanSet() {
		return 0, fmt.Errorf("field %s cannot be set", flag.field.Name)
	}
	if rField.Kind() == reflect.Bool {
		err := setFieldValue(flag, "true", rField)
		if err != nil {
			return 0, err
		}
		return 1, nil
	} else if i+1 < len(args) {
		// Set the value for the option
		err := setFieldValue(flag, args[i+1], rField)
		if err != nil {
			return 0, err
		}
		return 2, nil
	} else {
		return 0, fmt.Errorf("missing value for argument %s", name)
	}
}

func setFieldValue(flag Flag, value string, rField reflect.Value) error {
	// Set the value for the option
	switch rField.Kind() {
	case reflect.Bool:
		if value == "true" || value == "1" {
			rField.SetBool(true)
		} else if value == "false" || value == "0" {
			rField.SetBool(false)
		}
	case reflect.String:
		rField.SetString(value)
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid value for integer field %s: %s", flag.field.Name, value)
		}
		rField.SetInt(intValue)
	case reflect.Float64, reflect.Float32:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value for float field %s: %s", flag.field.Name, value)
		}
		rField.SetFloat(floatValue)
	default:
		return fmt.Errorf("unsupported field type %s for field %s", rField.Type().Name(), flag.field.Name)
	}

	return nil
}
