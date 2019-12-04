# Policy CLI

The policy command groups subcommands for interacting with policies. Users can write, read, and list policies in Sherpa. The write and delete command will only work if the Sherpa server is running using the API policy engine enabled.

## Examples

Create an example job group scaling policy:
```bash
$ sherpa policy init
```

List all policies:
```bash
$ sherpa policy list
```

Read a policy for a job named example:
```bash
$ sherpa policy read example
```

Read a policy for a job named example and group named cache:
```bash
$ sherpa policy read --policy-group-name=cache example
```

Create a policy for a job named example:
```bash
$ sherpa policy write example policy.json
```

Create a policy for a job named example and group named cache:
```bash
$ sherpa policy write --policy-group-name=cache example policy.json
```

Delete the policy for a job named example:
```bash
$ sherpa policy delete example
```

## Usage
```bash
Usage:
  sherpa policy [flags]
  sherpa policy [command]

Available Commands:
  delete      Deletes a scaling policy from Sherpa
  init        Creates an example job group scaling policy
  list        Lists all scaling policies
  read        Details scaling policies associated to a job
  write       Uploads a policy from file
```
