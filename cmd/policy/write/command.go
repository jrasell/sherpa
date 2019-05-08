package write

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jrasell/sherpa/pkg/api"
	policyCfg "github.com/jrasell/sherpa/pkg/config/policy"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "write",
		Short: "Uploads a policy from file",
		Run: func(cmd *cobra.Command, args []string) {
			runWrite(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runWrite(_ *cobra.Command, args []string) {
	switch {
	case len(args) < 2:
		fmt.Println("Not enough arguments, expected 2 args got", len(args))
		os.Exit(sysexits.Usage)
	case len(args) > 2:
		fmt.Println("Too many arguments, expected 2 args got", len(args))
		os.Exit(sysexits.Usage)
	}

	path := strings.TrimSpace(args[1])

	b, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading scaling policy file:", err)
		os.Exit(sysexits.Software)
	}

	clientCfg := api.DefaultConfig()

	client, err := api.NewClient(clientCfg)
	if err != nil {
		fmt.Println("Error setting up Sherpa client:", err)
		os.Exit(sysexits.Software)
	}

	name := strings.TrimSpace(strings.ToLower(args[0]))

	policyConfig := policyCfg.GetConfig()
	if policyConfig.GroupName != "" {
		os.Exit(runJobGroupWrite(client, name, policyConfig.GroupName, b))
	}
	os.Exit(runJobWrite(client, name, b))
}

func runJobWrite(c *api.Client, job string, policy []byte) int {
	if err := c.Policies().WriteJobPolicy(job, policy); err != nil {
		fmt.Println("Error writing job scaling policy:", err)
		return sysexits.Software
	}

	fmt.Println("Successfully wrote job scaling policy")
	return sysexits.OK
}

func runJobGroupWrite(c *api.Client, job, group string, policy []byte) int {
	if err := c.Policies().WriteJobGroupPolicy(job, group, policy); err != nil {
		fmt.Println("Error writing job group scaling policy:", err)
		return sysexits.Software
	}

	fmt.Println("Successfully wrote job group scaling policy")
	return sysexits.OK
}
