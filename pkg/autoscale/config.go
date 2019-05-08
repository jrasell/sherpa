package autoscale

import "github.com/jrasell/sherpa/pkg/config/autoscale"

type Config struct {
	Scaling        *autoscale.Config
	StrictChecking bool
}

const (
	defaultCPUPercentageScaleOutThreshold    = 80
	defaultMemoryPercentageScaleOutThreshold = 80
	defaultCPUPercentageScaleInThreshold     = 20
	defaultMemoryPercentageScaleInThreshold  = 20
)
