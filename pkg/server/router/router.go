package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type RouteTable []Routes

type Routes []Route

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

func WithRoutes(l zerolog.Logger, routes RouteTable) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for rsI := range routes {
		for rI := range routes[rsI] {
			l.Info().
				Str("method", routes[rsI][rI].Method).
				Str("path", routes[rsI][rI].Pattern).
				Msgf("mounting route endpoint %s", routes[rsI][rI].Name)
			router.Path(routes[rsI][rI].Pattern).
				Methods(routes[rsI][rI].Method).
				Name(routes[rsI][rI].Name).
				HandlerFunc(routes[rsI][rI].Handler)
		}
	}

	return router
}
