package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/DanielOsorio01/enron-email-search/back/handlers"
	"github.com/DanielOsorio01/enron-email-search/back/repository/email"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// Allow only your frontend's origin
		AllowedOrigins: []string{"http://localhost:8080", "http://127.0.0.1:8080", "https://your-production-site.com"},
		// Restrict allowed HTTP methods
		AllowedMethods: []string{"GET"},
		// Allow specific headers (Content-Type for JSON requests, Authorization for tokens, etc.)
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		// Allow credentials (cookies, authorization headers, etc.) only if necessary
		AllowCredentials: false,
		// Cache preflight responses for better performance
		MaxAge: 300,
	}))
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
