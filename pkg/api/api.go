package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
)

type Config struct {
	Address    string
	httpClient *http.Client
}

// QueryOptions are used to create a query which includes query params. This is used for GET, POST
// and PUT calls.
type QueryOptions struct {

	// Params are HTTP parameters on the query URL.
	Params map[string]string
}

type Client struct {
	config Config
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

func DefaultConfig() *Config {
	return &Config{
		Address:    "http://127.0.0.1:8000",
		httpClient: cleanhttp.DefaultClient(),
	}
}

func NewClient(config *Config) (*Client, error) {

	defConfig := DefaultConfig()

	if config.Address == "" {
		config.Address = defConfig.Address
	} else if _, err := url.Parse(config.Address); err != nil {
		return nil, fmt.Errorf("invalid address '%s': %v", config.Address, err)
	}

	if config.httpClient == nil {
		config.httpClient = defConfig.httpClient
	}

	return &Client{
		config: *config,
	}, nil
}

func (c *Client) get(endpoint string, out interface{}) error {
	r, err := c.newRequest(http.MethodGet, endpoint)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(r)
	resp, err = requireOK(resp, err, http.StatusOK)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := decodeBody(&resp.Body, out); err != nil {
		return err
	}
	return nil
}

func (c *Client) put(endpoint string, in, out interface{}, q *QueryOptions) error {
	r, err := c.newRequest(http.MethodPut, endpoint)
	if err != nil {
		return err
	}

	r.setQueryOptions(q)

	r.obj = in
	resp, err := c.doRequest(r)
	resp, err = requireOK(resp, err, http.StatusOK)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := decodeBody(&resp.Body, out); err != nil {
		return err
	}

	return nil
}

func (c *Client) post(endpoint string, in, out interface{}) error {
	r, err := c.newRequest(http.MethodPost, endpoint)
	if err != nil {
		return err
	}
	r.obj = in
	resp, err := c.doRequest(r)
	resp, err = requireOK(resp, err, http.StatusCreated)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if out != nil {
		if err := decodeBody(&resp.Body, out); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) delete(endpoint string, out interface{}) error {
	r, err := c.newRequest(http.MethodDelete, endpoint)
	if err != nil {
		return err
	}

	resp, err := c.doRequest(r)
	resp, err = requireOK(resp, err, http.StatusNoContent)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if out != nil {
		if err := decodeBody(&resp.Body, &out); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) newRequest(method, path string) (*request, error) {
	base, _ := url.Parse(c.config.Address)
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	r := &request{
		config: &c.config,
		method: method,
		url: &url.URL{
			Scheme: base.Scheme,
			User:   base.User,
			Host:   base.Host,
			Path:   u.Path,
		},
		params: make(map[string][]string),
	}

	return r, nil
}

func (c *Client) doRequest(r *request) (*http.Response, error) {
	req, err := r.toHTTP()
	if err != nil {
		return nil, err
	}
	return c.config.httpClient.Do(req)
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

func decodeBody(body *io.ReadCloser, out interface{}) error {
	return json.NewDecoder(*body).Decode(out)
}

func encodeBody(obj interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}

func requireOK(resp *http.Response, e error, expected int) (*http.Response, error) {
	if e != nil {
		if resp != nil {
			resp.Body.Close() // nolint:errcheck
		}
		return nil, e
	}
	if resp.StatusCode != expected {
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body) // nolint:errcheck
		resp.Body.Close()        // nolint:errcheck
		return nil, fmt.Errorf(strings.TrimSpace(fmt.Sprintf("unexpected response code %d: %s",
			resp.StatusCode, buf.Bytes())))
	}
	return resp, nil
}
