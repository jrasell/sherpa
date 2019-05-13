package server

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyBindAddrDefault                     = "127.0.0.1"
	configKeyBindPortDefault                     = 8000
	configKeyStorageBackendConsulPathDefault     = "sherpa/policies/"
	configKeyAutoscalerEvaluationIntervalDefault = 60

	configKeyBindAddr                          = "bind-addr"
	configKeyBindPort                          = "bind-port"
	configKeyAutoscalerEnabled                 = "autoscaler-enabled"
	configKeyAutoscalerEvaluationInterval      = "autoscaler-evaluation-interval"
	configKeyPolicyEngineAPIEnabled            = "policy-engine-api-enabled"
	configKeyPolicyEngineNomadMetaEnabled      = "policy-engine-nomad-meta-enabled"
	configKeyPolicyEngineStrictCheckingEnabled = "policy-engine-strict-checking-enabled"
	configKeyStorageBackendConsulEnabled       = "storage-consul-enabled"
	configKeyStorageBackendConsulPath          = "storage-consul-path"
)

type Config struct {
	Bind                         string
	ConsulStorageBackendPath     string
	Port                         uint16
	APIPolicyEngine              bool
	NomadMetaPolicyEngine        bool
	StrictPolicyChecking         bool
	InternalAutoScaler           bool
	ConsulStorageBackend         bool
	InternalAutoScalerEvalPeriod int
}

func (c *Config) MarshalZerologObject(e *zerolog.Event) {
	e.Str(configKeyBindAddr, c.Bind).
		Uint16(configKeyBindPort, c.Port).
		Bool(configKeyPolicyEngineAPIEnabled, c.APIPolicyEngine).
		Bool(configKeyPolicyEngineNomadMetaEnabled, c.NomadMetaPolicyEngine).
		Bool(configKeyPolicyEngineStrictCheckingEnabled, c.StrictPolicyChecking).
		Bool(configKeyAutoscalerEnabled, c.InternalAutoScaler).
		Int(configKeyAutoscalerEvaluationInterval, c.InternalAutoScalerEvalPeriod).
		Bool(configKeyStorageBackendConsulEnabled, c.ConsulStorageBackend).
		Str(configKeyStorageBackendConsulPath, c.ConsulStorageBackendPath)
}

func GetConfig() Config {
	return Config{
		Bind:                         viper.GetString(configKeyBindAddr),
		Port:                         uint16(viper.GetInt(configKeyBindPort)),
		APIPolicyEngine:              viper.GetBool(configKeyPolicyEngineAPIEnabled),
		NomadMetaPolicyEngine:        viper.GetBool(configKeyPolicyEngineNomadMetaEnabled),
		StrictPolicyChecking:         viper.GetBool(configKeyPolicyEngineStrictCheckingEnabled),
		InternalAutoScaler:           viper.GetBool(configKeyAutoscalerEnabled),
		InternalAutoScalerEvalPeriod: viper.GetInt(configKeyAutoscalerEvaluationInterval),
		ConsulStorageBackend:         viper.GetBool(configKeyStorageBackendConsulEnabled),
		ConsulStorageBackendPath:     viper.GetString(configKeyStorageBackendConsulPath),
	}
}

func RegisterConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyBindAddr
			longOpt      = "bind-addr"
			defaultValue = configKeyBindAddrDefault
			description  = "The HTTP server address to bind to"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyBindPort
			longOpt      = "bind-port"
			defaultValue = configKeyBindPortDefault
			description  = "The HTTP server port to bind to"
		)

		flags.Uint16(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyPolicyEngineAPIEnabled
			longOpt      = "policy-engine-api-enabled"
			defaultValue = true
			description  = "Enable the Sherpa API to manage scaling policies"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyPolicyEngineNomadMetaEnabled
			longOpt      = "policy-engine-nomad-meta-enabled"
			defaultValue = false
			description  = "Enable Nomad job meta lookups to manage scaling policies"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyPolicyEngineStrictCheckingEnabled
			longOpt      = "policy-engine-strict-checking-enabled"
			defaultValue = true
			description  = "When enabled, all scaling activities must pass through policy checks"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyAutoscalerEnabled
			longOpt      = "autoscaler-enabled"
			defaultValue = false
			description  = "Enable the internal autoscaling engine"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyAutoscalerEvaluationInterval
			longOpt      = "autoscaler-evaluation-interval"
			defaultValue = configKeyAutoscalerEvaluationIntervalDefault
			description  = "The time period in seconds between autoscaling evaluation runs"
		)

		flags.Int(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyStorageBackendConsulEnabled
			longOpt      = "storage-consul-enabled"
			defaultValue = false
			description  = "Use Consul as a storage backend when using the API policy engine"
		)

		flags.Bool(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyStorageBackendConsulPath
			longOpt      = "storage-consul-path"
			defaultValue = configKeyStorageBackendConsulPathDefault
			description  = "The Consul KV path that will be used to store policies"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
