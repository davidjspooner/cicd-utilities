package command

import (
	"fmt"
	"strings"
)

// CommandGroup represents a collection of commands.
type CommandGroup struct {
	commands []Command // List of commands in the group.
}

func (g *CommandGroup) MustAdd(cmds ...any) {
	err := g.Add(cmds...)
	if err != nil {
		panic(err)
	}
}

// Add adds one or more commands to the group.
func (g *CommandGroup) Add(cmds ...any) error {
	for _, cmd := range cmds {
		if cmd == nil {
			continue
		}
		switch cmd := cmd.(type) {
		case Command:
			_, ok := cmd.(getCommonImpl)
			if !ok {
				return fmt.Errorf("command %s is not of type commandImpl. This is a bug, contact the developers", cmd.Name())
			}
			g.commands = append(g.commands, cmd)
		case []Command:
			for _, c := range cmd {
				_, ok := c.(getCommonImpl)
				if !ok {
					return fmt.Errorf("command %s is not of type commandImpl. This is a bug, contact the developers", c.Name())
				}
				g.commands = append(g.commands, c)
			}
			g.commands = append(g.commands, cmd...)
		default:
			return fmt.Errorf("invalid command type: %T", cmd)
		}
	}
	return nil
}

// findBestCommandMatch finds the best matching command for the given target string.
func (g CommandGroup) findBestCommandMatch(target string, path []string) Match {
	best := Match{Path: nil, Score: 1000}

	for _, cmd := range g.commands {
		currentPath := append(path, cmd.Name())
		score := levenshtein(target, cmd.Name())
		if score < best.Score {
			best = Match{Path: currentPath, Score: score}
		}
		sub := cmd.SubCommands()
		subBest := sub.findBestCommandMatch(target, currentPath)
		if subBest.Score < best.Score {
			best = subBest
		}
	}

	return best
}

func (g CommandGroup) findCommand(cmdName string) Command {
	for _, cmd := range g.commands {
		get, ok := cmd.(getCommonImpl)
		if !ok {
			panic("command is not of type getCommon")
		}
		if get.common().logicalGroup {
			subgroup := cmd.SubCommands()
			if subgroup != nil {
				subCmd := subgroup.findCommand(cmdName)
				if subCmd != nil {
					return subCmd
				}
			}
		}
		if cmd.Name() == cmdName {
			return cmd
		}
	}
	return nil
}

// findCommandOrBestMatch finds a command by its name.
func (g CommandGroup) findCommandOrBestMatch(cmdName string) (Command, error) {
	cmd := g.findCommand(cmdName)
	if cmd != nil {
		return cmd, nil
	}

	match := g.findBestCommandMatch(cmdName, nil)
	if match.Score == 1000 {
		return nil, fmt.Errorf("command %s not found. This is a bug", cmdName)
	}
	return nil, fmt.Errorf("command %s not found, did you mean %s?", cmdName, strings.Join(match.Path, " "))
}

// Count returns the number of commands in the group.
func (g CommandGroup) Count() int {
	return len(g.commands)
}
