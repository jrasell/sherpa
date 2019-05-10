package main

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/cmd/policy"
	"github.com/jrasell/sherpa/cmd/scale"
	"github.com/jrasell/sherpa/cmd/server"
	"github.com/jrasell/sherpa/cmd/system"
	"github.com/jrasell/sherpa/pkg/build"
	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
	envCfg "github.com/jrasell/sherpa/pkg/config/env"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "sherpa",
		Short: `
Sherpa is a fast and flexible job scaler for HashiCorp Nomad, capable of
running in a number of different modes to suit your needs.
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Version: build.GetVersion(),
	}
	envCfg.RegisterCobra(rootCmd)
	clientCfg.RegisterConfig(rootCmd)

	if err := registerCommands(rootCmd); err != nil {
		fmt.Println("error registering commands:", err)
		os.Exit(sysexits.Software)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(sysexits.Software)
	}
}

func registerCommands(rootCmd *cobra.Command) error {
	if err := server.RegisterCommand(rootCmd); err != nil {
		return err
	}

	if err := system.RegisterCommand(rootCmd); err != nil {
		return err
	}

	if err := scale.RegisterCommand(rootCmd); err != nil {
		return err
	}

	return policy.RegisterCommand(rootCmd)
}
