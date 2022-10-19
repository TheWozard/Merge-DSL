package main

import (
	"merge-dsl/pkg/api"
	"merge-dsl/pkg/config"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.GetNewConfig()
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		if cfg.DocumentAPI {
			// This route defines predefined document resolution.
			// Routes will be of the form /document/{name} and will resolve then named document.
			api.SetupREST(r)
		}
		if cfg.GraphQLAPI {
			// This route is identical to document except it is done through graphql
			api.SetupGraphQL(r)
		}
		if cfg.ReferenceAPI {
			// This route provides a sandbox endpoint for
		}
		if cfg.SandboxAPI {
			// TODO: This
		}
	})

	r.Get("/health", config.HealthEndpoint{Config: cfg}.ServeHealth)
	http.ListenAndServe(":8000", r) //nolint:errcheck
}
