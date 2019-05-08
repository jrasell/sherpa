package policy

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyPolicyGroupName = "policy-group-name"
)

type Config struct {
	GroupName string
}

func GetConfig() Config {
	return Config{
		GroupName: viper.GetString(configKeyPolicyGroupName),
	}
}

func RegisterConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyPolicyGroupName
			longOpt      = "policy-group-name"
			defaultValue = ""
			description  = "The job group to interact with"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
