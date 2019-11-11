package prometheus

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	sendMetrics "github.com/armon/go-metrics"
	"github.com/jrasell/sherpa/pkg/autoscale/metrics"
	"github.com/jrasell/sherpa/pkg/helper"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	"github.com/rs/zerolog"
)

const (
	// queryEndpoint is the Prometheus API endpoint used and currently supported for querying
	// metric values.
	queryEndpoint = "/api/v1/query?query="
)

type queryResp struct {
	Status string        `json:"status"`
	Data   queryRespData `json:"data"`
}

type queryRespData struct {
	ResultType string            `json:"resultType"`
	Result     []queryRespResult `json:"result"`
}

type queryRespResult struct {
	Value []interface{} `json:"value"`
}

// Client is a Prometheus metrics backend wrapper.
type Client struct {
	logger           zerolog.Logger
	prometheusClient api.Client
	queryAddr        string
}

// NewClient takes the base Prometheus API address and build the client for use in retrieving
// metric values.
func NewClient(addr string, log zerolog.Logger) (metrics.Provider, error) {
	client, err := api.NewClient(api.Config{Address: addr})
	if err != nil {
		return nil, err
	}
	return &Client{
		logger:           log.With().Str("metric-provider", policy.ProviderPrometheus.String()).Logger(),
		prometheusClient: client,
		queryAddr:        addr + queryEndpoint,
	}, nil
}

// GetValue satisfies the GetValue function of the metrics.Provider interface.
func (c *Client) GetValue(query string) (*float64, error) {
	defer sendMetrics.MeasureSince([]string{"autoscale", "prometheus", "get_value"}, time.Now())

	// Gather the value and any error returned from attempting to call Prometheus; handling the
	// error result via Sherpa telemetry.
	value, err := c.getValue(query)
	if err != nil {
		sendMetrics.IncrCounter([]string{"autoscale", "prometheus", "error"}, 1)
	} else {
		sendMetrics.IncrCounter([]string{"autoscale", "prometheus", "success"}, 1)
	}
	return value, err
}

// getValue performs the Prometheus query work, allowing the interface implementation to handle end
// state activities.
func (c *Client) getValue(query string) (*float64, error) {
	ctx := context.Background()

	parsedURL, err := url.Parse(c.queryAddr + url.QueryEscape(query))
	if err != nil {
		return nil, err
	}
	c.logger.Debug().Str("url", parsedURL.String()).Msg("successfully built query URL")

	req, err := http.NewRequest(http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	_, bytes, err := c.prometheusClient.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	// Unmarshal the prometheus API response and get the metric value.
	var unmarshalResp queryResp
	if err := json.Unmarshal(bytes, &unmarshalResp); err != nil {
		return nil, err
	}
	return c.getValueFromResp(&unmarshalResp)
}

// getValueFromResp is used to get the single metric value from the Prometheus response.
func (c *Client) getValueFromResp(resp *queryResp) (*float64, error) {

	// If we do not have the correct number of results, do not guess, inform the client this is an
	// error so they can fix the query.
	if len(resp.Data.Result) < 1 || len(resp.Data.Result) > 1 {
		return nil, errors.New("received incorrect length result list from Prometheus")
	}

	floatVal, err := strconv.ParseFloat(resp.Data.Result[0].Value[1].(string), 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert Prometheus metric value to float64")
	}
	return helper.Float64ToPointer(floatVal), nil
}
