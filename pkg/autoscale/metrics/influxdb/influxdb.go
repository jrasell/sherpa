package influxdb

import (
	"encoding/json"
	"time"

	sendMetrics "github.com/armon/go-metrics"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/jrasell/sherpa/pkg/autoscale/metrics"
	"github.com/jrasell/sherpa/pkg/helper"
	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Client is an InfluxDB metrics backend wrapper.
type Client struct {
	logger         zerolog.Logger
	influxDBClient client.Client
}

// NewClient takes the InfluDB connection config and builds the client for use in retrieving
// metric values.
func NewClient(addr string, user string, pass string, ins bool, log zerolog.Logger) (metrics.Provider, error) {
	client, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               addr,
		Username:           user,
		Password:           pass,
		InsecureSkipVerify: ins,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		logger:         log.With().Str("metric-provider", policy.ProviderInfluxDB.String()).Logger(),
		influxDBClient: client,
	}, nil
}

// GetValue satisfies the GetValue function of the metrics.Provider interface.
func (c *Client) GetValue(query string) (*float64, error) {
	defer sendMetrics.MeasureSince([]string{"autoscale", "influxdb", "get_value"}, time.Now())

	// Gather the value and any error returned from attempting to call InfluxDB; handling the
	// error result via Sherpa telemetry.
	value, err := c.getValue(query)
	if err != nil {
		sendMetrics.IncrCounter([]string{"autoscale", "influxdb", "error"}, 1)
	} else {
		sendMetrics.IncrCounter([]string{"autoscale", "influxdb", "success"}, 1)
	}
	return value, err
}

// getValue performs the InfluxDB query work, allowing the interface implementation to handle end
// state activities. The data returned from InfluxDB comes back as a series, with columns and rows.
// the first column is always "time", followed by the actual queried value.
func (c *Client) getValue(query string) (*float64, error) {

	queryInfluxDB := client.NewQuery(query, "", "s")
	resp, err := c.influxDBClient.Query(queryInfluxDB)
	if err != nil {
		return nil, err
	}

	for _, result := range resp.Results {
		// The query should only return 1 row
		if len(result.Series) == 1 {
			for _, item := range result.Series {
				for _, row := range item.Values {
					// example timestamp, cpu -> Values:[[1582671262 11.690490326184905]]
					metric := row[1].(json.Number)
					floatVal, err := metric.Float64()
					if err != nil {
						return nil, errors.Wrap(err, "Metric value could not be converted to float64")
					}
					return helper.Float64ToPointer(floatVal), nil
				}
			}
		}
		return nil, errors.Wrap(err, "Query returned incorrect length result list from InfluxDB")
	}
	return nil, errors.Wrap(err, "Query did not return valid data from InfluxDB")
}
