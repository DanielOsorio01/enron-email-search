package handlers

import (
	"fmt"
	"net/http"
)

type Email struct{}

func (o *Email) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("email list")
}
