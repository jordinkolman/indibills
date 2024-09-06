package main

import (
	"fmt"
	"indibills/internal/data"
	"net/http"
)


var (
	healthCheckUrl = fmt.Sprintf("/v%v/healthcheck", data.VERSION)
	usersUrl = fmt.Sprintf("/v%v/users", data.VERSION)
	accountsUrl = fmt.Sprintf("/v%v/accounts", data.VERSION)
)

func (app *application) route() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(healthCheckUrl, app.healthcheck)
	mux.HandleFunc(usersUrl, app.getCreateUsersHandler)
	mux.HandleFunc(accountsUrl, app.getCreateAccountsHandler)
	return mux
}
