package server

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyMetricProviderPrometheusAddr = "metric-provider-prometheus-addr"
)

type MetricProviderConfig struct {
	Prometheus *MetricProviderPrometheusConfig
}

type MetricProviderPrometheusConfig struct {
	Addr string
}

// MarshalZerologObject is the Zerolog marshaller which allow us to log the object.
func (mpc *MetricProviderConfig) MarshalZerologObject(e *zerolog.Event) {}

func GetMetricProviderConfig() *MetricProviderConfig {
	mpc := &MetricProviderConfig{}

	if promAddr := viper.GetString(configKeyMetricProviderPrometheusAddr); promAddr != "" {
		mpc.Prometheus = &MetricProviderPrometheusConfig{Addr: promAddr}
	}

	return mpc
}

func RegisterMetricProviderConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyMetricProviderPrometheusAddr
			longOpt      = "metric-provider-prometheus-addr"
			defaultValue = ""
			description  = "The address of the Prometheus endpoint in the form <protocol>://<addr>:<port>"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
