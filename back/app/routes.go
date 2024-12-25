package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/DanielOsorio01/enron-email-search/back/handlers"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/emails", loadEmailRoutes)

	return router
}

func loadEmailRoutes(router chi.Router) {
	email := &handlers.Email{}

	router.Get("/", email.List)
}
