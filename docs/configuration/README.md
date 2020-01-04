# Sherpa Configuration

The Sherpa server can be configured by supplying either CLI flags or using environment variables.

## Parameters

* `--autoscaler-enabled` (bool: false) - Enable the internal autoscaling engine.
* `--autoscaler-evaluation-interval` (int: 60) - The time period in seconds between autoscaling evaluation runs.
* `--autoscaler-num-threads` (int: 3) - Specifies the number of parallel autoscaler threads to run.
* `--bind-addr` (string: "127.0.0.1") - The HTTP server address to bind to.
* `--bind-port` (uint16: 8000) - The HTTP server port to bind to.
* `--cluster-advertise-addr` (string: "http://127.0.0.1:8000") - The Sherpa server advertise address used for NAT traversal on HTTP redirects.
* `--cluster-name` (string: "") - Specifies the identifier for the Sherpa cluster.
* `--debug-enabled` (bool: false) - Specifies if the debugging HTTP endpoints should be enabled.
* `--log-format` (string: "auto") - Specify the log format ("auto", "zerolog" or "human").
* `--log-level` (string: "info") - Change the level used for logging.
* `--log-use-color` (bool: true) - Use ANSI colors in logging output.
* `--metric-provider-prometheus-addr` (string: "") The address of the Prometheus endpoint in the form <protocol>://<addr>:<port>.
* `--policy-engine-api-enabled` (bool: true) - Enable the Sherpa API to manage scaling policies.
* `--policy-engine-nomad-meta-enabled` (bool: false) - Enable Nomad job meta lookups to manage scaling policies.
* `--policy-engine-strict-checking-enabled` (bool: true) - When enabled, all scaling activities must pass through policy checks.
* `--storage-consul-enabled` (bool: false) - Use Consul as the storage backend for state.
* `--storage-consul-path` (string: "sherpa/") - The Consul KV path that will be used to store policies and state.
* `--telemetry-prometheus` (bool: false) - Specifies whether Prometheus formatted metrics are available.
* `--telemetry-statsd-address` (string: "") - Specifies the address of a statsd server to forward metrics to.
* `--telemetry-statsite-address` (string: "") - Specifies the address of a statsite server to forward metrics data to.
* `--tls-cert-key-path` (string: "") - Path to the TLS certificate key for the Sherpa server.
* `--tls-cert-path` (string: "") - Path to the TLS certificate for the Sherpa server.
* `--ui` (bool: false) - Run the Sherpa user interface.

### Environment Variables

When specifying environment variables, the CLI flag should be converted like follows:
* `--bind-addr` becomes `SHERPA_BIND_ADDR`

## Client Parameters

Nomad and Consul clients can be configured using the native environment variables which are available through the HashiCorp SDKs. Using these keeps the setup simple and consistent.

### Nomad Client Parameters

The Nomad client environment variables documentation can be found on the [Nomad general options](https://github.com/hashicorp/nomad/blob/22fd62753510a4a41c1b8f1d117ea1a90b48df06/website/source/docs/commands/_general_options.html.md) GitHub document. For ease of use this document is reproduced below:

* `NOMAD_ADDR` (string: "http://127.0.0.1:4646") - The address of the Nomad server.
* `NOMAD_REGION` (string: "") - The region of the Nomad servers to forward commands to.
* `NOMAD_NAMESPACE` (string "default") - The target namespace for queries and actions bound to a namespace.
* `NOMAD_CACERT` (string: "") - Path to a PEM encoded CA cert file to use to verify the Nomad server SSL certificate.
* `NOMAD_CAPATH` (string: "") - Path to a directory of PEM encoded CA cert files to verify the Nomad server SSL certificate.
* `NOMAD_CLIENT_CERT` (string: "") - Path to a PEM encoded client certificate for TLS authentication to the Nomad server.
* `NOMAD_CLIENT_KEY` (string: "") - Path to an unencrypted PEM encoded private key matching the client certificate.
* `NOMAD_SKIP_VERIFY` (bool: false) - Do not verify TLS certificate.
* `NOMAD_TOKEN` (string: "") - The SecretID of an ACL token to use to authenticate API requests with.

### Consul Client Parameters

The Consul client environment variables documentation can be found on the [Consul commands page](https://www.consul.io/docs/commands/index.html#environment-variables). For ease of use this document is reproduced below:

* `CONSUL_HTTP_ADDR` (string: "127.0.0.1:8500") - This is the HTTP API address to the local Consul agent (not the remote server) specified as a URI with optional scheme.
* `CONSUL_HTTP_TOKEN` (string: "") - This is the API access token required when access control lists (ACLs) are enabled.
* `CONSUL_HTTP_AUTH` (string: "") - This specifies HTTP Basic access credentials as a username:password pair.
* `CONSUL_HTTP_SSL` (bool: false) - This is a boolean value that enables the HTTPS URI scheme and SSL connections to the HTTP API.
* `CONSUL_HTTP_SSL_VERIFY` (bool: true) - This is a boolean value to specify SSL certificate verification.
* `CONSUL_CACERT` (string: "") - Path to a CA file to use for TLS when communicating with Consul.
* `CONSUL_CAPATH` (string: "") - Path to a directory of CA certificates to use for TLS when communicating with Consul.
* `CONSUL_CLIENT_CERT` (string: "") - Path to a client cert file to use for TLS.
* `CONSUL_CLIENT_KEY` (string: "") - Path to a client key file to use for TLS.
* `CONSUL_TLS_SERVER_NAME` (string: "") - The server name to use as the SNI host when connecting via TLS.
