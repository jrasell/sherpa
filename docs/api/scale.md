# Scale API

## Scale Out Job Group

This endpoint can be used to scale a Nomad job group out, therefore increasing its count.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/v1/scale/out/:job_id/:group`              | `200 application/binary` |

#### Parameters

* `:job_id` (string: required) - Specifies the ID of the job and is specified as part of the path.
* `:group` (string: required) - Specifies the group name within the job and is specified as part of the path.
* `count` (int: 0) - Specifies the count which to scale the job group by. If this is not passed, Sherpa will attempt to use the value within the scaling policy.

### Sample Request

```
$ curl \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8000/v1/scale/out/my-job/my-job-group?count=2
```

### Sample Response

```json
{
  "ID": "036e4bd6-8f7d-4a8c-bf90-790790bbdc2a",
  "EvaluationID": "d092fdc0-e1fe-2536-67d8-43af8ca798ac"
}
```

## Scale In Job Group

This endpoint can be used to scale a Nomad job group in, therefore decreasing its count.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `PUT`    | `/v1/scale/in/:job_id/:group`              | `200 application/binary` |

#### Parameters

* `:job_id` (string: required) - Specifies the ID of the job and is specified as part of the path.
* `:group` (string: required) - Specifies the group name within the job and is specified as part of the path.
* `count` (int: 0) - Specifies the count which to scale the job group by. If this is not passed, Sherpa will attempt to use the value detailed within the scaling policy.

### Sample Request

```
$ curl \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8000/v1/scale/in/my-job/my-job-group?count=2
```

### Sample Response

```json
{
  "ID": "036e4bd6-8f7d-4a8c-bf90-790790bbdc2a",
  "EvaluationID": "d092fdc0-e1fe-2536-67d8-43af8ca798ac"
}
```

## List Scaling Events

This endpoint can be used to list the recent scaling events.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/scale/status`              | `200 application/binary` |

### Sample Request

```
$ curl \
    --request GET \
    http://127.0.0.1:8000/v1/scale/status
```

### Sample Response

```json
{
  "036e4bd6-8f7d-4a8c-bf90-790790bbdc2a": {
    "example2:cache": {
      "EvalID": "e05a8d0f-87f8-bda8-eb3e-885caaf50c36",
      "Source": "InternalAutoscaler",
      "Time": 1568538833630403000,
      "Status": "Completed",
      "Details": {
        "Count": 1,
        "Direction": "in"
      }
    }
  },
  "3bc8190e-b9fc-4997-bb39-3749eed5affd": {
    "example1:cache": {
      "EvalID": "ec38990e-81e2-1c99-fbf2-725e8ca6ad70",
      "Source": "InternalAutoscaler",
      "Time": 1568538893629872000,
      "Status": "Completed",
      "Details": {
        "Count": 1,
        "Direction": "in"
      }
    }
  }
}
```

## Read Scaling Event

This endpoint can be used to query a scaling event.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/scale/status/:id`              | `200 application/binary` |

### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/scale/status/3bc8190e-b9fc-4997-bb39-3749eed5affd
```

### Sample Response

```json
"3bc8190e-b9fc-4997-bb39-3749eed5affd": {
  "example1:cache": {
    "EvalID": "ec38990e-81e2-1c99-fbf2-725e8ca6ad70",
    "Source": "InternalAutoscaler",
    "Time": 1568538893629872000,
    "Status": "Completed",
    "Details": {
      "Count": 1,
      "Direction": "in"
    }
  }
}
```
