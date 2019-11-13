package metrics

// Provider is the interface which all external metric providers must implement.
type Provider interface {

	// GetValue takes a query string and returns the resulting metric value as a float64 along with
	// an error if one was encountered. When implementing this interface function, it should handle
	// sending any Sherpa telemetry which directly reference to implementation name.
	GetValue(query string) (*float64, error)
}
