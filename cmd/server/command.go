package server

import (
	"fmt"
	"os"

	autoscaleCfg "github.com/jrasell/sherpa/pkg/config/autoscale"
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
	autoscaleCfg.RegisterConfig(cmd)
	logCfg.RegisterConfig(cmd)
	rootCmd.AddCommand(cmd)

	return nil
}

func runServer(_ *cobra.Command, _ []string) {
	serverConfig := serverCfg.GetConfig()
	autoscaleConfig := autoscaleCfg.GetConfig()

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

	cfg := &server.Config{Server: &serverConfig, AutoScale: &autoscaleConfig}

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

	if cfg.NomadMetaPolicyEngine && cfg.ConsulStorageBackend {
		return errors.New("Consul storage backend is not compatible with Nomad meta policy engine")
	}

	return nil
}
