package log

import (
	"os"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel  string
	LogFormat string
	EnableDev bool
	UseColor  bool
}

const (
	configKeyLogLevel     = "log-level"
	configKeyLogFormat    = "log-format"
	configKeyLogEnableDev = "log-enable-dev"
	configKeyUseColor     = "log-use-color"
)

func GetConfig() Config {
	return Config{
		LogLevel:  viper.GetString(configKeyLogLevel),
		LogFormat: viper.GetString(configKeyLogFormat),
		EnableDev: viper.GetBool(configKeyLogEnableDev),
		UseColor:  viper.GetBool(configKeyUseColor),
	}
}

func RegisterConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyLogLevel
			longOpt      = "log-level"
			defaultValue = "info"
			description  = "Change the level used for logging"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyLogFormat
			longOpt      = "log-format"
			defaultValue = "auto"
			description  = `Specify the log format ("auto", "zerolog" or "human")`
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyLogEnableDev
			longOpt      = "log-enable-dev"
			defaultValue = false
			description  = "Log with file:line of the caller"
		)
		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key         = configKeyUseColor
			longOpt     = "log-use-color"
			description = "Use ANSI colors in logging output"
		)
		defaultValue := false
		if isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd()) {
			defaultValue = true
		}

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
