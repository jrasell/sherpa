package read

import (
	"fmt"
	"os"
	"strings"

	"github.com/jrasell/sherpa/pkg/api"
	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
	"github.com/ryanuber/columnize"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

const (
	outputHeader = "Group|Enabled|MinCount|MaxCount|ScaleInCount|ScaleOutCount"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Details the scaling policy",
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

	out := []string{outputHeader}

	for group, pol := range *resp {
		out = append(out, fmt.Sprintf("%s|%v|%v|%v|%v|%v",
			group, pol.Enabled, pol.MinCount, pol.MaxCount, pol.ScaleInCount, pol.ScaleOutCount))
	}
	fmt.Println(formatList(out))
}

func formatList(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	return columnize.Format(in, columnConf)
}
