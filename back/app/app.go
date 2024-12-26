package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/DanielOsorio01/enron-email-search/back/db"
)

type App struct {
	router   http.Handler
	dbClient *db.ZincClient
	config   Config
}

func New(config Config) *App {
	app := &App{
		dbClient: db.NewZincClient(
			config.dbAddr,
			config.dbUser,
			config.dbPassword),
		config: config,
	}
	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.serverPort),
		Handler: a.router,
	}
	// Ping the database to check if it's up
	err := a.dbClient.Ping()
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
