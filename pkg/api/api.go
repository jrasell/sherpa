package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-rootcerts"
	clientCfg "github.com/jrasell/sherpa/pkg/config/client"
)

type Config struct {
	Address    string
	TLSConfig  *TLSConfig
	httpClient *http.Client
}

type TLSConfig struct {
	CACert        string
	ClientCert    string
	ClientCertKey string
}

type Client struct {
	config Config
}

func DefaultConfig(cfg *clientCfg.Config) *Config {
	config := Config{
		Address:    "http://127.0.0.1:8000",
		TLSConfig:  &TLSConfig{},
		httpClient: cleanhttp.DefaultClient(),
	}
	transport := config.httpClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if cfg.Addr != "" {
		config.Address = cfg.Addr
	}
	if cfg.CAPath != "" {
		config.TLSConfig.CACert = cfg.CAPath
	}
	if cfg.CertPath != "" {
		config.TLSConfig.ClientCert = cfg.CertPath
	}
	if cfg.CertKeyPath != "" {
		config.TLSConfig.ClientCertKey = cfg.CertKeyPath
	}
	return &config
}

func NewClient(config *Config) (*Client, error) {
	if _, err := url.Parse(config.Address); err != nil {
		return nil, fmt.Errorf("invalid address '%s': %v", config.Address, err)
	}

	// Configure the TLS configurations
	if err := config.ConfigureTLS(); err != nil {
		return nil, err
	}

	return &Client{
		config: *config,
	}, nil
}

func (c *Config) ConfigureTLS() error {
	if c.TLSConfig == nil {
		return nil
	}

	var clientCert tls.Certificate

	foundClientCert := false
	if c.TLSConfig.ClientCert != "" || c.TLSConfig.ClientCertKey != "" {
		if c.TLSConfig.ClientCert != "" && c.TLSConfig.ClientCertKey != "" {
			var err error
			clientCert, err = tls.LoadX509KeyPair(c.TLSConfig.ClientCert, c.TLSConfig.ClientCertKey)
			if err != nil {
				return err
			}
			foundClientCert = true
		} else {
			return fmt.Errorf("client cert and client key must be provided")
		}
	}

	clientTLSConfig := c.httpClient.Transport.(*http.Transport).TLSClientConfig
	rootConfig := &rootcerts.Config{CAPath: c.TLSConfig.CACert}

	if err := rootcerts.ConfigureTLS(clientTLSConfig, rootConfig); err != nil {
		return err
	}

	if foundClientCert {
		clientTLSConfig.Certificates = []tls.Certificate{clientCert}
	}
	return nil
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
