package info

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/pkg/api"
	"github.com/ryanuber/columnize"
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
	clientCfg := api.DefaultConfig()

	client, err := api.NewClient(clientCfg)
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

	fmt.Println(formatList(out))
}

func formatList(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	return columnize.Format(in, columnConf)
}
