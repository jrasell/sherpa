package leader

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
		Use:   "leader",
		Short: "Check the HA status and current leader",
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

	leader, err := client.System().Leader()
	if err != nil {
		fmt.Println("Error calling server info:", err)
		os.Exit(sysexits.Software)
	}

	var out []string
	out = append(out, fmt.Sprintf("%s|%v", "Is Self", leader.IsSelf))
	out = append(out, fmt.Sprintf("%s|%s", "Leader Address", leader.LeaderAddress))
	out = append(out, fmt.Sprintf("%s|%s", "Leader Cluster Address", leader.LeaderClusterAddress))
	out = append(out, fmt.Sprintf("%s|%v", "HA Enabled", leader.HAEnabled))

	fmt.Println(helper.FormatList(out))
}
