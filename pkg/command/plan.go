package command

import (
	"context"
	"fmt"
	"strings"
)

type plan struct {
	steps        []step
	currentIndex int
	unparsedArgs []string
}

func (plan *plan) checkForUnparsedFlags() error {
	for _, arg := range plan.unparsedArgs {
		if arg == "--" {
			break
		}
		if arg == "-" { //some commands use "-" to represent stdin
			continue
		}
		if strings.HasPrefix(arg, "-") {
			//find the best match for the flag
			bestMatch, alternatives := plan.findBestFlagMatch(arg)
			if bestMatch == "" {
				return fmt.Errorf("unknown flag %s, did you mean one of %s?", arg, strings.Join(alternatives, ", "))
			} else {
				return fmt.Errorf("unknown flag %s, did you mean %s?", arg, bestMatch)
			}
		}
	}
	return nil
}

// run runs all frames in the execution plan sequentially.
func (plan *plan) run(ctx context.Context) error {
	var err error
	for plan.currentIndex = 0; plan.currentIndex < len(plan.steps); plan.currentIndex++ {
		frame := plan.steps[plan.currentIndex]
		if plan.currentIndex == len(plan.steps)-1 {
			err = frame.run(ctx, plan.unparsedArgs)
		} else {
			err = frame.run(ctx, nil)
		}
		if err != nil {
			return err
		}
	}
	for i := len(plan.steps) - 1; i >= 0; i-- {
		frame := plan.steps[i]
		if plan.currentIndex == len(plan.steps)-1 {
			err = frame.postRun(ctx, plan.unparsedArgs)
		} else {
			err = frame.postRun(ctx, nil)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// addStep adds a new frame to the execution plan.
func (plan *plan) addStep(cmd Command) error {
	common, ok := cmd.(getCommonImpl)
	if !ok {
		return fmt.Errorf("command is not of type getCommon. This is a bug, contact the developers")
	}

	step, err := common.newStep()
	if err != nil {
		return fmt.Errorf("failed to parse args: %v", err)
	}

	args, err := step.parseEnvAndArgs(plan.unparsedArgs)
	if err != nil {
		return fmt.Errorf("failed to parse env: %v", err)
	}
	plan.steps = append(plan.steps, step)
	plan.unparsedArgs = args
	return nil
}

func (p *plan) findBestFlagMatch(arg string) (string, []string) {
	bestMatch := ""
	alternatives := []string{}
	lowestDistance := -1

	for _, step := range p.steps {
		flags, _ := step.command().Flags()
		for _, flag := range flags {
			for _, alias := range flag.Aliases() {
				if strings.HasPrefix(alias, "$") {
					continue
				}
				distance := levenshtein(arg, alias)

				// Track best match
				if lowestDistance == -1 || distance < lowestDistance {
					bestMatch = alias
					lowestDistance = distance
				}

				// Collect similar alternatives (within a threshold)
				if distance <= 2 {
					alternatives = append(alternatives, alias)
				}
			}
		}
	}

	return bestMatch, alternatives
}

// buildPlan constructs an execution plan from the given arguments.
func buildPlan(root Command, args []string) (*plan, error) {

	_, ok := root.(getCommonImpl)
	if !ok {
		return nil, fmt.Errorf("root command is not of type commandImpl. This is a bug, contact the developers")
	}

	plan := &plan{
		steps:        make([]step, 0),
		currentIndex: -1,
		unparsedArgs: args,
	}
	var err error
	plan.unparsedArgs, err = normalizeArgs(args)
	if err != nil {
		return nil, fmt.Errorf("failed to expand args: %v", err)
	}

	// Recursively iterate over the commands starting at TopLevel and find the exact match.
	curentFrame := root

	err = plan.addStep(curentFrame)
	if err != nil {
		return nil, err
	}

	var cmdName string
	for len(plan.unparsedArgs) > 0 {
		subCommands := curentFrame.SubCommands()
		if subCommands.Count() == 0 {
			break
		}
		cmdName, plan.unparsedArgs = extractFirstNonFlagArg(plan.unparsedArgs)
		if cmdName == "" {
			break
		}
		prefix := strings.Builder{}
		for _, cmd := range plan.steps {
			prefix.WriteString(cmd.command().Name())
			prefix.WriteString(" ")
		}
		cmdDef, err := subCommands.findCommandOrBestMatch(prefix.String(), cmdName)
		if err != nil {
			return nil, err
		}

		err = plan.addStep(cmdDef)
		if err != nil {
			return nil, err
		}
		curentFrame = cmdDef
	}
	return plan, nil
}
