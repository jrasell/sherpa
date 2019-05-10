package sherpa

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeySherpaAddr        = "addr"
	configKeySherpaAddrDefault = "http://127.0.0.1:8000"

	configKeySherpaClientCertPath    = "client-cert-path"
	configKeySherpaClientCertKeyPath = "client-cert-key-path"
	configKeySherpaCAPath            = "client-ca-path"
)

type Config struct {
	Addr        string
	CertPath    string
	CertKeyPath string
	CAPath      string
}

func GetConfig() Config {
	return Config{
		Addr:        viper.GetString(configKeySherpaAddr),
		CertPath:    viper.GetString(configKeySherpaClientCertPath),
		CertKeyPath: viper.GetString(configKeySherpaClientCertKeyPath),
		CAPath:      viper.GetString(configKeySherpaCAPath),
	}
}

func RegisterConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeySherpaAddr
			longOpt      = "addr"
			defaultValue = configKeySherpaAddrDefault
			description  = "The HTTP(S) address of the sherpa server"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeySherpaClientCertPath
			longOpt      = "client-cert-path"
			defaultValue = ""
			description  = "Path to a PEM encoded client certificate for TLS authentication to the Sherpa server"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeySherpaClientCertKeyPath
			longOpt      = "client-cert-key-path"
			defaultValue = ""
			description  = "Path to an unencrypted PEM encoded private key matching the client certificate"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeySherpaCAPath
			longOpt      = "client-ca-path"
			defaultValue = ""
			description  = "Path to a PEM encoded CA cert file to use to verify the Sherpa server SSL certificate"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
