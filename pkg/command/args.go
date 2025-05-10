package command

import (
	"slices"
	"strings"
)

func normalizeArgs(argsIn []string) ([]string, error) {

	argsOut := make([]string, 0, len(argsIn))
	stopParsing := false
	for _, arg := range argsIn {
		if stopParsing {
			argsOut = append(argsOut, arg)
			continue
		}

		if arg == "--" {
			argsOut = append(argsOut, arg)
			stopParsing = true
			continue
		}

		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg, "=", 2)
			argsOut = append(argsOut, parts[0])
			if len(parts) == 2 {
				argsOut = append(argsOut, parts[1])
			}
		} else if strings.HasPrefix(arg, "-") {
			if len(arg) == 1 {
				argsOut = append(argsOut, arg)
				continue
			}
			parts := strings.SplitN(arg, "=", 2)
			for _, letter := range parts[0][1:] {
				argsOut = append(argsOut, "-"+string(letter))
			}
			if len(parts) == 2 {
				argsOut = append(argsOut, parts[1])
			}
		} else {
			argsOut = append(argsOut, arg)
		}
	}
	return argsOut, nil
}

func extractFirstNonFlagArg(args []string) (string, []string) {
	for i, arg := range args {
		if arg == "--" {
			return "", args
		}
		if !strings.HasPrefix(arg, "-") {
			return arg, append(args[:i], args[i+1:]...)
		}
	}
	return "", args
}

var helpTriggers []string = []string{"--help", "-h", "help"}

// SetHelpTriggers enables custom help commands.
func SetHelpTriggers(args ...string) {
	helpTriggers = args
}

func extractHelpTriggers(args []string) ([]string, bool) {
	foundHelp := false
	for i := 0; i < len(args); i++ {
		if slices.Contains(helpTriggers, args[i]) {
			foundHelp = true
			args = append(args[:i], args[i+1:]...)
			break
		}
		if args[i] == "--" {
			break
		}
	}
	return args, foundHelp
}
