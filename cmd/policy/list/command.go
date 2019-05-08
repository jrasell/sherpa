package list

import (
	"fmt"
	"os"

	"github.com/jrasell/sherpa/pkg/api"
	"github.com/ryanuber/columnize"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

const (
	outputHeader = "Job:Group|Enabled|MinCount|MaxCount|ScaleInCount|ScaleOutCount"
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

	clientCfg := api.DefaultConfig()

	client, err := api.NewClient(clientCfg)
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

	var out []string
	out = append(out, outputHeader)

	for job, v := range *resp {
		for group, pol := range v {
			out = append(out, fmt.Sprintf("%s:%s|%v|%v|%v|%v|%v",
				job, group, pol.Enabled, pol.MinCount, pol.MaxCount, pol.ScaleInCount, pol.ScaleOutCount))
		}
	}
	fmt.Println(formatList(out))
}

func formatList(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	return columnize.Format(in, columnConf)
}
