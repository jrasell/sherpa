package server

import (
	"net/http"
	"net/url"

	"github.com/jrasell/sherpa/pkg/server/cluster"
)

// leaderProtectedHandler is a HTTP handler to be used on all endpoints which require a response
// from the current Sherpa cluster leader.
func leaderProtectedHandler(mem *cluster.Member, handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		isLeader, _, advAddr, err := mem.Leader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If the Sherpa server is determined to be the leader, continue and handle the request.
		if isLeader {
			handler.ServeHTTP(w, r)
			return
		}

		// If we are not the leader and there is no registered leader address, then return an error
		// to the client.
		if advAddr == "" {
			http.Error(w, "no cluster leader found", http.StatusServiceUnavailable)
			return
		}
		standbyResponseHandler(mem, w, r.URL)
	})
}

// standbyResponseHandler is responsible to handling the response to requests where the Sherpa
// server questioned is not the leader, but does know the advertised address of the leader. In this
// case the responding Sherpa server will send a redirect response to the client.
func standbyResponseHandler(mem *cluster.Member, w http.ResponseWriter, reqURL *url.URL) {

	_, _, advAddr, err := mem.Leader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If there is no leader advertise address, we cannot forward.
	if advAddr == "" {
		http.Error(w, "no cluster leader found", http.StatusServiceUnavailable)
		return
	}

	advertiseURL, err := url.Parse(advAddr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirectURL := url.URL{
		Scheme:   advertiseURL.Scheme,
		Host:     advertiseURL.Host,
		Path:     reqURL.Path,
		RawQuery: reqURL.RawPath,
	}

	if redirectURL.Scheme == "" {
		redirectURL.Scheme = "https"
	}

	w.Header().Set("Location", redirectURL.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}
