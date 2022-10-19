package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func SetupGraphQL(r chi.Router) http.Handler {
	r.Group(func(r chi.Router) {
		SetupGraphQLMiddleware(r)
		SetupGraphQLEndpoints(r)
	})
	return r
}

func SetupGraphQLMiddleware(r chi.Router) http.Handler {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	return r
}

func SetupGraphQLEndpoints(r chi.Router) http.Handler {
	h := handler.New(&handler.Config{
		Schema:   BuildGraphQLSchema(),
		Pretty:   true,
		GraphiQL: true,
	})

	r.Get("/graphql", h.ServeHTTP)
	return r
}

func BuildGraphQLSchema() *graphql.Schema {
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, _ := graphql.NewSchema(schemaConfig)
	return &schema
}
