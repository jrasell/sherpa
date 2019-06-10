# Policy Guide

Scaling policies allow for the tight and close control of scaling for Nomad job groups. A job group scaling policy can currently consistent of the following fields:
* `Enabled` (bool: false) - Whether the job group is enabled for scaling to take place, used in strict checking.
* `MinCount` (int: 2) - The minimum job group count which should be running, used in strict checking.
* `MaxCount` (int: 10)  - The maximum job group count which should be running, used in strict checking.
* `ScaleInCount` (int: 1) - The number by which to decrement the job group count by when performing a scaling in action.
* `ScaleOutCount` (int: 1) - The number by which to increment the job group count by when performing a scaling in action.
* `ScaleOutCPUPercentageThreshold` (int: 80) - The CPU usage threshold in percent that will trigger an auto-scaling out event.
* `ScaleOutMemoryPercentageThreshold` (int: 80) - The memory usage threshold in percent that will trigger an auto-scaling out event.
* `ScaleInCPUPercentageThreshold` (int: 20) - The CPU usage threshold in percent that will trigger an auto-scaling in event.
* `ScaleInMemoryPercentageThreshold` (int: 20) - The memory usage threshold in percent that will trigger an auto-scaling in event.

If a partial job group scaling policy is submitted, Sherpa will merge this with the defaulted values to ensure safety and consistency. 

If using the Nomad meta policy engine all meta key fields should be prefixed with `sherpa_`, be lowercase and split the camel case using underscores. This results in the `ScaleOutCount` JSON field becoming `sherpa_scale_out_count`.

### Examples

An example of a job group scaling policy is JSON:
```json
{
  "Enabled": true,
  "MaxCount": 16,
  "MinCount": 4,
  "ScaleOutCount": 2,
  "ScaleInCount": 2,
  "ScaleOutCPUPercentageThreshold": 75,
  "ScaleOutMemoryPercentageThreshold": 75,
  "ScaleInCPUPercentageThreshold": 35,
  "ScaleInMemoryPercentageThreshold": 35
}
```

It is also possible to write a policy which represents many job groups which are part of a single job:
```json
{
  "cache": {
    "Enabled": true,
    "MaxCount": 20,
    "MinCount": 10,
    "ScaleOutCount": 3,
    "ScaleInCount": 1,
    "ScaleOutCPUPercentageThreshold": 75,
    "ScaleOutMemoryPercentageThreshold": 75,
    "ScaleInCPUPercentageThreshold": 35,
    "ScaleInMemoryPercentageThreshold": 35
  },
  "db": {
    "Enabled": true,
    "MaxCount": 5,
    "MinCount": 2,
    "ScaleOutCount": 2,
    "ScaleInCount": 2,
    "ScaleOutCPUPercentageThreshold": 85,
    "ScaleOutMemoryPercentageThreshold": 85,
    "ScaleInCPUPercentageThreshold": 15,
    "ScaleInMemoryPercentageThreshold": 15
  }
}
```

## Policy Storage Backend

In order to persist policy documents configured via the CLI and API, Sherpa can be configured to use different types of backends.

### In-Memory

The in-memory backend is the default backend and ideal for development, but its use is not suggested for any path-to-production environments. Any Sherpa server failures or restarts will result in all policy data loss.

### Consul

Consul KV provides a scalable and robust backend store for Sherpa policies. All CRUD operations will be sanitized and then passed through for action within Consul using the official SDK. Policies will be stored under the root KV as configured when running the Sherpa server, and can be browsed either using the Sherpa CLI, API or directly via Consul.

The Consul backend is preferable to in-memory as Sherpa server restarts or failures will not result in data loss. Instead the policy data relies on Consul distributed KV persistence which is proven as the highest scale.
