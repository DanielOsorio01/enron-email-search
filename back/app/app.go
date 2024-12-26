package app

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type App struct {
	router   http.Handler
	dbClient *http.Client
}

func New() *App {
	return &App{
		router:   loadRoutes(),
		dbClient: &http.Client{},
	}
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	// Ping the database to check if it's up
	_, err := a.dbClient.Get("http://localhost:4080")
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Database is up. Starting server...")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(ctxTimeout)
	}

}
