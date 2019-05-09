package autoscale

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyAutoscalingEvaluationInterval        = "autoscaling-evaluation-interval"
	configKeyAutoscalingEvaluationIntervalDefault = 60
)

type Config struct {
	AutoscalingEvalInt int
}

func GetConfig() Config {
	return Config{
		AutoscalingEvalInt: viper.GetInt(configKeyAutoscalingEvaluationInterval),
	}
}

func RegisterConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyAutoscalingEvaluationInterval
			longOpt      = "autoscaling-evaluation-interval"
			defaultValue = configKeyAutoscalingEvaluationIntervalDefault
			description  = "The time in seconds between each run on the internal autoscaler"
		)

		flags.Int(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
