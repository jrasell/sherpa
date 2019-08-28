package metrics

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/cmd/helper"
	"github.com/jrasell/sherpa/pkg/api"
	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

const (
	outputHeader = "Name|Type|Value"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Retrieve metrics from a Sherpa server",
		Run: func(cmd *cobra.Command, args []string) {
			runInfo(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runInfo(_ *cobra.Command, _ []string) {
	clientConfig := clientCfg.GetConfig()
	mergedConfig := api.DefaultConfig(&clientConfig)

	client, err := api.NewClient(mergedConfig)
	if err != nil {
		fmt.Println("Error setting up Sherpa client:", err)
		os.Exit(sysexits.Software)
	}

	metrics, err := client.System().Metrics()
	if err != nil {
		fmt.Println("Error calling server metrics:", err)
		os.Exit(sysexits.Software)
	}

	out := []string{outputHeader}

	for i := range metrics.Gauges {
		out = append(out, fmt.Sprintf("%s|%s|%v",
			metrics.Gauges[i].Name, "Gauge", metrics.Gauges[i].Value))
	}

	for i := range metrics.Counters {
		out = append(out, fmt.Sprintf("%s|%s|%v",
			metrics.Counters[i].Name, "Counter", metrics.Counters[i].Mean))
	}

	for i := range metrics.Samples {
		out = append(out, fmt.Sprintf("%s|%s|%v",
			metrics.Samples[i].Name, "Counter", metrics.Samples[i].Mean))
	}

	// If there are no metrics to print (happens during initial server startup)
	// then we don't want to just print the header so perform a check so the
	// CLI is nice and tidy.
	if len(out) > 1 {
		fmt.Println(helper.FormatList(out))
	}
}
