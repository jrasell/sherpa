package sherpa

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeySherpaAddr        = "sherpa-addr"
	configKeySherpaAddrDefault = "http://127.0.0.1:8000"
)

type Config struct {
	Addr string
}

func GetSherpaConfig() Config {
	return Config{
		Addr: viper.GetString(configKeySherpaAddr),
	}
}

func RegisterConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeySherpaAddr
			longOpt      = "sherpa-addr"
			defaultValue = configKeySherpaAddrDefault
			description  = "The HTTP(S) address of the sherpa server"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
