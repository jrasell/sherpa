package health

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/pkg/api"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Retrieve health information of a Sherpa server",
		Run: func(cmd *cobra.Command, args []string) {
			runHealth(cmd, args)
		},
	}

	rootCmd.AddCommand(cmd)

	return nil
}

func runHealth(_ *cobra.Command, _ []string) {
	clientCfg := api.DefaultConfig()

	client, err := api.NewClient(clientCfg)
	if err != nil {
		fmt.Println("Error setting up Sherpa client:", err)
		os.Exit(sysexits.Software)
	}

	health, err := client.System().Health()
	if err != nil {
		fmt.Println("Error calling server health:", err)
		os.Exit(sysexits.Software)
	}

	fmt.Println("Sherpa server status:", health.Status)
}
