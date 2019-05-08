# System API

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
  "PolicyStorageBackend": "Consul",
  "StrictPolicyChecking": false
}
```
