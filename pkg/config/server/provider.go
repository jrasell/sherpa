package server

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyMetricProviderPrometheusAddr   = "metric-provider-prometheus-addr"
	configKeyMetricProviderInfluxDBAddr     = "metric-provider-influxdb-addr"
	configKeyMetricProviderInfluxDBUsername = "metric-provider-influxdb-username"
	configKeyMetricProviderInfluxDBPassword = "metric-provider-influxdb-password"
	configKeyMetricProviderInfluxDBInsecure = "metric-provider-influxdb-insecure"
)

type MetricProviderConfig struct {
	Prometheus *MetricProviderPrometheusConfig
	InfluxDB   *MetricProviderInfluxDBConfig
}

type MetricProviderPrometheusConfig struct {
	Addr string
}

type MetricProviderInfluxDBConfig struct {
	Addr     string
	Username string
	Password string
	Insecure bool
}

// MarshalZerologObject is the Zerolog marshaller which allow us to log the object.
func (mpc *MetricProviderConfig) MarshalZerologObject(e *zerolog.Event) {}

func GetMetricProviderConfig() *MetricProviderConfig {
	mpc := &MetricProviderConfig{}

	if promAddr := viper.GetString(configKeyMetricProviderPrometheusAddr); promAddr != "" {
		mpc.Prometheus = &MetricProviderPrometheusConfig{Addr: promAddr}
	}
	if influxDBAddr := viper.GetString(configKeyMetricProviderInfluxDBAddr); influxDBAddr != "" {
		mpc.InfluxDB = &MetricProviderInfluxDBConfig{
			Addr:     influxDBAddr,
			Username: viper.GetString(configKeyMetricProviderInfluxDBUsername),
			Password: viper.GetString(configKeyMetricProviderInfluxDBPassword),
			Insecure: viper.GetBool(configKeyMetricProviderInfluxDBInsecure),
		}
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
	{
		const (
			key          = configKeyMetricProviderInfluxDBAddr
			longOpt      = "metric-provider-influxdb-addr"
			defaultValue = ""
			description  = "The address of the InfluxDB server in the form <protocol>://<addr>:<port>"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
	{
		const (
			key          = configKeyMetricProviderInfluxDBUsername
			longOpt      = "metric-provider-influxdb-username"
			defaultValue = ""
			description  = "InfluxDB username"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
	{
		const (
			key          = configKeyMetricProviderInfluxDBPassword
			longOpt      = "metric-provider-influxdb-password"
			defaultValue = ""
			description  = "InfluxDB password"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
	{
		const (
			key          = configKeyMetricProviderInfluxDBInsecure
			longOpt      = "metric-provider-influxdb-insecure"
			defaultValue = false
			description  = "Skip TLS validation of InfluxDB server certificate"
		)

		flags.Bool(longOpt, false, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
