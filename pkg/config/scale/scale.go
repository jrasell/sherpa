package scale

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyScaleCount     = "count"
	configKeyScaleGroupName = "group-name"
)

type Config struct {
	Count     int
	GroupName string
}

func GetScaleConfig() Config {
	return Config{
		Count:     viper.GetInt(configKeyScaleCount),
		GroupName: viper.GetString(configKeyScaleGroupName),
	}
}

func RegisterScaleConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyScaleCount
			longOpt      = "count"
			defaultValue = 0
			description  = "The number by which to increment or decrement the job group"
		)

		flags.Int(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyScaleGroupName
			longOpt      = "group-name"
			defaultValue = ""
			description  = "The job group name to scale (required)"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
