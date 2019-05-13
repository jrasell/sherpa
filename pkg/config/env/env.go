package env

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RegisterCobra(cmd *cobra.Command) {
	viper.SetEnvPrefix(strings.ToUpper(cmd.Name()))
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}
