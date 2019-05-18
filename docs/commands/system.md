# System CLI

The system command groups subcommands gaining insights and information about a running Sherpa server.

## Examples

Detail the health of the server:
```bash
$ sherpa system health
```

Show configuration details of the running server:
```bash
$ sherpa system info
```

Output the latest metric data points for the running server:
```bash
$ sherpa system metrics
```

## Usage
```bash
Usage:
  sherpa system [flags]
  sherpa system [command]

Available Commands:
  health      Retrieve health information of a Sherpa server
  info        Retrieve information about a Sherpa server
  metrics     Retrieve metrics from a Sherpa server
```
