# Sherpa Configuration

The Sherpa server can be configured by supplying either CLI flags or using environment variables.

## Parameters

* `--autoscaling-evaluation-interval` (int: 60) - If using the built-in autoscaler, this is the time period in seconds between evaluation runs.
* `--bind-addr` (string: "127.0.0.1") - Specifies the address which the Sherpa HTTP server will bind to for network connectivity. 
* `--bind-port` (uint16: 8000) - The address the HTTP server is bound to.
* `--api-policy-engine-enabled` (bool: true) - Enables configuring scaling policies via the CLI and API. Conflicts with the Nomad meta engine.
* `--consul-storage-backend-enabled` (bool: false) - Enables storing polices created via the API in Consul KV as a durable backend.
* `--consul-storage-backend-path` (string: "sherpa/policies/") - The Consul KV path that will be used to store policies.
* `--internal-auto-scaler-enabled` (bool: false) - Enables the Sherpa internal autoscaler, which will run periodically evaluating job resource utilization and making scaling decisions.
* `--nomad-meta-policy-engine-enabled` (bool: false) - If enabled, job scaling policies will be read and updated using job meta parameters. Conflicts with the API engine.
* `--strict-policy-checking-enabled` (bool: true) - Specifies whether Sherpa should strictly check all scaling requests against scaling policies.
* `--log-format` (string: "auto") - Is the format at which to log to and can be "auto", "zerolog" or "human".
* `--log-level` (string: "info") - The level at which Sherpa should log to. Valid log levels include WARN, INFO, or DEBUG in increasing order of verbosity.
* `--log-use-color` (bool: true) - Specifies whether to use ANSI colour logging.
* `--tls-cert-key-path` (string: "") - Path to the TLS certificate key for the Sherpa server.
* `--tls-cert-path` (string: "") - Path to the TLS certificate for the Sherpa server.

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
