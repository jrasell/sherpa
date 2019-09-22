package server

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyClusterAdvertiseAddrDefault = "http://127.0.0.1:8000"
	configKeyClusterAdvertiseAddr        = "cluster-advertise-addr"
	configKeyClusterName                 = "cluster-name"
)

type ClusterConfig struct {
	Addr string
	Name string
}

// MarshalZerologObject is the Zerolog marshaller which allow us to log the
// object.
func (c *ClusterConfig) MarshalZerologObject(e *zerolog.Event) {
	e.Str(configKeyClusterAdvertiseAddr, c.Addr).
		Str(configKeyClusterName, c.Name)
}

func GetClusterConfig() ClusterConfig {
	return ClusterConfig{
		Addr: viper.GetString(configKeyClusterAdvertiseAddr),
		Name: viper.GetString(configKeyClusterName),
	}
}

func RegisterClusterConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyClusterAdvertiseAddr
			longOpt      = "cluster-advertise-addr"
			defaultValue = configKeyClusterAdvertiseAddrDefault
			description  = "The Sherpa server advertise address used for NAT traversal on HTTP redirects"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyClusterName
			longOpt      = "cluster-name"
			defaultValue = ""
			description  = "Specifies the identifier for the Sherpa cluster"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
