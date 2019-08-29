package info

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/cmd/helper"
	"github.com/jrasell/sherpa/pkg/api"
	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Retrieve information about a Sherpa server",
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

	info, err := client.System().Info()
	if err != nil {
		fmt.Println("Error calling server info:", err)
		os.Exit(sysexits.Software)
	}

	var out []string
	out = append(out, fmt.Sprintf("%s|%s", "Nomad Address", info.NomadAddress))
	out = append(out, fmt.Sprintf("%s|%s", "Policy Engine", info.PolicyEngine))
	out = append(out, fmt.Sprintf("%s|%s", "Policy Storage Backend", info.PolicyStorageBackend))
	out = append(out, fmt.Sprintf("%s|%v", "Internal AutoScaling Engine", info.InternalAutoScalingEngine))
	out = append(out, fmt.Sprintf("%s|%v", "Strict Policy Checking", info.StrictPolicyChecking))

	fmt.Println(helper.FormatList(out))
}
