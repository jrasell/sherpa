package server

import (
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/armon/go-metrics/prometheus"
	"github.com/jrasell/sherpa/pkg/build"
)

func (h *HTTPServer) setupTelemetry() error {

	inm := metrics.NewInmemSink(telemetryInterval*time.Second, time.Minute)
	metrics.DefaultInmemSignal(inm)

	metricsConf := metrics.DefaultConfig(build.ProgramName())

	var fanout metrics.FanoutSink

	// Configure the statsite sink
	if h.cfg.Telemetry.StatsiteAddr != "" {
		sink, err := metrics.NewStatsiteSink(h.cfg.Telemetry.StatsiteAddr)
		if err != nil {
			return err
		}
		fanout = append(fanout, sink)
	}

	// Configure the statsd sink
	if h.cfg.Telemetry.StatsdAddr != "" {
		sink, err := metrics.NewStatsdSink(h.cfg.Telemetry.StatsdAddr)
		if err != nil {
			return err
		}
		fanout = append(fanout, sink)
	}

	// Configure prometheus formatted metrics.
	if h.cfg.Telemetry.Prometheus {
		promSink, err := prometheus.NewPrometheusSink()
		if err != nil {
			return err
		}
		fanout = append(fanout, promSink)
	}

	// Initialize the global sink
	if len(fanout) > 0 {
		fanout = append(fanout, inm)
		if _, err := metrics.NewGlobal(metricsConf, fanout); err != nil {
			return err
		}
	} else {
		metricsConf.EnableHostname = false
		if _, err := metrics.NewGlobal(metricsConf, inm); err != nil {
			return err
		}
	}

	h.telemetry = inm
	return nil
}
