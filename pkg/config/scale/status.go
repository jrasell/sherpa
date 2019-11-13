package scale

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyScaleStatusLatest = "latest"
)

type StatusConfig struct {
	Latest bool
}

func GetScaleStatusConfig() *StatusConfig {
	return &StatusConfig{
		Latest: viper.GetBool(configKeyScaleStatusLatest),
	}
}

func RegisterScaleStatusConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyScaleStatusLatest
			longOpt      = "latest"
			defaultValue = false
			description  = "List the latest scaling event for each job group only"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
