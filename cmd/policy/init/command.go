package init

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	initPolicyCountSection = "{\"Enabled\":true,\"Cooldown\":120,\"MaxCount\":16,\"MinCount\":4,\"ScaleOutCount\":2,\"ScaleInCount\":2,"
	thresholdPolicySection = "\"ScaleOutCPUPercentageThreshold\":75,\"ScaleOutMemoryPercentageThreshold\":75,\"ScaleInCPUPercentageThreshold\":30,\"ScaleInMemoryPercentageThreshold\":30}"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates an example job group scaling policy",
		Run: func(cmd *cobra.Command, args []string) {
			runInit(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runInit(_ *cobra.Command, _ []string) {
	fmt.Println(initPolicyCountSection + thresholdPolicySection)
}
