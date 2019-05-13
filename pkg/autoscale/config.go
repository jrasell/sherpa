package autoscale

type Config struct {
	ScalingInterval int
	StrictChecking  bool
}

const (
	defaultCPUPercentageScaleOutThreshold    = 80
	defaultMemoryPercentageScaleOutThreshold = 80
	defaultCPUPercentageScaleInThreshold     = 20
	defaultMemoryPercentageScaleInThreshold  = 20
)
