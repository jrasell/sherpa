package server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const configKeyDebugEnable = "debug-enabled"

// GetDebugEnabled is used to identify whether the operator has enabled the server debug API
// endpoints.
func GetDebugEnabled() bool { return viper.GetBool(configKeyDebugEnable) }

// RegisterDebugConfig registers the CLI flags used to alter the server debug API endpoints.
func RegisterDebugConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyDebugEnable
			longOpt      = "debug-enabled"
			defaultValue = false
			description  = "Specifies if the debugging HTTP endpoints should be enabled"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
