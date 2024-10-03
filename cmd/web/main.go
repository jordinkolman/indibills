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
	"os"
	"os/signal"
	"syscall"
	"time"

	"indibills/internal/data"
	"indibills/internal/models"

	"github.com/gorilla/sessions"
)

// TODO: Add Transaction List and remove userList
// App should store current user, accounts, and transactions

type application struct {
	logger   *log.Logger
	user     *models.UserModel
	accounts *models.AccountListModel
}

var store = sessions.NewCookieStore([]byte("STRING"))

func init() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 1,
		HttpOnly: true,
		}
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	endpoint := flag.String("endpoint", fmt.Sprintf("http://localhost:42069/v%v", data.VERSION), "Endpoint for Indibills Users")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		logger:   logger,
		user:     &models.UserModel{Endpoint: *endpoint},
		accounts: &models.AccountListModel{Endpoint: *endpoint},
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting the server on port %s", *addr)
	go func() {
		<-sigs
		logger.Fatal("server interrupted by user")
		os.Exit(1)
	}()

	err := srv.ListenAndServe()
	logger.Fatal(err)
}
