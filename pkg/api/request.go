package api

import (
	"io"
	"net/http"
	"net/url"
)

// QueryOptions are used to create a query which includes query params. This is used for GET, POST
// and PUT calls.
type QueryOptions struct {

	// Params are HTTP parameters on the query URL.
	Params map[string]string
}

type request struct {
	config *Config
	method string
	url    *url.URL
	params url.Values
	body   io.Reader
	obj    interface{}
}

// setQueryOptions is used to annotate the request with additional query options.
func (r *request) setQueryOptions(q *QueryOptions) {
	if q == nil {
		return
	}
	for k, v := range q.Params {
		r.params.Set(k, v)
	}
}

func (r *request) toHTTP() (*http.Request, error) {
	// Encode the get parameters
	r.url.RawQuery = r.params.Encode()

	// Check if we should encode the body
	if r.body == nil && r.obj != nil {
		b, err := encodeBody(r.obj)
		if err != nil {
			return nil, err
		}
		r.body = b
	}

	// Create the HTTP request
	req, err := http.NewRequest(r.method, r.url.RequestURI(), r.body)
	if err != nil {
		return nil, err
	}

	req.URL.Host = r.url.Host
	req.URL.Scheme = r.url.Scheme
	req.Host = r.url.Host
	return req, nil
}
