package command

import (
	"encoding"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
var textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

type Flag struct {
	aliases      []string
	help         string
	metaVar      string
	defaultValue string
	fieldPath    []reflect.StructField
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
	return flag.leaf().Type.Kind()
}
func (flag *Flag) TypeName() string {
	return flag.leaf().Type.Name()
}
func (flag *Flag) leaf() *reflect.StructField {
	if len(flag.fieldPath) == 0 {
		return nil
	}
	return &flag.fieldPath[len(flag.fieldPath)-1]
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
	return getRFlagDefinitions(rDefaults, nil)
}

func getRFlagDefinitions(rDefaults reflect.Value, fieldPath []reflect.StructField) ([]Flag, error) {
	t := rDefaults.Type()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", t.Kind())
	}

	var args []Flag
	seenaliases := map[string]bool{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		currentFieldPath := append(fieldPath, field)

		if field.Type.Kind() == reflect.Struct {
			// Recurse into sub-structures, including anonymous ones
			subStruct := rDefaults.Field(i)
			subArgs, err := getRFlagDefinitions(subStruct, currentFieldPath)
			if err != nil {
				return nil, fmt.Errorf("failed to process sub-structure %s: %v", field.Name, err)
			}
			args = append(args, subArgs...)
			continue
		}

		arg := Flag{}
		arg.fieldPath = currentFieldPath
		if !field.IsExported() {
			continue
		}

		tag := field.Tag.Get("flag")
		if strings.TrimSpace(tag) == "" {
			continue // no tag, skip
		}

		tagParts := strings.SplitN(tag, ",", 2)
		if len(tagParts) != 2 {
			return nil, fmt.Errorf("invalid tag format for field %s: %q. Expected something like '--opt|-o|$ENV,Help text with optional <metavar>'", field.Name, tag)
		}

		aliases := strings.Split(tagParts[0], "|")
		if len(aliases) == 0 {
			return nil, fmt.Errorf("no names provided for field %s", field.Name)
		}

		allowedPrefixes := []string{"--", "-", "$"}
		lowerAlphaNumeric := regexp.MustCompile(`^[a-z,0-9-]+$`)
		upperCase := regexp.MustCompile(`^[A-Z,0-9_]+$`)
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
				return nil, fmt.Errorf("invalid alias %q in tag for field %s", alias, field.Name)
			}
			if strings.HasPrefix(alias, "$") {
				if !upperCase.MatchString(aliasWithoutPrefix) {
					return nil, fmt.Errorf("invalid alias %q in tag for field %s", alias, field.Name)
				}
			} else {
				if !lowerAlphaNumeric.MatchString(aliasWithoutPrefix) {
					return nil, fmt.Errorf("invalid alias %q in tag for field %s", alias, field.Name)
				}
			}
			if _, exists := seenaliases[aliasWithoutPrefix]; exists {
				return nil, fmt.Errorf("duplicate name %q in tag for field %s", alias, field.Name)
			}
			seenaliases[aliasWithoutPrefix] = true
		}

		arg.aliases = aliases
		var err error
		arg.metaVar, arg.help, err = extractMetaVar(tagParts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to extract metavar for field %s: %v", field.Name, err)
		}
		if arg.help == "" {
			return nil, fmt.Errorf("missing help text in tag for field %s", field.Name)
		}
		kind := field.Type.Kind()
		if kind == reflect.Bool {
			if arg.metaVar != "" {
				return nil, fmt.Errorf("metavar is not allowed for boolean flags %s", field.Name)
			}
		} else {
			if arg.metaVar == "" {
				arg.metaVar = fmt.Sprintf("<%s>", strings.ToLower(field.Type.Name()))
			}
		}

		// Check if the field implements encoding.TextUnmarshaler or encoding.TextMarshaler
		fieldValue := rDefaults.FieldByName(field.Name)
		if fieldValue.Addr().Type().Implements(textUnmarshalerType) {
			arg.defaultValue = "(implements TextUnmarshaler)"
		} else if fieldValue.Addr().Type().Implements(textMarshalerType) {
			marshaler := fieldValue.Addr().Interface().(encoding.TextMarshaler)
			if marshaler != nil {
				text, err := marshaler.MarshalText()
				if err == nil {
					arg.defaultValue = string(text)
				}
			}
		} else if (!fieldValue.IsZero()) && fieldValue.CanInterface() {
			arg.defaultValue = fmt.Sprintf("%v", fieldValue.Interface())
		}

		args = append(args, arg)
	}

	return args, nil
}
