package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type RouteTable []Routes

type Routes []Route

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Handler     http.Handler
	HandlerFunc http.HandlerFunc
}

func WithRoutes(l zerolog.Logger, routes RouteTable) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for rsI := range routes {
		for rI := range routes[rsI] {
			l.Info().
				Str("method", routes[rsI][rI].Method).
				Str("path", routes[rsI][rI].Pattern).
				Msgf("mounting route endpoint %s", routes[rsI][rI].Name)

			if routes[rsI][rI].Handler != nil {
				router.Path(routes[rsI][rI].Pattern).
					Methods(routes[rsI][rI].Method).
					Name(routes[rsI][rI].Name).
					Handler(routes[rsI][rI].Handler)
			}
			if routes[rsI][rI].HandlerFunc != nil {
				router.Path(routes[rsI][rI].Pattern).
					Methods(routes[rsI][rI].Method).
					Name(routes[rsI][rI].Name).
					HandlerFunc(routes[rsI][rI].HandlerFunc)
			}
		}
	}

	return router
}
