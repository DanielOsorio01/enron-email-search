package main

import (
	"context"
	"fmt"

	"github.com/DanielOsorio01/enron-email-search/back/app"
)

func main() {
	app := app.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app: %w", err)
	}
}
