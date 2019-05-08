# HTTP API

The Sherpa HTTP API gives you full access to a Sherpa server via HTTP. Every aspect of Sherpa can be controlled via this API.

All API routes are prefixed with /v1/, which is the current API version.

## HTTP Status Codes
The following HTTP status codes are used throughout the API. Sherpa tries to adhere to these whenever possible.

* `200` - Success with data.
* `204` - Success created without return content.
* `404` - Not found.
* `422` - Unprocessable request. An error where the supplied payload or query params are incorrect.
* `500` - Internal server error. An internal error has occurred, try again later.
