package command

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Flag struct {
	aliases      []string
	help         string
	metaVar      string
	defaultValue string
	field        reflect.StructField
}

func (flag *Flag) Aliases() []string {
	return flag.aliases
}
func (flag *Flag) Help() string {
	return flag.help
}
func (flag *Flag) MetaVar() string {
	return flag.metaVar
}
func (flag *Flag) DefaultValue() string {
	return flag.defaultValue
}
func (flag *Flag) Kind() reflect.Kind {
	return flag.field.Type.Kind()
}

func (o Flag) String() string {
	help := o.help
	if help == "" {
		help = "No help available"
	}
	return fmt.Sprintf("%s: %s (default: %s)", strings.Join(o.aliases, ", "), help, o.defaultValue)
}

var metaVarFormat = regexp.MustCompile(`<[a-z0-9\|]+>`)

var spaceWithLineBreak = regexp.MustCompile("\\s*\n\\s*")

func extractMetaVar(help string) (string, string, error) {
	help = strings.TrimSpace(help)

	//for each position of the string, find the pattern spaceWithLineBreak
	help = spaceWithLineBreak.ReplaceAllString(help, " ")
	help = strings.TrimSpace(help)
	help = strings.ReplaceAll(help, "\\n", "\n")

	if help == "" {
		return "", help, nil
	}
	var metaVar string
	for n, example := range metaVarFormat.FindAllString(help, -1) {
		if n == 0 {
			metaVar = example
		} else if metaVar != example {
			// If we find multiple meta vars, return an error as this is not expected
			return "", help, fmt.Errorf("different <metavar> found in help text: %s", help)
		}
		metaVar = example
	}
	return metaVar, help, nil
}

func getFlagDefinitions[T any](defaults *T) ([]Flag, error) {
	if defaults == nil {
		return nil, fmt.Errorf("defaults cannot be nil")
	}

	rDefaults := reflect.ValueOf(defaults).Elem()
	t := rDefaults.Type()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", t.Kind())
	}

	var args []Flag
	seenaliases := map[string]bool{}

	for i := range t.NumField() {
		arg := Flag{}
		arg.field = t.Field(i)
		if !arg.field.IsExported() {
			continue
		}

		tag := arg.field.Tag.Get("flag")
		if strings.TrimSpace(tag) == "" {
			continue // no tag, skip
		}

		tagParts := strings.SplitN(tag, ",", 2)
		if len(tagParts) != 2 {
			return nil, fmt.Errorf("invalid tag format for field %s: %q. Expected something like '--opt|-o|$ENV,Help text with optional <metavar>'", arg.field.Name, tag)
		}

		aliases := strings.Split(tagParts[0], "|")
		if len(aliases) == 0 {
			return nil, fmt.Errorf("no names provided for field %s", arg.field.Name)
		}

		allowedPrefixes := []string{"--", "-", "$"}
		lowerAlphaNumeric := regexp.MustCompile(`^[a-z,0-9]+$`)
		for _, alias := range aliases {
			alias = strings.TrimSpace(alias)
			var aliasWithoutPrefix string
			for _, prefix := range allowedPrefixes {
				if strings.HasPrefix(alias, prefix) {
					aliasWithoutPrefix = strings.TrimPrefix(alias, prefix)
					break
				}
			}
			if aliasWithoutPrefix == "" {
				return nil, fmt.Errorf("invalid alias %q in tag for field %s", alias, arg.field.Name)
			}
			if !lowerAlphaNumeric.MatchString(aliasWithoutPrefix) {
				return nil, fmt.Errorf("invalid alias %q in tag for field %s", alias, arg.field.Name)
			}
			if _, exists := seenaliases[aliasWithoutPrefix]; exists {
				return nil, fmt.Errorf("duplicate name %q in tag for field %s", alias, arg.field.Name)
			}
			seenaliases[aliasWithoutPrefix] = true
		}

		arg.aliases = aliases
		var err error
		arg.metaVar, arg.help, err = extractMetaVar(tagParts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to extract metavar for field %s: %v", arg.field.Name, err)
		}
		if arg.help == "" {
			return nil, fmt.Errorf("missing help text in tag for field %s", arg.field.Name)
		}
		kind := arg.field.Type.Kind()
		if kind == reflect.Bool {
			if arg.metaVar != "" {
				return nil, fmt.Errorf("metavar is not allowed for boolean flags %s", arg.field.Name)
			}
		} else {
			if arg.metaVar == "" {
				arg.metaVar = fmt.Sprintf("<%s>", strings.ToLower(arg.field.Type.Name()))
			}
		}

		if err != nil {
			return nil, fmt.Errorf("failed to extract metavar for field %s: %v", arg.field.Name, err)
		}

		args = append(args, arg)
	}

	return args, nil
}
