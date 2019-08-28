package v1

import (
	"net/http"
	"strconv"

	"github.com/jrasell/sherpa/pkg/policy"
	"github.com/jrasell/sherpa/pkg/scale"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func getCountFromQueryParam(r *http.Request) int {
	count := r.FormValue("count")
	if count == "" {
		return countFailed
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		return countFailed
	}
	return countInt
}

func payloadOrPolicyCount(payloadCount int, policy *policy.GroupScalingPolicy, direction scale.Direction) (int, error) {
	if payloadCount > 0 {
		return payloadCount, nil
	}

	if policy == nil && payloadCount == 0 {
		return payloadCount, errors.New("no policy configured, specify a count to scale by")
	}

	switch direction {
	case scale.DirectionIn:
		return policy.ScaleInCount, nil
	case scale.DirectionOut:
		return policy.ScaleOutCount, nil
	}

	return 0, errors.New("all possible checks failed to obtain correct count")
}

func writeJSONResponse(w http.ResponseWriter, bytes []byte, statusCode int) { // nolint:unparam
	w.Header().Set(headerKeyContentType, headerValueContentTypeJSON)
	w.WriteHeader(statusCode)
	if _, err := w.Write(bytes); err != nil {
		log.Error().Err(err).Msg("failed to write JSON response")
	}
}
