# Sherpa Commands (CLI)

In addition to a verbose HTTP API, Sherpa features a command-line interface that wraps common functionality and formats output.

To get help, run:
```bash
$ sherpa -h
```

To get help for a subcommand, run:
```bash
$ sherpa <subcommand> -h
```

## CLI Command Structure

There are a number of command and subcommand options available. Construct your Vault CLI command such that the command options precede its and arguments if any:

```bash
sherpa <command> <subcommand> [options] [args]
```

## General Options

The following options are available to all Sherpa CLI commands and help set client connection variables.

* `--addr` (string: "http://127.0.0.1:8000") - The HTTP(S) address of the sherpa server
* `--client-ca-path` (string: "") - Path to a PEM encoded CA cert file to use to verify the Sherpa server SSL certificate.
* `--client-cert-key-path` (string: "") - Path to an unencrypted PEM encoded private key matching the client certificate
* `--client-cert-path string` (string: "") - Path to a PEM encoded client certificate for TLS authentication to the Sherpa server

## Exit Codes

The Sherpa CLI aims to be consistent and well-behaved using [sysexits](https://github.com/sean-/sysexits) for CLI exit codes.

* exit code `64`: represents local errors such as incorrect flags, failed validation or an incorrect number of passed arguments.
* exit code `70`: represents an internal failure such as API failures.
