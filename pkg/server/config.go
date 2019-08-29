package server

import (
	serverCfg "github.com/jrasell/sherpa/pkg/config/server"
)

type Config struct {
	Server    *serverCfg.Config
	TLS       *serverCfg.TLSConfig
	Telemetry *serverCfg.TelemetryConfig
}

const (
	routeSystemHealthName          = "GetSystemHealth"
	routeSystemHealthPattern       = "/v1/system/health"
	routeSystemInfoName            = "GetSystemInfo"
	routeSystemInfoPattern         = "/v1/system/info"
	routeScaleOutJobGroupName      = "ScaleOutJobGroup"
	routeScaleOutJobGroupPattern   = "/v1/scale/out/{job_id}/{group}"
	routeScaleInJobGroupName       = "ScaleInJobGroup"
	routeScaleInJobGroupPattern    = "/v1/scale/in/{job_id}/{group}"
	routeGetJobScalingPoliciesName = "GetJobScalingPolicies"

	routeGetScalingStatusPattern = "/v1/scale/status"
	routeGetScalingStatusName    = "GetScalingStatus"

	routeGetScalingInfoPattern = "/v1/scale/status/{id}"
	routeGetScalingInfoName    = "GetScalingInfo"

	routeGetJobScalingPoliciesPattern       = "/v1/policies"
	routeGetJobScalingPolicyName            = "GetJobScalingPolicy"
	routeGetJobScalingPolicyPattern         = "/v1/policy/{job_id}"
	routeGetJobGroupScalingPolicyName       = "GetJobGroupScalingPolicy"
	routeGetJobGroupScalingPolicyPattern    = "/v1/policy/{job_id}/{group}"
	routePostJobScalingPolicyName           = "PostJobScalingPolicy"
	routePutJobScalingPolicyPattern         = "/v1/policy/{job_id}"
	routePostJobGroupScalingPolicyName      = "PostJobGroupScalingPolicy"
	routePutJobGroupScalingPolicyPattern    = "/v1/policy/{job_id}/{group}"
	routeDeleteJobGroupScalingPolicyName    = "DeleteJobGroupScalingPolicy"
	routeDeleteJobGroupScalingPolicyPattern = "/v1/policy/{job_id}/{group}"
	routeDeleteJobScalingPolicyName         = "DeleteJobScalingPolicy"
	routeDeleteJobScalingPolicyPattern      = "/v1/policy/{job_id}"
	routeGetMetricsName                     = "GetSystemMetrics"
	routeGetMetricsPattern                  = "/v1/system/metrics"

	telemetryInterval = 10
)
