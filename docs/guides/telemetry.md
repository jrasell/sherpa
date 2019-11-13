# Sherpa Telemetry

The Sherpa server collects various runtime metrics about the performance that are retained for one minute. This data can be viewed either by sending the Sherpa server process a signal, or [configuring](./README.md) the server to stream data to [statsite](https://github.com/statsite/statsite), [statsd](https://github.com/statsd/statsd), or be scrapped by [Prometheus](https://prometheus.io/).

To view this data via sending a signal to the Sherpa process: on Unix, this is `USR1` while on Windows it is `BREAK`. Once Nomad receives the signal, it will dump the current telemetry information to the server's `stderr`:

```bash
[2019-05-14 18:32:50 +0100 BST][G] 'sherpa.lluna.local.runtime.sys_bytes': 72220920.000
[2019-05-14 18:32:50 +0100 BST][G] 'sherpa.lluna.local.runtime.malloc_count': 76736.000
[2019-05-14 18:32:50 +0100 BST][G] 'sherpa.lluna.local.runtime.free_count': 41066.000
[2019-05-14 18:32:50 +0100 BST][G] 'sherpa.lluna.local.runtime.heap_objects': 35670.000
[2019-05-14 18:32:50 +0100 BST][G] 'sherpa.lluna.local.runtime.total_gc_pause_ns': 39109.000
[2019-05-14 18:32:50 +0100 BST][G] 'sherpa.lluna.local.runtime.total_gc_runs': 1.000
[2019-05-14 18:32:50 +0100 BST][G] 'sherpa.lluna.local.runtime.num_goroutines': 7.000
[2019-05-14 18:32:50 +0100 BST][G] 'sherpa.lluna.local.runtime.alloc_bytes': 3044160.000
[2019-05-14 18:32:50 +0100 BST][S] 'sherpa.runtime.gc_pause_ns': Count: 1 Sum: 39109.000 LastUpdated: 2019-05-14 18:32:54.504907 +0100 BST m=+1.125110442
```

# Runtime Metrics

Runtime metrics allow operators to get insight into how the Sherpa server process is functioning.

<table class="table table-bordered table-striped">
  <tr>
    <th>Metric</th>
    <th>Description</th>
    <th>Unit</th>
    <th>Type</th>
  </tr>
  <tr>
    <td>`sherpa.runtime.num_goroutines`</td>
    <td>Number of goroutines and general load pressure indicator</td>
    <td>Number of goroutines</td>
    <td>Gauge</td>
  </tr>
  <tr>
    <td>`sherpa.runtime.alloc_bytes`</td>
    <td>Number of bytes allocated to the Sherpa process which should keep a steady state</td>
    <td>Number of bytes</td>
    <td>Gauge</td>
  </tr>
  <tr>
    <td>`sherpa.runtime.sys_bytes`</td>
    <td>This includes what is being used by Sherpa's heap and what has been reclaimed but not given back to the operating system</td>
    <td>Number of bytes</td>
    <td>Gauge</td>
  </tr>
  <tr>
    <td>`sherpa.runtime.malloc_count`</td>
    <td>Cumulative count of allocated heap objects</td>
    <td>Number of heap objects</td>
    <td>Gauge</td>
  </tr>
  <tr>
    <td>`sherpa.runtime.free_count`</td>
    <td>Number of freed objects from the heap and should steadily increase over time</td>
    <td>Number of freed objects</td>
    <td>Gauge</td>
  </tr>
  <tr>
    <td>`sherpa.runtime.heap_objects`</td>
    <td>This is a good general memory pressure indicator worth establishing a baseline and thresholds for alerting</td>
    <td>Number of objects in the heap</td>
    <td>Gauge</td>
  </tr>
  <tr>
    <td>`sherpa.runtime.total_gc_pause_ns`</td>
    <td>The total garbage collector pause time since Sherpa was last started</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.runtime.total_gc_runs`</td>
    <td>Total number of garbage collection runs since Sherpa was last started</td>
    <td>Number of operations</td>
    <td>Gauge</td>
  </tr>
</table>

# Policy Backend Metrics

Policy backend metrics allow operators to get insight into how the policy storage backend is functioning.

<table class="table table-bordered table-striped">
  <tr>
    <th>Metric</th>
    <th>Description</th>
    <th>Unit</th>
    <th>Type</th>
  </tr>
  <tr>
    <td>`sherpa.policy.memory.get_policies`</td>
    <td>Time taken to list all stored scaling policies from the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.memory.get_job_policy`</td>
    <td>Time taken to get a job scaling policy from the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.memory.get_job_group_policy`</td>
    <td>Time taken to get a job group scaling policy from the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.memory.put_job_policy`</td>
    <td>Time taken to put a job scaling policy in the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.memory.put_job_group_policy`</td>
    <td>Time taken to put a job group scaling policy in the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.memory.delete_job_policy`</td>
    <td>Time taken to delete a job scaling policy from the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.memory.delete_job_group_policy`</td>
    <td>Time taken to delete a job group scaling policy from the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.consul.get_policies`</td>
    <td>Time taken to list all stored scaling policies from the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.consul.get_job_policy`</td>
    <td>Time taken to get a job scaling policy from the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.consul.get_job_group_policy`</td>
    <td>Time taken to get a job group scaling policy from the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.consul.put_job_policy`</td>
    <td>Time taken to put a job scaling policy in the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.consul.put_job_group_policy`</td>
    <td>Time taken to put a job group scaling policy in the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.consul.delete_job_policy`</td>
    <td>Time taken to delete a job scaling policy from the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.policy.consul.delete_job_group_policy`</td>
    <td>Time taken to delete a job group scaling policy from the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
</table>


# Scaling State Backend Metrics

Scaling state backend metrics allow operators to get insight into how the scaling state backend is functioning.

<table class="table table-bordered table-striped">
  <tr>
    <th>Metric</th>
    <th>Description</th>
    <th>Unit</th>
    <th>Type</th>
  </tr>
  <tr>
    <td>`sherpa.scale.state.memory.get_events`</td>
    <td>Time taken to list all stored scaling activities from the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.scale.state.memory.get_event`</td>
    <td>Time taken to get a stored scaling activity from the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.scale.state.memory.get_latest_events`</td>
    <td>Time taken to list the latest stored scaling activities from the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.scale.state.memory.get_latest_event`</td>
    <td>Time taken to get the latest scaling activity for a job group from the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.scale.state.memory.put_event`</td>
    <td>Time taken to put a scaling activity in the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.scale.state.memory.gc`</td>
    <td>Time taken to run the scaling state garbage collector for the memory backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.scale.state.consul.get_events`</td>
    <td>Time taken to list all stored scaling activities from the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.scale.state.consul.get_event`</td>
    <td>Time taken to get a stored scaling activity from the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.scale.state.consul.get_latest_events`</td>
    <td>Time taken to list the latest stored scaling activities from the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.scale.state.consul.get_latest_event`</td>
    <td>Time taken to get the latest scaling activity for a job group from the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.scale.state.consul.put_event`</td>
    <td>Time taken to put a scaling activity in the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.scale.state.consul.gc`</td>
    <td>Time taken to run the scaling state garbage collector for the Consul backend</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
</table>

# Autoscale Metrics

Autoscale metrics allow operators to get insight into how the autoscaler is functioning.

<table class="table table-bordered table-striped">
  <tr>
    <th>Metric</th>
    <th>Description</th>
    <th>Unit</th>
    <th>Type</th>
  </tr>
  <tr>
    <td>`sherpa.autoscale.{job}.evaluation`</td>
    <td>The time taken to perform the autoscaling evaluation for the job named {job}</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.autoscale.{job}.{group}.evaluation`</td>
    <td>The time taken to perform the autoscaling evaluation for the job named {job} and group named {group}</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.autoscale.trigger.error`</td>
    <td>Number of autoscaling scale trigger errors across all jobs</td>
    <td>Number of errors</td>
    <td>Counter</td>
  </tr>
  <tr>
    <td>`sherpa.autoscale.{job}.trigger.error`</td>
    <td>Number of autoscaling scale trigger errors for the job named {job}</td>
    <td>Number of errors</td>
    <td>Counter</td>
  </tr>
  <tr>
    <td>`sherpa.autoscale.trigger.success`</td>
    <td>Number of autoscaling scale trigger successes across all jobs</td>
    <td>Number of successes</td>
    <td>Counter</td>
  </tr>
  <tr>
    <td>`sherpa.autoscale.{job}.trigger.success`</td>
    <td>Number of autoscaling scale trigger successes for the job named {job}</td>
    <td>Number of successes</td>
    <td>Counter</td>
  </tr>
  <tr>
    <td>`sherpa.autoscale.prometheus.get_value`</td>
    <td>The time taken to query Prometheus for a metric value</td>
    <td>Milliseconds</td>
    <td>Summary</td>
  </tr>
  <tr>
    <td>`sherpa.autoscale.prometheus.error`</td>
    <td>Number of errors querying Prometheus for a metric value</td>
    <td>Number of errors</td>
    <td>Counter</td>
  </tr>
    <tr>
      <td>`sherpa.autoscale.prometheus.success`</td>
      <td>Number of successful queries of Prometheus for a metric value</td>
      <td>Number of successes</td>
      <td>Counter</td>
    </tr>
</table>
