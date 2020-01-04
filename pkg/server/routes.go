package server

import (
	"net/http"
	"net/http/pprof"

	policyV1 "github.com/jrasell/sherpa/pkg/policy/v1"
	scaleV1 "github.com/jrasell/sherpa/pkg/scale/v1"
	v1 "github.com/jrasell/sherpa/pkg/server/endpoints/v1"
	"github.com/jrasell/sherpa/pkg/server/router"
)

type routes struct {
	System *v1.SystemServer
	Policy *policyV1.Policy
	Scale  *scaleV1.Scale
	UI     *v1.UIServer
}

func (h *HTTPServer) setupRoutes() *router.RouteTable {
	h.logger.Debug().Msg("setting up HTTP server routes")

	var r router.RouteTable

	// Setup the scaling routes.
	scaleRoutes := h.setupScaleRoutes()
	r = append(r, scaleRoutes)

	// Setup the system routes.
	systemRoutes := h.setupSystemRoutes()
	r = append(r, systemRoutes)

	// Setup the base policy routes.
	policyRoutes := h.setupPolicyRoutes()
	r = append(r, policyRoutes)

	// Setup the server debug routes if enabled.
	if h.cfg.Debug {
		debugRoutes := h.setupDebugRoutes()
		r = append(r, debugRoutes)
	}

	// Setup the UI routes if it is enabled.
	if h.cfg.Server.UI {
		uiRoutes := h.setupUIRoutes()
		r = append(r, uiRoutes)
	}

	// Setup the policy engine API route if it is enabled.
	if h.cfg.Server.APIPolicyEngine {
		apiPolicyRoutes := h.setupAPIPolicyRoutes()
		r = append(r, apiPolicyRoutes)
	}

	return &r
}

func (h *HTTPServer) setupUIRoutes() []router.Route {
	h.logger.Debug().Msg("setting up server UI routes")

	h.routes.UI = v1.NewUIServer()

	return router.Routes{
		router.Route{
			Name:        routeUIName,
			Method:      http.MethodGet,
			Pattern:     routeUIPattern,
			HandlerFunc: h.routes.UI.Get,
		},
		router.Route{
			Name:        routeUIRedirectName,
			Method:      http.MethodGet,
			Pattern:     routeUIRedirectPattern,
			HandlerFunc: h.routes.UI.Redirect,
		},
	}
}

func (h *HTTPServer) setupScaleRoutes() []router.Route {
	h.logger.Debug().Msg("setting up server scale routes")

	h.routes.Scale = scaleV1.NewScaleServer(h.cfg.Server.StrictPolicyChecking, &scaleV1.ScaleConfig{
		Logger: h.logger,
		Policy: h.policyBackend,
		Scale:  h.scaleBackend,
		State:  h.stateBackend,
	})

	return router.Routes{
		// Deprecated: the PUT method is deprecated in favour of POST and will be removed in a
		// future release.
		router.Route{
			Name:    routeScaleOutJobGroupName,
			Method:  http.MethodPut,
			Pattern: routeScaleOutJobGroupPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Scale.OutJobGroup),
		},
		// Deprecated: the PUT method is deprecated in favour of POST and will be removed in a
		// future release.
		router.Route{
			Name:    routeScaleInJobGroupName,
			Method:  http.MethodPut,
			Pattern: routeScaleInJobGroupPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Scale.InJobGroup),
		},
		router.Route{
			Name:    routePostScaleOutJobGroupName,
			Method:  http.MethodPost,
			Pattern: routePostScaleOutJobGroupPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Scale.OutJobGroup),
		},
		router.Route{
			Name:    routePostScaleInJobGroupName,
			Method:  http.MethodPost,
			Pattern: routePostScaleInJobGroupPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Scale.InJobGroup),
		},
		router.Route{
			Name:    routeGetScalingStatusName,
			Method:  http.MethodGet,
			Pattern: routeGetScalingStatusPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Scale.StatusList),
		},
		router.Route{
			Name:    routeGetScalingInfoName,
			Method:  http.MethodGet,
			Pattern: routeGetScalingInfoPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Scale.StatusInfo),
		},
	}
}

