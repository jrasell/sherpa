package out

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/pkg/api"
	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
	scaleCfg "github.com/jrasell/sherpa/pkg/config/scale"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "out",
		Short: "Perform scaling out actions on Nomad jobs and groups.",
		Run: func(cmd *cobra.Command, args []string) {
			runOut(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runOut(_ *cobra.Command, args []string) {
	switch {
	case len(args) < 1:
		fmt.Println("Not enough arguments, expected 1 arg got", len(args))
		os.Exit(sysexits.Usage)
	case len(args) > 1:
		fmt.Println("Too many arguments, expected 1 arg got", len(args))
		os.Exit(sysexits.Usage)
	}

	scaleConfig := scaleCfg.GetScaleConfig()

	clientConfig := clientCfg.GetConfig()
	mergedConfig := api.DefaultConfig(&clientConfig)

	client, err := api.NewClient(mergedConfig)
	if err != nil {
		fmt.Println("Error setting up Sherpa client:", err)
		os.Exit(sysexits.Software)
	}

	if scaleConfig.GroupName == "" {
		fmt.Println("Please specify a job group to scale")
		os.Exit(sysexits.Usage)
	}

	os.Exit(runJobGroupScaleOut(client, args[0], scaleConfig.GroupName, scaleConfig.Count))
}

func runJobGroupScaleOut(c *api.Client, job, group string, count int) int {
	resp, err := c.Scale().JobGroupOut(job, group, count)
	if err != nil {
		fmt.Println("Error scaling out job group:", err)
		return sysexits.Software
	}

	fmt.Println("Evaluation ID:", resp.EvaluationID)
	return sysexits.OK
}
