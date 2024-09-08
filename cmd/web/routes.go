package main

import (
	"net/http"
)

/*
	TODO: Implement the following routes:
	Transactions (/accounts/<account_id>/transactions)
	Dashboard (/ for authenticated users)
	Accounts (/accounts)
	Account Details (/accounts/<account_id>)
	Account Create (/accounts/create)
	Transaction Details (/accounts/<account_id>/transactions/<transaction_id>)
	Transaction Create (/accounts/<account_id>/transactions/create)
	Log In (/ will redirect here for unauthenticated users)
	Log Out
	Sign Up
	other pages:
	- budget (/user/<user_id>/tools/budget)
		- assign limits for various categories (and maybe specific vendors?) and see how your spending compares
		- create savings and payoff plans and see how you're doing along the way
	- assets (/user/<user_id>/assets)
		- non-monetary assets, such as investments, vehicles, homes, etc.
	- income & employment
		- total / average annual and monthly income
		- tools for 1099 workers
			- expense tracking
			- mileage tracking
			- tax payment date reminders and estimates
			- DISCLAIMER: NOT TAX OR FINANCIAL ADVICE, JUST AN ESTIMATE BASED ON THE GIVEN DATA
	- liabilities (/user/<user_id>/liabilities)
		- financial liabilites, such as auto, home, and other loans, credit cards, etc.
		- totals, payments, amortization tables? etc.
		- !!! DATA PROCESSING OPPORTUNITY !!!
	- summaries (/user/<user_id/summaries)
		- spending totals and trends by categories
		- transactions records for the given month
		- annual summaries of totals and trends only; transaction lists are only on monthly summaries
*/

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// index and login routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/login", app.loginHandler)
	mux.HandleFunc("/logout", app.logoutHandler)
	// user routes
	mux.HandleFunc("/users", app.usersView)
	mux.HandleFunc("/signup", app.userCreate)
	// authenticated routes (accounts etc.)
	return mux
}
