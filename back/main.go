package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/DanielOsorio01/enron-email-search/back/app"
)

func main() {
	// Crea la instancia de la aplicación con la configuración cargada
	app := app.New(*app.LoadConfig())

	// Define a context that will be canceled when a SIGINT is sent
	// to have graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Start the app
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
