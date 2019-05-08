package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"os"
	"time"

	logCfg "github.com/jrasell/sherpa/pkg/config/log"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// Use a log format that resembles time.RFC3339Nano but includes all trailing
	// zeros so that we get fixed-width logging.
	logTimeFormat = "2006-01-02T15:04:05.000000000Z07:00"
)

var stdLogger *stdlog.Logger

func init() {
	// Initialize zerolog with a set set of defaults.  Re-initialization of
	// logging with user-supplied configuration parameters happens in Setup().

	// os.Stderr isn't guaranteed to be thread-safe, wrap in a sync writer.  Files
	// are guaranteed to be safe, terminals are not.
	w := zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: true,
	}
	zlog := zerolog.New(zerolog.SyncWriter(w)).With().Timestamp().Logger()

	zerolog.DurationFieldUnit = time.Microsecond
	zerolog.DurationFieldInteger = true
	zerolog.TimeFieldFormat = logTimeFormat
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = zlog

	stdlog.SetFlags(0)
	stdlog.SetOutput(zlog)
}

// Setup configures the sherpa logging based on user configuration.
func Setup(config logCfg.Config) error {
	logLevel, err := setLogLevel(config.LogLevel)
	if err != nil {
		return errors.Wrap(err, "unable to set log level")
	}

	var logWriter io.Writer = os.Stderr

	logFmt, err := getLogFormat(config.LogFormat)
	if err != nil {
		return errors.Wrap(err, "unable to parse log format")
	}

	if logFmt == FormatAuto {
		if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
			logFmt = FormatHuman
		} else {
			logFmt = FormatZerolog
		}
	}

	var zlog zerolog.Logger
	switch logFmt {
	case FormatZerolog:
		zlog = zerolog.New(logWriter).With().Timestamp().Logger()
	case FormatHuman:
		useColor := config.UseColor
		w := zerolog.ConsoleWriter{
			Out:     logWriter,
			NoColor: !useColor,
		}
		zlog = zerolog.New(w).With().Timestamp().Logger()
	default:
		return fmt.Errorf("unsupported log format: %q", logFmt)
	}

	log.Logger = zlog

	stdlog.SetFlags(0)
	stdlog.SetOutput(zlog)
	stdLogger = &stdlog.Logger{}

	// In order to prevent random libraries from hooking the standard logger and
	// filling the logger with garbage, discard all log entries.  At debug level,
	// however, let it all through.
	if logLevel != LevelDebug {
		stdLogger.SetOutput(ioutil.Discard)
	} else {
		stdLogger.SetOutput(zlog)
	}

	if config.EnableDev {
		log.Logger = log.With().Caller().Logger()
	}

	return nil
}
