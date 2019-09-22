package system

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/cmd/system/health"
	"github.com/jrasell/sherpa/cmd/system/info"
	"github.com/jrasell/sherpa/cmd/system/leader"
	"github.com/jrasell/sherpa/cmd/system/metrics"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "Retrieve system information about a Sherpa server",
		Run: func(cmd *cobra.Command, args []string) {
			runSystem(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	if err := registerCommands(cmd); err != nil {
		fmt.Println("Error registering commands:", err)
		os.Exit(sysexits.Software)
	}

	return nil
}

func runSystem(cmd *cobra.Command, _ []string) {
	_ = cmd.Usage()
}

func registerCommands(rootCmd *cobra.Command) error {
	if err := info.RegisterCommand(rootCmd); err != nil {
		return err
	}

	if err := leader.RegisterCommand(rootCmd); err != nil {
		return err
	}

	if err := metrics.RegisterCommand(rootCmd); err != nil {
		return err
	}

	return health.RegisterCommand(rootCmd)
}
