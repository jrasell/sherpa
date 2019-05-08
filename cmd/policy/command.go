package policy

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/cmd/policy/delete"
	"github.com/jrasell/sherpa/cmd/policy/list"
	"github.com/jrasell/sherpa/cmd/policy/read"
	"github.com/jrasell/sherpa/cmd/policy/write"
	policyCfg "github.com/jrasell/sherpa/pkg/config/policy"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Interact with scaling policies",
		Run: func(cmd *cobra.Command, args []string) {
			runPolicy(cmd, args)
		},
	}

	rootCmd.AddCommand(cmd)
	policyCfg.RegisterConfig(cmd)

	if err := registerCommands(cmd); err != nil {
		fmt.Println("Error registering commands:", err)
		os.Exit(sysexits.Software)
	}
	return nil
}

func runPolicy(cmd *cobra.Command, _ []string) {
	_ = cmd.Usage()
}

func registerCommands(cmd *cobra.Command) error {
	if err := list.RegisterCommand(cmd); err != nil {
		return err
	}

	if err := delete.RegisterCommand(cmd); err != nil {
		return err
	}

	if err := write.RegisterCommand(cmd); err != nil {
		return err
	}

	return read.RegisterCommand(cmd)
}
