package server

import (
	"net/http"

	"github.com/rs/zerolog"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// middlewareLogger ensures a log line is generated per API call and is formatted in the correct
// format according to Sherpas configuration.
func middlewareLogger(h http.Handler, l zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		respWriter := newLoggingResponseWriter(w)
		h.ServeHTTP(respWriter, req)
		l.Info().
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Str("remote-addr", req.RemoteAddr).
			Int("response-code", respWriter.statusCode).
			Msg("server responded to request")
	})
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(lrw.statusCode)
}
