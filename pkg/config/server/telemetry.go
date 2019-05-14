package server

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyTelemetryStatsiteAddress = "telemetry-statsite-address"
	configKeyTelemetryStatsdAddress   = "telemetry-statsd-address"
)

// TelemetryConfig is the server Telemetry configuration struct.
type TelemetryConfig struct {
	StatsiteAddr string
	StatsdAddr   string
}

// MarshalZerologObject is the Zerolog marshaller which allow us to log the
// object.
func (c *TelemetryConfig) MarshalZerologObject(e *zerolog.Event) {
	e.Str(configKeyTelemetryStatsiteAddress, c.StatsiteAddr).
		Str(configKeyTelemetryStatsdAddress, c.StatsdAddr)
}

// GetTelemetryConfig hydrates the telemetry config struct.
func GetTelemetryConfig() TelemetryConfig {
	return TelemetryConfig{
		StatsiteAddr: viper.GetString(configKeyTelemetryStatsiteAddress),
		StatsdAddr:   viper.GetString(configKeyTelemetryStatsdAddress),
	}
}

// RegisterTelemetryConfig is used by a Cobra command to register the Telemetry
// CLI flags.
func RegisterTelemetryConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyTelemetryStatsiteAddress
			longOpt      = "telemetry-statsite-address"
			defaultValue = ""
			description  = "Specifies the address of a statsite server to forward metrics data to"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyTelemetryStatsdAddress
			longOpt      = "telemetry-statsd-address"
			defaultValue = ""
			description  = "Specifies the address of a statsd server to forward metrics to"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
