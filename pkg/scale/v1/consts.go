package v1

import "github.com/pkg/errors"

type ScaleReq struct {
	Count int `json:"count"`
}

type ScaleResp struct {
	EvalID string `json:"eval_id"`
}

const (
	countFailed                = 0
	headerKeyContentType       = "Content-Type"
	headerValueContentTypeJSON = "application/json; charset=utf-8"
)

var (
	errInternalScaleOutNoPolicy = errors.New("scale out forbidden, no scaling policy found")
	errInternalScaleInNoPolicy  = errors.New("scale in forbidden, no scaling policy found")
)
