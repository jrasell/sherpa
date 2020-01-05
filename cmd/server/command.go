package server

import (
	"fmt"
	"os"

	logCfg "github.com/jrasell/sherpa/pkg/config/log"
	serverCfg "github.com/jrasell/sherpa/pkg/config/server"
	"github.com/jrasell/sherpa/pkg/logger"
	"github.com/jrasell/sherpa/pkg/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/sysexits"
	"github.com/spf13/cobra"
)

func RegisterCommand(rootCmd *cobra.Command) error {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start a Sherpa server",
		Run: func(cmd *cobra.Command, args []string) {
			runServer(cmd, args)
		},
	}
	serverCfg.RegisterConfig(cmd)
	serverCfg.RegisterTLSConfig(cmd)
	serverCfg.RegisterTelemetryConfig(cmd)
	serverCfg.RegisterClusterConfig(cmd)
	serverCfg.RegisterMetricProviderConfig(cmd)
	serverCfg.RegisterDebugConfig(cmd)
	logCfg.RegisterConfig(cmd)
	rootCmd.AddCommand(cmd)

	return nil
}

func runServer(_ *cobra.Command, _ []string) {
	serverConfig := serverCfg.GetConfig()
	tlsConfig := serverCfg.GetTLSConfig()
	telemetryConfig := serverCfg.GetTelemetryConfig()
	clusterConfig := serverCfg.GetClusterConfig()
	metricProviderConfig := serverCfg.GetMetricProviderConfig()

	if err := verifyServerConfig(serverConfig); err != nil {
		fmt.Println(err)
		os.Exit(sysexits.Usage)
	}

	// Setup the server logging.
	logConfig := logCfg.GetConfig()
	if err := logger.Setup(logConfig); err != nil {
		fmt.Println(err)
		os.Exit(sysexits.Software)
	}

	cfg := &server.Config{
		Debug:          serverCfg.GetDebugEnabled(),
		Cluster:        &clusterConfig,
		MetricProvider: metricProviderConfig,
		Server:         &serverConfig,
		TLS:            &tlsConfig,
		Telemetry:      &telemetryConfig,
	}
	srv := server.New(log.Logger, cfg)

	if err := srv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(sysexits.Software)
	}
}

func verifyServerConfig(cfg serverCfg.Config) error {
	if cfg.NomadMetaPolicyEngine && cfg.APIPolicyEngine {
		return errors.New("Please only enable one policy engine")
	}
	return nil
}
