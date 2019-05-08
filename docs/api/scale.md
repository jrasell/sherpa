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
  "EvaluationID: "d092fdc0-e1fe-2536-67d8-43af8ca798ac"
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
  "EvaluationID: "d092fdc0-e1fe-2536-67d8-43af8ca798ac"
}
```
