package delete

import (
	"fmt"
	"os"
	"strings"

	"github.com/jrasell/sherpa/pkg/api"
	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
	policyCfg "github.com/jrasell/sherpa/pkg/config/policy"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes a scaling policy from Sherpa",
		Run: func(cmd *cobra.Command, args []string) {
			runDelete(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runDelete(_ *cobra.Command, args []string) {
	switch {
	case len(args) < 1:
		fmt.Println("Not enough arguments, expected 1 arg got", len(args))
		os.Exit(sysexits.Usage)
	case len(args) > 1:
		fmt.Println("Too many arguments, expected 1 arg got", len(args))
		os.Exit(sysexits.Usage)
	}

	clientConfig := clientCfg.GetConfig()
	mergedConfig := api.DefaultConfig(&clientConfig)

	client, err := api.NewClient(mergedConfig)
	if err != nil {
		fmt.Println("Error setting up Sherpa client:", err)
		os.Exit(sysexits.Software)
	}

	name := strings.TrimSpace(strings.ToLower(args[0]))

	policyConfig := policyCfg.GetConfig()
	if policyConfig.GroupName != "" {
		os.Exit(runDeleteJobGroupPolicy(client, name, policyConfig.GroupName))
	}
	os.Exit(runDeleteJobPolicy(client, name))
}

func runDeleteJobPolicy(c *api.Client, job string) int {
	err := c.Policies().DeleteJobPolicy(job)
	if err != nil {
		fmt.Println("Error deleting job scaling policy:", err)
		return sysexits.Software
	}

	fmt.Println("Successfully deleted job scaling policy")
	return sysexits.OK
}

func runDeleteJobGroupPolicy(c *api.Client, job, group string) int {
	err := c.Policies().DeleteJobGroupPolicy(job, group)
	if err != nil {
		fmt.Println("Error deleting job group scaling policy:", err)
		return sysexits.Software
	}

	fmt.Println("Successfully deleted job group scaling policy")
	return sysexits.OK
}
