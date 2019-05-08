package logger

import (
	"fmt"
	"strings"
)

// Format represents our log format type.
type Format uint

// Our log format options.
const (
	FormatAuto Format = iota
	FormatZerolog
	FormatHuman
)

// String takes a Format type and converts to a string.
func (f Format) String() string {
	switch f {
	case FormatAuto:
		return "auto"
	case FormatZerolog:
		return "zerolog"
	case FormatHuman:
		return "human"
	default:
		panic(fmt.Sprintf("unknown log format: %d", f))
	}
}

func getLogFormat(logFormat string) (Format, error) {
	switch strLogFormat := strings.ToLower(logFormat); strLogFormat {
	case "auto":
		return FormatAuto, nil
	case "json", "zerolog":
		return FormatZerolog, nil
	case "human":
		return FormatHuman, nil
	default:
		return FormatAuto, fmt.Errorf("unsupported log format: %q", logFormat)
	}
}
