package main

import (
	"fmt"
	"indibills/internal/data"

	"github.com/gorilla/mux"
)

// TODO: figure out how to convert to HTTPS and encrypted transmission

var (
	healthCheckUrl = fmt.Sprintf("/v%v/healthcheck", data.VERSION)
	userUrl       = fmt.Sprintf("/v%v/users/{email}", data.VERSION)
	accountsUrl    = fmt.Sprintf("/v%v/{user_id}/accounts", data.VERSION)
)

func (app *application) route() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(healthCheckUrl, app.healthcheck)
	r.HandleFunc(userUrl, app.getCreateUserHandler)
	r.HandleFunc(accountsUrl, app.getCreateAccountsHandler)
	return r
}
