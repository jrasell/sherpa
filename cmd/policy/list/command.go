package list

import (
	"fmt"
	"os"
	"sort"

	"github.com/jrasell/sherpa/cmd/helper"
	"github.com/jrasell/sherpa/pkg/api"
	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

const (
	outputHeader = "Job:Group|Enabled|MinCount|MaxCount|Cooldown|ScaleInCount|ScaleOutCount"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all scaling policies",
		Run: func(cmd *cobra.Command, args []string) {
			runList(cmd, args)
		},
	}
	rootCmd.AddCommand(cmd)

	return nil
}

func runList(_ *cobra.Command, args []string) {
	switch {
	case len(args) > 0:
		fmt.Println("Too many arguments, expected 0 args got", len(args))
		os.Exit(sysexits.Usage)
	}

	clientConfig := clientCfg.GetConfig()
	mergedConfig := api.DefaultConfig(&clientConfig)

	client, err := api.NewClient(mergedConfig)
	if err != nil {
		fmt.Println("Error setting up Sherpa client:", err)
		os.Exit(sysexits.Software)
	}

	resp, err := client.Policies().List()
	if err != nil {
		fmt.Println("Error querying policy list:", err)
		os.Exit(sysexits.Software)
	}

	if len(*resp) == 0 {
		os.Exit(sysexits.OK)
	}

	out := []string{outputHeader}
	out = append(out, produceSortedList(resp)...)

	fmt.Println(helper.FormatList(out))
}

func produceSortedList(input *map[string]map[string]*api.JobGroupPolicy) []string {
	var sorted []string
	for job, v := range *input {
		for group, pol := range v {
			sorted = append(sorted, fmt.Sprintf("%s:%s|%v|%v|%v|%v|%v|%v",
				job, group, pol.Enabled, pol.MinCount, pol.MaxCount, pol.Cooldown, pol.ScaleInCount, pol.ScaleOutCount))
		}
	}
	sort.Strings(sorted)
	return sorted
}