func (h *HTTPServer) setupSystemRoutes() []router.Route {
	h.logger.Debug().Msg("setting up server system routes")

	h.routes.System = v1.NewSystemServer(h.logger, h.nomad, h.cfg.Server, h.telemetry, h.clusterMember)

	return router.Routes{
		router.Route{
			Name:        routeSystemHealthName,
			Method:      http.MethodGet,
			Pattern:     routeSystemHealthPattern,
			HandlerFunc: h.routes.System.GetHealth,
		},
		router.Route{
			Name:        routeSystemInfoName,
			Method:      http.MethodGet,
			Pattern:     routeSystemInfoPattern,
			HandlerFunc: h.routes.System.GetInfo,
		},
		router.Route{
			Name:        routeGetMetricsName,
			Method:      http.MethodGet,
			Pattern:     routeGetMetricsPattern,
			HandlerFunc: h.routes.System.GetMetrics,
		},
		router.Route{
			Name:        routeGetSystemLeaderName,
			Method:      http.MethodGet,
			Pattern:     routeGetSystemLeaderPattern,
			HandlerFunc: h.routes.System.GetLeader,
		},
	}
}

func (h *HTTPServer) setupPolicyRoutes() []router.Route {
	h.logger.Debug().Msg("setting up server policy routes")

	h.routes.Policy = policyV1.NewPolicyServer(h.logger, h.policyBackend)

	return router.Routes{
		router.Route{
			Name:    routeGetJobScalingPoliciesName,
			Method:  http.MethodGet,
			Pattern: routeGetJobScalingPoliciesPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Policy.GetJobPolicies),
		},
		router.Route{
			Name:    routeGetJobScalingPolicyName,
			Method:  http.MethodGet,
			Pattern: routeGetJobScalingPolicyPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Policy.GetJobPolicy),
		},
		router.Route{
			Name:    routeGetJobGroupScalingPolicyName,
			Method:  http.MethodGet,
			Pattern: routeGetJobGroupScalingPolicyPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Policy.GetJobGroupPolicy),
		},
	}
}

func (h *HTTPServer) setupAPIPolicyRoutes() []router.Route {
	h.logger.Debug().Msg("setting up server API policy engine routes")

	return router.Routes{
		router.Route{
			Name:    routePostJobScalingPolicyName,
			Method:  http.MethodPost,
			Pattern: routePutJobScalingPolicyPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Policy.PutJobPolicy),
		},
		router.Route{
			Name:    routePostJobGroupScalingPolicyName,
			Method:  http.MethodPost,
			Pattern: routePutJobGroupScalingPolicyPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Policy.PutJobGroupPolicy),
		},
		router.Route{
			Name:    routeDeleteJobGroupScalingPolicyName,
			Method:  http.MethodDelete,
			Pattern: routeDeleteJobGroupScalingPolicyPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Policy.DeleteJobGroupPolicy),
		},
		router.Route{
			Name:    routeDeleteJobScalingPolicyName,
			Method:  http.MethodDelete,
			Pattern: routeDeleteJobScalingPolicyPattern,
			Handler: leaderProtectedHandler(h.clusterMember, h.routes.Policy.DeleteJobPolicy),
		},
	}
}

func (h *HTTPServer) setupDebugRoutes() []router.Route {
	h.logger.Debug().Msg("setting up server Debug routes")

	return router.Routes{
		router.Route{
			Name:        routeGetDebugPPROFName,
			Method:      http.MethodGet,
			Pattern:     routeGetDebugPPROFPattern,
			HandlerFunc: pprof.Index,
		},
		router.Route{
			Name:        routeGetDebugPPROFCMDLineName,
			Method:      http.MethodGet,
			Pattern:     routeGetDebugPPROFCMDLinePattern,
			HandlerFunc: pprof.Cmdline,
		},
		router.Route{
			Name:        routeGetDebugPPROFProfileName,
			Method:      http.MethodGet,
			Pattern:     routeGetDebugPPROFProfilePattern,
			HandlerFunc: pprof.Profile,
		},
		router.Route{
			Name:        routeGetDebugPPROFSymbolName,
			Method:      http.MethodGet,
			Pattern:     routeGetDebugPPROFSymbolPattern,
			HandlerFunc: pprof.Symbol,
		},
		router.Route{
			Name:        routeGetDebugPPROFTraceName,
			Method:      http.MethodGet,
			Pattern:     routeGetDebugPPROFTracePattern,
			HandlerFunc: pprof.Trace,
		},
	}
}
