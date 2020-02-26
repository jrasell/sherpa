# Sherpa Policies

Scaling policies allow for the tight and close control of scaling for Nomad job groups. A job group scaling policy has a number of required, and optional parameters; as well as defaults if some required parameters are not configured.

### Required Params
* `Enabled` (bool: false) - Whether the job group is enabled for scaling to take place.
* `MinCount` (int: 2) - The minimum job group count which should be running.
* `MaxCount` (int: 10)  - The maximum job group count which should be running.
* `Cooldown` (int: 180) - Cooldown is a time period in seconds. Once a scaling action has been triggered on the desired group, another action will not be triggered until the cooldown period has passed.
* `ScaleInCount` (int: 1) - The number by which to decrement the job group count by when performing a scaling in action.
* `ScaleOutCount` (int: 1) - The number by which to increment the job group count by when performing a scaling in action.

### Optional Nomad Check Params
The Nomad checks parameters tell the autoscaler to check the resource consumption of the job group using metrics gathered from the Nomad API. It compares the actual resource usage against the allocated resources as configured within the job specification.

* `ScaleOutCPUPercentageThreshold` (float64) - The percentage utilisation threshold of CPU, which if broken will result in a scaling out of the job group.
* `ScaleOutMemoryPercentageThreshold` (float64) - The percentage utilisation threshold of memory, which if broken will result in a scaling out of the job group.
* `ScaleInCPUPercentageThreshold` (float64) - The percentage utilisation threshold of CPU, which if broken will result in a scaling in of the job group.
* `ScaleInMemoryPercentageThreshold` (float64) - The percentage utilisation threshold of memory, which if broken will result in a scaling in of the job group.

### Optional External Checks Params
The optional external checks are a map of checks which utilise external sources for metrics values. The obtained value is then compared via the `ComparisonOperator` to the `ComparisonValue`. The map key is a free-form name, operators should use to clearly identify the check.

* `Enabled` (bool) - Whether this check should be run or not.
* `Provider` (string) - The metrics provider to utilise for obtaining the value for comparison. Currently only `prometheus` is supported.
* `Query` (string) - The query which can be run against the provider. The style is specific to the provider; examples of which can be seen below. It is important to note that this query should result in the return of a single data-point.
* `ComparisonOperator` (string) - The equality operator used to compare the metric value with the threshold. Currently this supports `greater-than` and `less-than`.
* `ComparisonValue` (string) - The threshold value which the metric value will be compared against.
* `Action` (string) - The action to take if the threshold check is broken. This can be either `scale-in` or `scale-out`.

## Nomad Meta Policies
Scaling policies can be configured within Nomad job specification [meta stanzas](https://www.nomadproject.io/docs/job-specification/meta.html). When this features is enabled, Sherpa will monitor jobs, and update its internal policies to match those found on the cluster. The parameter names are prefixed within sherpa, use lowercase and break the camel case with underscores.  
* `sherpa_enabled`
* `sherpa_cooldown`
* `sherpa_max_count`
* `sherpa_min_count`
* `sherpa_scale_in_count`
* `sherpa_scale_out_count`
* `sherpa_scale_out_cpu_percentage_threshold`
* `sherpa_scale_out_memory_percentage_threshold`
* `sherpa_scale_in_cpu_percentage_threshold`
* `sherpa_scale_in_memory_percentage_threshold`
* `sherpa_external_checks`

Due to the string:string nature of Nomad meta keys, the `sherpa_external_checks` needs to be formatted and escaped correctly to be decoded. The below example shows the Nomad meta value for an external check using Prometheus.
```
"sherpa_external_checks": "{\"ExternalChecks\":{\"prometheus_test\":{\"Enabled\":true,\"Provider\":\"prometheus\",\"Query\":\"job:nomad_redis_cache_memory:percentage\",\"ComparisonOperator\":\"less-than\",\"ComparisonValue\":30,\"Action\":\"scale-in\"}}}
```

## Examples
An example job group policy which configures Sherpa to perform all the Nomad checks and no external checks.
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

An example job group policy which configures Sherpa to perform two external Prometheus checks and no Nomad resource checks.
```json
{
  "Enabled": true,
  "MaxCount": 16,
  "MinCount": 1,
  "ScaleOutCount": 1,
  "ScaleInCount": 1,
  "ExternalChecks": {
    "prometheus_memory_in": {
      "Enabled": true,
      "Provider": "prometheus",
      "Query": "sum(nomad_client_allocs_memory_usage{task_group='cache'})/sum(nomad_client_allocs_memory_allocated{task_group='cache'})*100",
      "ComparisonOperator": "less-than",
      "ComparisonValue": 30,
      "Action": "scale-in"
    },
    "prometheus_memory_out": {
      "Enabled": true,
      "Provider": "prometheus",
      "Query": "sum(nomad_client_allocs_memory_usage{task_group='cache'})/sum(nomad_client_allocs_memory_allocated{task_group='cache'})*100",
      "ComparisonOperator": "greater-than",
      "ComparisonValue": 80,
      "Action": "scale-out"
    }
  }
}
```

An example job group policy which configures Sherpa to perform two external InfluxDB checks and no Nomad resource checks. The query below is gathering the mean average over the last 10 minutes, and is querying from the telegraf database using the default retention strategy. Additional details on querying InfluxDB can be found [here](https://docs.influxdata.com/influxdb/v1.7/query_language/data_exploration/#the-basic-select-statement).
```json
{
  "Enabled": true,
  "MaxCount": 16,
  "MinCount": 1,
  "ScaleOutCount": 1,
  "ScaleInCount": 1,
  "ExternalChecks": {
    "influxdb_cpu_in": {
      "Enabled": true,
      "Provider": "influxdb",
      "Query": "SELECT mean(gauge) as cpu FROM telegraf..nomad_client_allocs_cpu_total_ticks WHERE job = 'cache' and time >= now() -10m",
      "ComparisonOperator": "less-than",
      "ComparisonValue": 30,
      "Action": "scale-in"
    },
    "influxdb_cpu_out": {
      "Enabled": true,
      "Provider": "influxdb",
      "Query": "SELECT mean(gauge) as cpu FROM telegraf..nomad_client_allocs_cpu_total_ticks WHERE job = 'cache' and time >= now() -10m",
      "ComparisonOperator": "greater-than",
      "ComparisonValue": 80,
      "Action": "scale-out"
    }
  }
}⏎
```

A Nomad meta stanza example configuring both Nomad and external checks.
```
"sherpa_enabled"                               = "true"
"sherpa_cooldown"                              = "120"
"sherpa_max_count"                             = "13"
"sherpa_min_count"                             = "3"
"sherpa_scale_in_count"                        = "1"
"sherpa_scale_out_count"                       = "1"
"sherpa_scale_out_cpu_percentage_threshold"    = "80"
"sherpa_scale_out_memory_percentage_threshold" = "80"
"sherpa_scale_in_cpu_percentage_threshold"     = "20"
"sherpa_scale_in_memory_percentage_threshold"  = "20"
"sherpa_external_checks"                       = "{\"prometheus_test\":{\"Enabled\":true,\"Provider\":\"prometheus\",\"Query\":\"job:nomad_redis_cache_memory:percentage\",\"ComparisonOperator\":\"less-than\",\"ComparisonValue\":30,\"Action\":\"scale-in\"}}"
```
