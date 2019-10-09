# System API

## Get Server Leader

This endpoint can  be used to identify the current cluster leader and the storage HA capability.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/system/leader`              | `200 application/binary` |

### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/system/leader
```

### Sample Response

```json
{
  "IsSelf": true,
  "HAEnabled": true,
  "LeaderAddress": "127.0.0.1:8000",
  "LeaderClusterAddress": "http://127.0.0.1:8000"
}
```

## Get Server Health

This endpoint can be used to query the Sherpa server health status.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/system/health`              | `200 application/binary` |


### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/system/health
```

### Sample Response

```json
{
  "status": "ok"
}
```

## Get Server Info

This endpoint can be used to query the Sherpa server configuration information.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/system/info`              | `200 application/binary` |


### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/system/info
```

### Sample Response

```json
{
  "InternalAutoScalingEngine": true,
  "NomadAddress": "http://localhost:4646",
  "PolicyEngine": "Sherpa API",
  "StorageBackend": "Consul",
  "StrictPolicyChecking": false
}
```

## Get Server Metrics

This endpoint can be used to query the Sherpa server for its latest telemetry data.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/system/metrics`              | `200 application/binary` |


### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/system/metrics
```

### Sample Response

```json
{
  "Timestamp": "2019-05-14 20:41:30 +0000 UTC",
  "Gauges": [
    {
      "Name": "sherpa.runtime.alloc_bytes",
      "Value": 1745680,
      "Labels": {}
    },
    {
      "Name": "sherpa.runtime.free_count",
      "Value": 230471,
      "Labels": {}
    },
    {
      "Name": "sherpa.runtime.heap_objects",
      "Value": 17573,
      "Labels": {}
    },
    {
      "Name": "sherpa.runtime.malloc_count",
      "Value": 248044,
      "Labels": {}
    },
    {
      "Name": "sherpa.runtime.num_goroutines",
      "Value": 5,
      "Labels": {}
    },
    {
      "Name": "sherpa.runtime.sys_bytes",
      "Value": 72810744,
      "Labels": {}
    },
    {
      "Name": "sherpa.runtime.total_gc_pause_ns",
      "Value": 660024,
      "Labels": {}
    },
    {
      "Name": "sherpa.runtime.total_gc_runs",
      "Value": 13,
      "Labels": {}
    }
  ],
  "Points": [],
  "Counters": [],
  "Samples": []
}
```
