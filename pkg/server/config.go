package server

import (
	serverCfg "github.com/jrasell/sherpa/pkg/config/server"
)

type Config struct {
	Debug          bool
	Cluster        *serverCfg.ClusterConfig
	MetricProvider *serverCfg.MetricProviderConfig
	Server         *serverCfg.Config
	TLS            *serverCfg.TLSConfig
	Telemetry      *serverCfg.TelemetryConfig
}

const (
	routeUIRedirectName    = "UIRedirect"
	routeUIRedirectPattern = "/"
	routeUIName            = "UI"
	routeUIPattern         = "/ui"

	routeGetScalingStatusPattern            = "/v1/scale/status"
	routeGetScalingStatusName               = "GetScalingStatus"
	routeGetScalingInfoPattern              = "/v1/scale/status/{id}"
	routeGetScalingInfoName                 = "GetScalingInfo"
	routeScaleOutJobGroupName               = "ScaleOutJobGroup"
	routeScaleOutJobGroupPattern            = "/v1/scale/out/{job_id}/{group}"
	routeScaleInJobGroupName                = "ScaleInJobGroup"
	routeScaleInJobGroupPattern             = "/v1/scale/in/{job_id}/{group}"
	routePostScaleOutJobGroupName           = "ScaleOutJobGroup"
	routePostScaleOutJobGroupPattern        = "/v1/scale/out/{job_id}/{group}"
	routePostScaleInJobGroupName            = "ScaleInJobGroup"
	routePostScaleInJobGroupPattern         = "/v1/scale/in/{job_id}/{group}"
	routeGetJobScalingPoliciesName          = "GetJobScalingPolicies"
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

// System server routes.
const (
	routeGetSystemLeaderName    = "GetSystemLeader"
	routeGetSystemLeaderPattern = "/v1/system/leader"
	routeSystemHealthName       = "GetSystemHealth"
	routeSystemHealthPattern    = "/v1/system/health"
	routeSystemInfoName         = "GetSystemInfo"
	routeSystemInfoPattern      = "/v1/system/info"
)

// Debug server routes.
const (
	routeGetDebugPPROFName           = "GetDebugPPROF"
	routeGetDebugPPROFPattern        = "/debug/pprof/"
	routeGetDebugPPROFCMDLineName    = "GetDebugPPROFCMDLine"
	routeGetDebugPPROFCMDLinePattern = "/debug/pprof/cmdline"
	routeGetDebugPPROFProfileName    = "GetDebugPPROFProfile"
	routeGetDebugPPROFProfilePattern = "/debug/pprof/profile"
	routeGetDebugPPROFSymbolName     = "GetDebugPPROFSymbol"
	routeGetDebugPPROFSymbolPattern  = "/debug/pprof/symbol"
	routeGetDebugPPROFTraceName      = "GetDebugPPROFTrace"
	routeGetDebugPPROFTracePattern   = "/debug/pprof/trace"
)
