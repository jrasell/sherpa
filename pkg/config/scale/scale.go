package scale

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configKeyScaleCount           = "count"
	configKeyScaleGroupName       = "group-name"
	configKeyScaleMeta            = "meta"
	configScaleMetaPairSeparator  = ":"
	configScaleMetaValueSeparator = "="
)

type Config struct {
	Count     int
	GroupName string
	Meta      map[string]string
}

func GetScaleConfig() Config {
	return Config{
		Count:     viper.GetInt(configKeyScaleCount),
		GroupName: viper.GetString(configKeyScaleGroupName),
		Meta:      parseMetaMap(viper.GetString(configKeyScaleMeta)),
	}
}

func parseMetaMap(raw string) map[string]string {
	m := make(map[string]string)

	pairs := strings.Split(raw, configScaleMetaPairSeparator)
	for _, pair := range pairs {
		v := strings.Split(pair, configScaleMetaValueSeparator)
		if len(v) != 2 {
			continue
		}

		m[v[0]] = v[1]
	}

	return m
}

func RegisterScaleConfig(cmd *cobra.Command) {
	flags := cmd.PersistentFlags()

	{
		const (
			key          = configKeyScaleCount
			longOpt      = "count"
			defaultValue = 0
			description  = "The number by which to increment or decrement the job group"
		)

		flags.Int(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyScaleGroupName
			longOpt      = "group-name"
			defaultValue = ""
			description  = "The job group name to scale (required)"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}

	{
		const (
			key          = configKeyScaleMeta
			longOpt      = "meta"
			defaultValue = ""
			description  = "The meta parameters to record in the scaling event (key=value:key2=value2)"
		)

		flags.String(longOpt, defaultValue, description)
		_ = viper.BindPFlag(key, flags.Lookup(longOpt))
		viper.SetDefault(key, defaultValue)
	}
}
