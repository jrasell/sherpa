package in

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/cmd/helper"
	"github.com/jrasell/sherpa/pkg/api"
	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
	scaleCfg "github.com/jrasell/sherpa/pkg/config/scale"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "in",
		Short: "Perform scaling in actions on Nomad jobs and groups",
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

	os.Exit(runJobGroupScaleIn(client, args[0], scaleConfig.GroupName, scaleConfig.Count, scaleConfig.Meta))
}

func runJobGroupScaleIn(c *api.Client, job, group string, count int, meta map[string]string) int {
	resp, err := c.Scale().JobGroupIn(job, group, count, meta)
	if err != nil {
		fmt.Println("Error scaling in job group:", err)
		return sysexits.Software
	}

	out := []string{
		fmt.Sprintf("ID|%s", resp.ID),
		fmt.Sprintf("EvalID|%v", resp.EvaluationID),
	}

	fmt.Println(helper.FormatKV(out))
	return sysexits.OK
}
