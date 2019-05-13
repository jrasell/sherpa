package server

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyServerTLSCertPath    = "tls-cert-path"
	configKeyServerTLSCertKeyPath = "tls-cert-key-path"
)

type TLSConfig struct {
	CertPath    string
	CertKeyPath string
}

func (c *TLSConfig) MarshalZerologObject(e *zerolog.Event) {
	e.Str(configKeyServerTLSCertPath, c.CertPath).
		Str(configKeyServerTLSCertKeyPath, c.CertKeyPath)
}

func GetTLSConfig() TLSConfig {
	return TLSConfig{
		CertPath:    viper.GetString(configKeyServerTLSCertPath),
		CertKeyPath: viper.GetString(configKeyServerTLSCertKeyPath),
	}
}

func RegisterTLSConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyServerTLSCertPath
			longOpt      = "tls-cert-path"
			defaultValue = ""
			description  = "Path to the TLS certificate for the Sherpa server"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyServerTLSCertKeyPath
			longOpt      = "tls-cert-key-path"
			defaultValue = ""
			description  = "Path to the TLS certificate key for the Sherpa server"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
