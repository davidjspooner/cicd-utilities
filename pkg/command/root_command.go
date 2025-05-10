package command

import (
	"context"
	"fmt"
	"os"
	"reflect"
)

type NoopOptions struct{}

type contextKey struct{}

// extractPlan retrieves the execution plan from the context.
func extractPlan(ctx context.Context) (*plan, bool) {
	stack, ok := ctx.Value(contextKey{}).(*plan)
	return stack, ok
}

// RootCommand is the top-level command for the application.
var RootCommand Command = NewCommand(
	os.Args[0],
	"undefined description of overall command",
	func(ctx context.Context, options *NoopOptions, args []string) error {
		//noop
		return nil
	},
	&NoopOptions{},
)

// Run executes the top-level command with the given arguments.
func Run(ctx context.Context, args []string) error {
	_, ok := extractPlan(ctx)
	if ok {
		return fmt.Errorf("command already executing")
	}

	args, showHelp := extractHelpTriggers(args)

	plan, err := buildPlan(RootCommand, args)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, contextKey{}, plan)

	if showHelp {
		//show help and exit
		plan.renderHelpText(ctx)
		os.Exit(0)
	}

	//if we parsed the command check there are no unparsed args
	err = plan.checkForUnparsedFlags()
	if err != nil {
		return err
	}
	return plan.run(ctx)
}

// FindOptionStruct finds the option struct of the given type in the execution plan.
func FindOptionStruct[T any](ctx context.Context) (*T, error) {
	plan, ok := extractPlan(ctx)
	if !ok {
		return nil, fmt.Errorf("no execution plan found in context")
	}
	for i := plan.currentIndex; i >= 0; i-- {
		frame, ok := plan.steps[i].(*stepImpl[T])
		if ok {
			return &frame.opts, nil
		}
	}
	return nil, fmt.Errorf("no options found in execution plan")
}

// FindOptionField finds the option field of the given type and name in the execution plan.
func FindOptionField[T any](ctx context.Context, fieldName string) (T, error) {
	plan, ok := extractPlan(ctx)
	if !ok {
		return *new(T), fmt.Errorf("no execution plan found in context")
	}
	for i := plan.currentIndex; i >= 0; i-- {
		options := plan.steps[i].options()
		if options == nil {
			continue
		}
		rOptions := reflect.ValueOf(options)
		ft, found := rOptions.Type().FieldByName(fieldName)
		if !found {
			continue
		}
		if !ft.IsExported() {
			return *new(T), fmt.Errorf("field %s in %s is not exported", fieldName, rOptions.Type().Name())
		}
		f := rOptions.FieldByName(fieldName)
		if f.IsValid() {
			if !f.CanInterface() {
				return *new(T), fmt.Errorf("field %s in %s is not exported", fieldName, rOptions.Type().Name())
			}
			cast, ok := f.Interface().(T)
			if ok {
				return cast, nil
			}
			return *new(T), fmt.Errorf("field %s in %s is not of type %T", fieldName, rOptions.Type().Name(), *new(T))
		}
	}
	return *new(T), fmt.Errorf("no field %s found in execution plan", fieldName)
}
