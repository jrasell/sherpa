package init

import (
	"fmt"

	"github.com/spf13/cobra"
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
	fmt.Println("{\"Enabled\": true,\"MaxCount\": 16,\"MinCount\": 4,\"ScaleOutCount\": 2,\"ScaleInCount\": 2}")
}
