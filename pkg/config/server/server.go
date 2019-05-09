package server

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyHTTPServerBindAddr              = "bind-addr"
	configKeyHTTPServerBindAddrDefault       = "127.0.0.1"
	configKeyHTTPServerPort                  = "bind-port"
	configKeyHTTPServerPortDefault           = 8000
	configKeyConsulStorageBackendPathDefault = "sherpa/policies/"

	configKeyEnableAPIPolicyEngine       = "api-policy-engine-enabled"
	configKeyConsulStorageBackendEnabled = "consul-storage-backend-enabled"
	configKeyConsulStorageBackendPath    = "consul-storage-backend-path"
	configKeyEnableNomadMetaPolicyEngine = "nomad-meta-policy-engine-enabled"
	configKeyStrictPolicyChecking        = "strict-policy-checking-enabled"
	configKeyEnableInternalAutoScaler    = "internal-auto-scaler-enabled"
)

type Config struct {
	Bind                     string
	Port                     uint16
	APIPolicyEngine          bool
	NomadMetaPolicyEngine    bool
	StrictPolicyChecking     bool
	InternalAutoScaler       bool
	ConsulStorageBackend     bool
	ConsulStorageBackendPath string
}

func (c *Config) MarshalZerologObject(e *zerolog.Event) {
	e.Str("bind", c.Bind).
		Uint16("port", c.Port).
		Bool("api-policy-engine-enabled", c.APIPolicyEngine).
		Bool("nomad-meta-policy-engine-enabled", c.NomadMetaPolicyEngine).
		Bool("strict-policy-checking-enabled", c.StrictPolicyChecking).
		Bool("internal-autoscaler-enabled", c.InternalAutoScaler).
		Bool("consul-storage-backend-enabled", c.ConsulStorageBackend).
		Str("consul-storage-backend-path", c.ConsulStorageBackendPath)
}

func GetConfig() Config {
	return Config{
		Bind:                     viper.GetString(configKeyHTTPServerBindAddr),
		Port:                     uint16(viper.GetInt(configKeyHTTPServerPort)),
		APIPolicyEngine:          viper.GetBool(configKeyEnableAPIPolicyEngine),
		NomadMetaPolicyEngine:    viper.GetBool(configKeyEnableNomadMetaPolicyEngine),
		StrictPolicyChecking:     viper.GetBool(configKeyStrictPolicyChecking),
		InternalAutoScaler:       viper.GetBool(configKeyEnableInternalAutoScaler),
		ConsulStorageBackend:     viper.GetBool(configKeyConsulStorageBackendEnabled),
		ConsulStorageBackendPath: viper.GetString(configKeyConsulStorageBackendPath),
	}
}

func RegisterConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyHTTPServerBindAddr
			longOpt      = "bind-addr"
			defaultValue = configKeyHTTPServerBindAddrDefault
			description  = "The HTTP server address to bind to"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyHTTPServerPort
			longOpt      = "bind-port"
			defaultValue = configKeyHTTPServerPortDefault
			description  = "The HTTP server port to bind to"
		)

		flags.Uint16(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyEnableAPIPolicyEngine
			longOpt      = "api-policy-engine-enabled"
			defaultValue = true
			description  = "Enable the Sherpa API to manage scaling policies"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyEnableNomadMetaPolicyEngine
			longOpt      = "nomad-meta-policy-engine-enabled"
			defaultValue = false
			description  = "Enable Nomad job meta lookups to manage scaling policies"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyStrictPolicyChecking
			longOpt      = "strict-policy-checking-enabled"
			defaultValue = true
			description  = "When enabled, all scaling activities must pass through policy checks"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyEnableInternalAutoScaler
			longOpt      = "internal-auto-scaler-enabled"
			defaultValue = false
			description  = "Enable the internal autoscaling engine"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyConsulStorageBackendEnabled
			longOpt      = "consul-storage-backend-enabled"
			defaultValue = false
			description  = "Use Consul as a storage backend when using the API policy engine"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyConsulStorageBackendPath
			longOpt      = "consul-storage-backend-path"
			defaultValue = configKeyConsulStorageBackendPathDefault
			description  = "The Consul KV path that will be used to store policies"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
