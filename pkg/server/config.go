package server

import (
	autoscalecfg "github.com/jrasell/sherpa/pkg/config/autoscale"
	serverCfg "github.com/jrasell/sherpa/pkg/config/server"
)

type Config struct {
	Server    *serverCfg.Config
	AutoScale *autoscalecfg.Config
}

const (
	routeSystemHealthName                   = "GetSystemHealth"
	routeSystemHealthPattern                = "/v1/system/health"
	routeSystemInfoName                     = "GetSystemInfo"
	routeSystemInfoPattern                  = "/v1/system/info"
	routeScaleOutJobGroupName               = "ScaleOutJobGroup"
	routeScaleOutJobGroupPattern            = "/v1/scale/out/{job_id}/{group}"
	routeScaleInJobGroupName                = "ScaleInJobGroup"
	routeScaleInJobGroupPattern             = "/v1/scale/in/{job_id}/{group}"
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
)
