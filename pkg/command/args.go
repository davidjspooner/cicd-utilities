package command

import (
	"fmt"
	"reflect"
	"strings"
)

type Option struct {
	Names   []string
	Help    string
	Default string
	field   reflect.StructField
}

func (o Option) Kind() string {
	switch o.field.Type.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Bool:
		return "boolean"
	default:
		return "unknown"
	}
}

func getDefinedOptions[T any](defaults *T) ([]Option, error) {
	// Get the type of the struct
	if defaults == nil {
		return nil, fmt.Errorf("defaults cannot be nil")
	}
	rDefaults := reflect.ValueOf(defaults).Elem()
	t := rDefaults.Type()
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", t.Kind())
	}

	// Create a slice to hold the defined arguments
	var args []Option

	// Iterate over the fields of the struct
	for i := 0; i < t.NumField(); i++ {
		arg := Option{}
		arg.field = t.Field(i)
		if !arg.field.IsExported() {
			continue
		}
		tag := arg.field.Tag.Get("arg")
		tagParts := strings.SplitN(tag, ",", 2)
		if len(tagParts) < 2 {
			return nil, fmt.Errorf("invalid tag format for field %s", arg.field.Name)
		}
		arg.Names = strings.Split(tagParts[0], "|")
		if len(arg.Names) == 0 {
			return nil, fmt.Errorf("no names provided for field %s", arg.field.Name)
		}
		arg.Help = tagParts[1]
		if arg.Help == "" {
			return nil, fmt.Errorf("no help text provided for field %s", arg.field.Name)
		}
		for i := 0; i < len(arg.Names); i++ {
			//count the number of dashes
			dashCount := 0
			for _, c := range arg.Names[i] {
				if c == '-' {
					dashCount++
				} else {
					break
				}
			}
			if dashCount > 2 {
				return nil, fmt.Errorf("too many dashes in name %s for field %s", arg.Names[i], arg.field.Name)
			}
			if dashCount == 0 {
				return nil, fmt.Errorf("no dashes in name %s for field %s", arg.Names[i], arg.field.Name)
			}
			if dashCount == 1 && len(arg.Names[i]) > 2 {
				return nil, fmt.Errorf("short name %s for field %s should be a single character", arg.Names[i], arg.field.Name)
			}
		}
		rdef := rDefaults.FieldByName(arg.field.Name)
		if !rdef.IsZero() {
			arg.Default = rdef.String()
		}

		args = append(args, arg)
	}
	return args, nil
}

func expandInputArgs(argsIn []string) (expandedArgs []string, err error) {

	expandedArgs = make([]string, 0, len(argsIn))

	for _, arg := range argsIn {
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg, "=", 2)
			expandedArgs = append(expandedArgs, parts[0])
			if len(parts) == 2 {
				expandedArgs = append(expandedArgs, parts[1])
			}
		} else if strings.HasPrefix(arg, "-") {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) == 1 {
				expandedArgs = append(expandedArgs, "-")
			}
			for _, letter := range parts[0][1:] {
				expandedArgs = append(expandedArgs, "-"+string(letter))
			}
			if len(parts) == 2 {
				expandedArgs = append(expandedArgs, parts[1])
			}
		} else {
			expandedArgs = append(expandedArgs, arg)
		}
	}
	return expandedArgs, nil
}
