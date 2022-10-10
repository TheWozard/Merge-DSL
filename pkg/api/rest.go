package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func SetupREST(r chi.Router) http.Handler {
	r.Group(func(r chi.Router) {
		SetupRESTMiddleware(r)
		SetupRESTEndpoints(r)
	})
	return r
}

func SetupRESTMiddleware(r chi.Router) http.Handler {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	return r
}

func SetupRESTEndpoints(r chi.Router) http.Handler {
	h := &RESTHandler{
		validator: validator.New(),
	}
	r.Route("/document", func(r chi.Router) {
		r.Get("/{name}", h.getDocument)
	})
	return r
}

type RESTHandler struct {
	validator *validator.Validate
}

func (h *RESTHandler) getDocument(w http.ResponseWriter, r *http.Request) {
	result := struct {
		Name string `validate:"required,numeric"`
	}{
		Name: chi.URLParam(r, "name"),
	}
	if err := h.validator.StructCtx(r.Context(), result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
