# Policy API

The policy API allows interaction with scaling policies registered with Sherpa. The write/update and delete endpoints are dependant on using the API policy engine and are disabled otherwise.

## List Job Scaling Policies

This endpoint lists all known job scaling policies in the system registered with Sherpa.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/policies`              | `200 application/binary` |

### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/policies
```

### Sample Response

```json
{
  "my-job": {
    "my-job-group": {
      "Enabled": false,
      "MinCount": 2,
      "MaxCount": 10,
      "ScaleOutCount": 1,
      "ScaleInCount": 1
    }
  },
  "my-other-job": {
    "my-other-job-group": {
      "Enabled": true,
      "MinCount": 2,
      "MaxCount": 10,
      "ScaleOutCount": 1,
      "ScaleInCount": 1
    }
  }
}
```

## Read A Job Scaling Policy

This endpoint is used to read the scaling policy for a job.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/policy/:job_id`              | `200 application/binary` |

#### Parameters

* `:job_id` (string: required) - Specifies the ID of the job and is specified as part of the path.

### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/policy/my-job
```

### Sample Response

```json
{
  "my-job-group": {
    "Enabled": true,
    "MinCount": 2,
    "MaxCount": 10,
    "ScaleOutCount": 1,
    "ScaleInCount": 1
  }
}
```

## Read A Job Group Scaling Policy

This endpoint is used to read the scaling policy for a job.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/policy/:job_id/:group`              | `200 application/binary` |

#### Parameters

* `:job_id` (string: required) - Specifies the ID of the job and is specified as part of the path.
* `:group` (string: required) - Specifies the group name within the job and is specified as part of the path.

### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/policy/my-job/my-job-group
```

### Sample Response

```json
{
  "Enabled": true,
  "MinCount": 2,
  "MaxCount": 10,
  "ScaleOutCount": 1,
  "ScaleInCount": 1
}
```

## Create/Update A Job Scaling Policy

This endpoint can be used to create or update the scaling policy for a job. This scaling policy can contain one or more task group policies for the job.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`    | `/v1/policy/:job_id`              | `200 application/binary` |

#### Parameters

* `:job_id` (string: required) - Specifies the ID of the job and is specified as part of the path.

### Sample Payload

```json
{
  "my-group": {
    "Enabled": true,
    "MinCount": 2,
    "MaxCount": 10,
    "ScaleOutCount": 1,
    "ScaleInCount": 1
  }
}
```

### Sample Request

```
$ curl \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8000/v1/policy/my-job
```

## Create/Update A Job Group Scaling Policy

This endpoint can be used to create or update the scaling policy for a job group.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`    | `/v1/policy/:job_id/:group`              | `200 application/binary` |

#### Parameters

* `:job_id` (string: required) - Specifies the ID of the job and is specified as part of the path.
* `:group` (string: required) - Specifies the group name within the job and is specified as part of the path.

### Sample Payload

```json
{
  "Enabled": true,
  "MinCount": 2,
  "MaxCount": 10,
  "ScaleOutCount": 1,
  "ScaleInCount": 1
}
```

### Sample Request

```
$ curl \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8000/v1/policy/my-job/my-job-group
```

## Delete A Job Scaling Policy

This endpoint can be used to delete the scaling policy for a job.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE`    | `/v1/policy/:job_id`              | `200 application/binary` |

#### Parameters

* `:job_id` (string: required) - Specifies the ID of the job and is specified as part of the path.

### Sample Request

```
$ curl \
    --request DELETE \
    http://127.0.0.1:8000/v1/policy/my-job
```

## Delete A Job Group Scaling Policy

This endpoint can be used to delete the scaling policy for a job group.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE`    | `/v1/policy/:job_id/:group`              | `200 application/binary` |

#### Parameters

* `:job_id` (string: required) - Specifies the ID of the job and is specified as part of the path.
* `:group` (string: required) - Specifies the group name within the job and is specified as part of the path.

### Sample Request

```
$ curl \
    --request DELETE \
    http://127.0.0.1:8000/v1/policy/my-job/my-job-group
```
