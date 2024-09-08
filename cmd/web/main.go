package main

// TODO:
// If a transaction is an ATM withdrawal, it should update the users cash account
// Generate a cash account on user creation
// this is where user totals and projections should be calculated
// For auto-linked accounts w/ or w/o Plaid (if you figure out how to do it without), add option to
// add a transaction thats missing

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"indibills/internal/data"
	"indibills/internal/models"

)

// TODO: Add Transaction List and remove userList
// App should store current user, accounts, and transactions

type application struct {
	user *models.UserModel
	userList *models.UserListModel
	accountList *models.AccountListModel
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	endpoint := flag.String("endpoint", fmt.Sprintf("http://localhost:42069/v%v", data.VERSION), "Endpoint for Indibills Users")

	app := &application{
		user: &models.UserModel{Endpoint: *endpoint},
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	log.Printf("Starting the server on port %s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
