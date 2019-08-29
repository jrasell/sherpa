package scale

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/cmd/scale/in"
	"github.com/jrasell/sherpa/cmd/scale/out"
	"github.com/jrasell/sherpa/cmd/scale/status"
	scaleCfg "github.com/jrasell/sherpa/pkg/config/scale"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "scale",
		Short: "Perform scaling actions against a Nomad job and check scaling status",
		Run: func(cmd *cobra.Command, args []string) {
			runScale(cmd, args)
		},
	}

	rootCmd.AddCommand(cmd)
	scaleCfg.RegisterScaleConfig(cmd)

	if err := registerCommands(cmd); err != nil {
		fmt.Println("Error registering commands:", err)
		os.Exit(sysexits.Software)
	}
	return nil
}

func runScale(cmd *cobra.Command, _ []string) {
	_ = cmd.Usage()
}

func registerCommands(cmd *cobra.Command) error {
	if err := in.RegisterCommand(cmd); err != nil {
		return err
	}

	if err := status.RegisterCommand(cmd); err != nil {
		return err
	}

	return out.RegisterCommand(cmd)
}
