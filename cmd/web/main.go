package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"indibills/internal/data"
	"indibills/internal/models"
)

type application struct {
	userList *models.UserListModel
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	usersEndpoint := flag.String("endpoint", fmt.Sprintf("http://localhost:42069/v%v/users", data.VERSION), "Endpoint for Indibills Users")

	app := &application{
		userList: &models.UserListModel{Endpoint: *usersEndpoint},
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	log.Printf("Starting the server on port %s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
