# Scale CLI

The scale command groups subcommands for actioning scaling requests. Users can scale in or out selected job groups by the count desired.

## Examples

Scale out job `example` and group `cache` by a count of `2`:
```bash
$ sherpa scale out --group-name=cache --count=2 example
```

Scale in job `example` and group `cache` using the count configured in the scaling policy:
```bash
$ sherpa scale in --group-name=cache example
```

## Usage
```bash
Usage:
  sherpa scale [flags]
  sherpa scale [command]

Available Commands:
  in          Perform scaling in actions on Nomad jobs and groups.
  out         Perform scaling out actions on Nomad jobs and groups.
```
