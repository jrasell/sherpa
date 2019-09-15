# Scale CLI

The scale command groups subcommands for actioning and detailing scaling requests. Users can scale in or out selected job groups by the count desired or view past events.

## Examples

Scale out job `example` and group `cache` by a count of `2`:
```bash
$ sherpa scale out --group-name=cache --count=2 example
```

Scale in job `example` and group `cache` using the count configured in the scaling policy:
```bash
$ sherpa scale in --group-name=cache example
```

List all the scaling events currently held with the Sherpa storage backend:
```bash
$ sherpa scale status
```

Read details about the scaling event with id `f7476465-4d6e-c0de-26d0-e383c49be941`:
```
$ sherpa scale status f7476465-4d6e-c0de-26d0-e383c49be941
```

## Usage
```bash
Usage:
  sherpa scale [flags]
  sherpa scale [command]

Available Commands:
  in          Perform scaling in actions on Nomad jobs and groups
  out         Perform scaling out actions on Nomad jobs and groups
  status      Display the status output for scaling activities
```
