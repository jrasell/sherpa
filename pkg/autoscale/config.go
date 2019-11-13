package autoscale

import (
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/sherpa/pkg/config/server"
	policyBackend "github.com/jrasell/sherpa/pkg/policy/backend"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/rs/zerolog"
)

type SetupConfig struct {
	ScalingInterval   int
	ScalingThreads    int
	StrictChecking    bool
	MetricProviderCfg *server.MetricProviderConfig

	Logger        zerolog.Logger
	PolicyBackend policyBackend.PolicyBackend
	Scale         scale.Scale
	Nomad         *api.Client
}

type Config struct {
	ScalingInterval   int
	ScalingThreads    int
	StrictChecking    bool
	MetricProviderCfg *server.MetricProviderConfig
}
