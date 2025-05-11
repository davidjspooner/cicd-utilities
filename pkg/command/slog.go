package command

import (
	"fmt"
	"log/slog"
)

type LogOptions struct {
	Level   string `flag:"--loglevel,Log level (debug|info|warn|error)"`
	Verbose bool   `flag:"--verbose,Verbose output (alias for --loglevel debug)"`
}

func (options *LogOptions) Parse() (slog.Level, error) {

	var level slog.Level
	var err error
	if options == nil {
		options = &LogOptions{Level: "info"}
	}
	if options.Level == "" && options.Verbose {
		options.Level = "debug"
	}
	switch options.Level {
	case "debug":
		level = slog.LevelDebug
	case "info", "":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	case "silent":
		level = 100
	default:
		err = fmt.Errorf("invalid log level %q", options.Level)
		level = slog.LevelInfo
	}
	return level, err
}
