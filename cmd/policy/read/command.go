package read

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jrasell/sherpa/cmd/helper"
	"github.com/jrasell/sherpa/pkg/api"
	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
	"github.com/liamg/tml"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

const (
	nomadCheckHeader    = "CPU In|CPU Out|Memory In|Memory Out"
	externalCheckHeader = "Name|Enabled|Provider|Operator|Value|Action|Query"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Details scaling policies associated to a job",
		Run: func(cmd *cobra.Command, args []string) {
			runRead(cmd, args)
		},
	}

	rootCmd.AddCommand(cmd)

	return nil
}

func runRead(_ *cobra.Command, args []string) {
	switch {
	case len(args) < 1:
		fmt.Println("Not enough arguments, expected 1 got", len(args))
		os.Exit(sysexits.Usage)
	case len(args) > 1:
		fmt.Println("Too many arguments, expected 1 got", len(args))
		os.Exit(sysexits.Usage)
	}

	clientConfig := clientCfg.GetConfig()
	mergedConfig := api.DefaultConfig(&clientConfig)

	client, err := api.NewClient(mergedConfig)
	if err != nil {
		fmt.Println("Error setting up Sherpa client:", err)
		os.Exit(sysexits.Software)
	}

	job := strings.ToLower(strings.TrimSpace(args[0]))

	resp, err := client.Policies().ReadJobPolicy(job)
	if err != nil {
		fmt.Println("Error reading scaling policy:", err)
		os.Exit(sysexits.Software)
	}

	if len(*resp) == 0 {
		os.Exit(sysexits.OK)
	}

	// Sort the keys so the output is ordered alphabetically by group name.
	keys := []string{}
	for k := range *resp {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		formatGroupOutput(k, (*resp)[k])
	}
}

// formatGroupOutput format all the details of a group scaling policy and outputs this to the user
// in a readable, easy to follow fashion.
func formatGroupOutput(group string, policy *api.JobGroupPolicy) {
	header := []string{
		fmt.Sprintf("Group|%s", group),
		fmt.Sprintf("MinCount|%v", policy.MinCount),
		fmt.Sprintf("MaxCount|%v", policy.MaxCount),
		fmt.Sprintf("Cooldown|%v", policy.Cooldown),
		fmt.Sprintf("ScaleInCount|%v", policy.ScaleInCount),
		fmt.Sprintf("ScaleOutCount|%v", policy.ScaleOutCount),
	}

	var nomadChecks []string
	var externalChecks []string

	// Check we have Nomad checks configured.
	if policy.ScaleInMemoryPercentageThreshold != nil || policy.ScaleOutMemoryPercentageThreshold != nil ||
		policy.ScaleInCPUPercentageThreshold != nil || policy.ScaleOutCPUPercentageThreshold != nil {

		// Create the Nomad check output.
		nomadChecks = append(nomadChecks, nomadCheckHeader)
		nomadChecks = append(nomadChecks, fmt.Sprintf("%v%%|%v%%|%v%%|%v%%",
			*policy.ScaleInCPUPercentageThreshold, *policy.ScaleOutCPUPercentageThreshold,
			*policy.ScaleInMemoryPercentageThreshold, *policy.ScaleOutMemoryPercentageThreshold))
	}

	// Check is there are external checks configured.
	if policy.ExternalChecks != nil {

		// Set the header.
		externalChecks = append(externalChecks, externalCheckHeader)

		// Iterate the checks and format the output.
		for name, check := range policy.ExternalChecks {
			externalChecks = append(externalChecks, fmt.Sprintf("%s|%v|%s|%s|%v|%s|%s",
				name, check.Enabled, check.Provider, check.ComparisonOperator, check.ComparisonValue, check.Action, check.Query))
		}
	}

	// Print our top header and include the core required parameters of a group scaling policy.
	tml.Println("<bold>Scaling Policy:</bold>")
	fmt.Println(helper.FormatKV(header))
	fmt.Println("")

	if len(nomadChecks) > 0 {
		fmt.Println("Nomad Checks:")
		fmt.Println(helper.FormatList(nomadChecks))
		fmt.Println("")
	}

	if len(externalChecks) > 0 {
		fmt.Println("External Checks:")
		fmt.Println(helper.FormatList(externalChecks))
		fmt.Println("")
	}
}
