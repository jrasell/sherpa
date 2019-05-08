package in

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/pkg/api"
	scaleCfg "github.com/jrasell/sherpa/pkg/config/scale"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "in",
		Short: "Perform scaling in actions on Nomad jobs and groups.",
		Run: func(cmd *cobra.Command, args []string) {
			runIn(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runIn(_ *cobra.Command, args []string) {
	switch {
	case len(args) < 1:
		fmt.Println("Not enough arguments, expected 1 arg got", len(args))
		os.Exit(sysexits.Usage)
	case len(args) > 1:
		fmt.Println("Too many arguments, expected 1 arg got", len(args))
		os.Exit(sysexits.Usage)
	}

	clientCfg := api.DefaultConfig()
	scaleConfig := scaleCfg.GetScaleConfig()

	client, err := api.NewClient(clientCfg)
	if err != nil {
		fmt.Println("Error setting up Sherpa client:", err)
		os.Exit(sysexits.Software)
	}

	if scaleConfig.GroupName == "" {
		fmt.Println("Please specify a job group to scale")
		os.Exit(sysexits.Usage)
	}

	os.Exit(runJobGroupScaleIn(client, args[0], scaleConfig.GroupName, scaleConfig.Count))
}

func runJobGroupScaleIn(c *api.Client, job, group string, count int) int {
	resp, err := c.Scale().JobGroupIn(job, group, count)
	if err != nil {
		fmt.Println("Error scaling in job group:", err)
		return sysexits.Software
	}

	fmt.Println("Evaluation ID:", resp.EvaluationID)
	return sysexits.OK
}
