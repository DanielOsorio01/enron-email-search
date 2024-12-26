package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/DanielOsorio01/enron-email-search/back/handlers"
	"github.com/DanielOsorio01/enron-email-search/back/repository/email"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/emails", a.loadEmailRoutes)

	a.router = router
}

func (a *App) loadEmailRoutes(router chi.Router) {
	email := &handlers.Email{
		Repo: &email.ZincsearchRepo{
			Client: a.dbClient,
		},
	}

	router.Get("/", email.List)
}
